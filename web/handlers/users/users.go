package users

import (
	"net/http"
	"path"

	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/cmelgarejo/minesweeper-svc/database/repo"
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	"github.com/cmelgarejo/minesweeper-svc/web/models/requests"
	"github.com/cmelgarejo/minesweeper-svc/web/services/common"
	"github.com/loopcontext/msgcat"
)

type UserHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	// For Admins
	// List(w http.ResponseWriter, r *http.Request)
	// Delete(w http.ResponseWriter, r *http.Request)
}

type UserHandlerSvc struct {
	log            logger.Logger
	catalog        msgcat.MessageCatalog
	authRepo       repo.AuthRepo
	requestHelper  common.RequestHelper
	responseHelper common.ResponseHelper
}

func NewUserHandlerSvc(log logger.Logger, catalog msgcat.MessageCatalog, authRepo repo.AuthRepo,
	requestHelper common.RequestHelper, responseHelper common.ResponseHelper) UserHandler {
	return &UserHandlerSvc{
		log:            log,
		catalog:        catalog,
		authRepo:       authRepo,
		responseHelper: responseHelper,
		requestHelper:  requestHelper,
	}
}

// SignIn godoc
// @Summary Sign in user of minesweeper
// @Description Sign in user of minesweeper and returns an API Key
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/auth/signIn [post]
// @Param credentials body requests.Credentials true "Credentials"
func (svc *UserHandlerSvc) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var credentials requests.Credentials
	err, status := svc.requestHelper.DecodeJSONBody(w, r, &credentials)
	if err != nil {
		svc.responseHelper.Error(w, r, status, err)

		return
	}
	apiKey, err := svc.authRepo.SignIn(ctx, credentials.Username, credentials.Password)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, apiKey)
}

// Create godoc
// @Summary Creates an user of minesweeper
// @Description Creates an user of minesweeper and returns an userID
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} responses.Response
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/auth/user [post]
// @Param userInput body requests.UserInput true "User Input"
func (svc *UserHandlerSvc) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user *models.User
	err, status := svc.requestHelper.DecodeJSONBody(w, r, &user)
	if err != nil {
		svc.responseHelper.Error(w, r, status, err)

		return
	}
	user, err = svc.authRepo.UpsertUser(ctx, nil, user)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusCreated, user)
}

// Read godoc
// @Summary Gets the information of a minesweeper user
// @Description Gets the information of a minesweeper user, fields and users
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/auth/user/{id} [get]
// @Param id path string true "User ID"
func (svc *UserHandlerSvc) Read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	// if err != nil {
	// 	svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
	// 	return
	// }
	userID := path.Base(r.URL.Path)
	user, err := svc.authRepo.Read(ctx, userID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, user)
}

// Update godoc
// @Summary Updates an user of minesweeper
// @Description Updates an user of minesweeper and returns an userID
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/auth/user/{id} [put]
// @Param id path string true "User ID"
// @Param userInput body requests.UserInput true "User Input"
func (svc *UserHandlerSvc) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	// if err != nil {
	// 	svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
	// 	return
	// }
	var user *models.User
	err, status := svc.requestHelper.DecodeJSONBody(w, r, &user)
	if err != nil {
		svc.responseHelper.Error(w, r, status, err)

		return
	}
	userID := path.Base(r.URL.Path)
	user, err = svc.authRepo.UpsertUser(ctx, &userID, user)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, user)
}

// func (svc *UserHandlerSvc) List(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	currentUser, err := svc.authSvc.GetCurrentUser(ctx)
// 	if err != nil {
// 		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

// 		return
// 	}
// 	list, err := svc.userRepo.ListByClientID(ctx, currentUser.Client.ID)
// 	if err != nil {
// 		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

// 		return
// 	}
// 	svc.responseHelper.Send(w, r, http.StatusOK, list)
// }
