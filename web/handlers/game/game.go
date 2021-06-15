package game

import (
	"net/http"
	"path"

	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/cmelgarejo/minesweeper-svc/database/repo"
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	"github.com/cmelgarejo/minesweeper-svc/web/services"
	"github.com/cmelgarejo/minesweeper-svc/web/services/common"
	"github.com/loopcontext/msgcat"
)

type GameHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	// Delete(w http.ResponseWriter, r *http.Request)
	// List(w http.ResponseWriter, r *http.Request)
}

type GameHandlerSvc struct {
	log            logger.Logger
	catalog        msgcat.MessageCatalog
	gameRepo       repo.GameRepo
	authSvc        services.AuthSvc
	requestHelper  common.RequestHelper
	responseHelper common.ResponseHelper
}

func NewGameHandlerSvc(log logger.Logger, catalog msgcat.MessageCatalog,
	gameRepo repo.GameRepo, authSvc services.AuthSvc,
	requestHelper common.RequestHelper, responseHelper common.ResponseHelper) GameHandler {

	return &GameHandlerSvc{
		log:            log,
		catalog:        catalog,
		gameRepo:       gameRepo,
		authSvc:        authSvc,
		responseHelper: responseHelper,
		requestHelper:  requestHelper,
	}
}

// Create godoc
// @Summary Creates a game of minesweeper
// @Description Creates a game of minesweeper and returns a gameID
// @Tags game
// @Accept json
// @Produce json
// @Success 201 {object} engine.Game
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /api/game [post]
func (svc *GameHandlerSvc) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	game, err := svc.gameRepo.UpsertGame(ctx, &models.Game{
		Rows:      1,
		Cols:      1,
		Mines:     1,
		CreatedBy: *currentUser,
	})
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, game.ID)
}

// Read godoc
// @Summary Updates a game of minesweeper
// @Description Updates a game of minesweeper and returns a gameID
// @Tags game
// @Accept json
// @Produce json
// @Success 201 {object} engine.Game
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /api/game/{id} [put]
// @Param id path string true "Game ID"
func (svc *GameHandlerSvc) Read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	// if err != nil {
	// 	svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

	// 	return
	// }
	clientID := path.Base(r.URL.Path)
	list, err := svc.gameRepo.Read(ctx, clientID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, list)
}

// Update godoc
// @Summary Gets the information of a minesweeper game
// @Description Gets the information of a minesweeper game, fields and users
// @Tags game
// @Accept json
// @Produce json
// @Success 201 {object} engine.Game
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /api/game/{id} [get]
// @Param id path string true "Game ID"
func (svc *GameHandlerSvc) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	game, err := svc.gameRepo.UpsertGame(ctx, &models.Game{
		Rows:      1,
		Cols:      1,
		Mines:     1,
		CreatedBy: *currentUser,
	})
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, game.ID)
}

// func (svc *GameHandlerSvc) List(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	currentUser, err := svc.authSvc.GetCurrentUser(ctx)
// 	if err != nil {
// 		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

// 		return
// 	}
// 	list, err := svc.gameRepo.ListByClientID(ctx, currentUser.Client.ID)
// 	if err != nil {
// 		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

// 		return
// 	}
// 	svc.responseHelper.Send(w, r, http.StatusOK, list)
// }
