package services

import (
	"context"

	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/cmelgarejo/minesweeper-svc/resources/messages/codes"
	"github.com/cmelgarejo/minesweeper-svc/utils"
	"github.com/loopcontext/msgcat"
	"github.com/valyala/fasthttp"
)

type AuthSvc interface {
	GetCurrentUser(ctx context.Context) (currentUser *models.User, err error)
}

type AuthSvcImpl struct {
	catalog msgcat.MessageCatalog
}

func NewAuthSvcImpl(catalog msgcat.MessageCatalog) AuthSvc {
	return &AuthSvcImpl{
		catalog: catalog,
	}
}

func (svc *AuthSvcImpl) GetCurrentUser(ctx context.Context) (currentUser *models.User, err error) {
	if ctx, ok := ctx.(*fasthttp.RequestCtx); ok {
		currentUser = ctx.UserValue(string(utils.CurrrentUserCtxKey)).(*models.User)
	} else {
		currentUser = ctx.Value(utils.CurrrentUserCtxKey).(*models.User)
	}
	if currentUser == nil {
		err = svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeHelperCurrentUserNotFound)
	}

	return
}

func (svc *AuthSvcImpl) SignIn(ctx context.Context) (apiKey string, err error) {

	// if apiKey == nil {
	// 	err = svc.catalog.GetErrorWithCtx(ctx, services.MsgCodeHelperCurrentUserNotFound)
	// }

	return
}
