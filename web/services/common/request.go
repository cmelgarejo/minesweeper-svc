package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/cmelgarejo/minesweeper-svc/resources/messages/codes"
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	"github.com/golang/gddo/httputil/header"
	"github.com/loopcontext/msgcat"
)

type RequestHelper interface {
	DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) (error, int)
}

type RequestHelperSvc struct {
	log     logger.Logger
	catalog msgcat.MessageCatalog
}

func NewRequestHelperSvc(log logger.Logger, catalog msgcat.MessageCatalog) RequestHelper {
	return &RequestHelperSvc{
		log:     log,
		catalog: catalog,
	}
}

func (svc *RequestHelperSvc) DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) (error, int) {
	ctx := r.Context()
	if r.Header.Get(ContentTypeKey) != "" {
		value, _ := header.ParseValueAndParams(r.Header, ContentTypeKey)
		if value != AppTypeJSON {
			return svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeReqHelperNotJSON), http.StatusUnsupportedMediaType
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeReqHelperBadlyFormedAtPos, syntaxError.Offset),
				http.StatusBadRequest

		case errors.Is(err, io.ErrUnexpectedEOF):
			return svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeReqHelperBadlyFormed), http.StatusBadRequest

		case errors.As(err, &unmarshalTypeError):
			return svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeReqHelperInvalidValue,
				unmarshalTypeError.Field, unmarshalTypeError.Offset), http.StatusBadRequest

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeReqHelperUnknownField, fieldName), http.StatusBadRequest

		case errors.Is(err, io.EOF):
			return svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeReqHelperReqBodyEmpty), http.StatusBadRequest

		case err.Error() == "http: request body too large":
			return svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeReqHelperLimitSize), http.StatusRequestEntityTooLarge

		default:
			return err, http.StatusInternalServerError
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeReqHelperLimit1Obj), http.StatusBadRequest
	}

	return nil, http.StatusOK
}
