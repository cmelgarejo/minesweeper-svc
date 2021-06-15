package games

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/cmelgarejo/minesweeper-svc/database/repo"
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	"github.com/cmelgarejo/minesweeper-svc/web/game/engine"
	"github.com/cmelgarejo/minesweeper-svc/web/game/service"
	"github.com/cmelgarejo/minesweeper-svc/web/models/requests"
	"github.com/cmelgarejo/minesweeper-svc/web/services"
	"github.com/cmelgarejo/minesweeper-svc/web/services/common"
	"github.com/loopcontext/msgcat"
)

type GameHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Click(w http.ResponseWriter, r *http.Request)
	// For Admins
	// List(w http.ResponseWriter, r *http.Request)
	// Delete(w http.ResponseWriter, r *http.Request)
}

type GameHandlerSvc struct {
	log            logger.Logger
	catalog        msgcat.MessageCatalog
	gameRepo       repo.GameRepo
	gameEngineSvc  service.MineSweeperGameSvc
	authSvc        services.AuthSvc
	requestHelper  common.RequestHelper
	responseHelper common.ResponseHelper
}

func NewGameHandlerSvc(log logger.Logger, catalog msgcat.MessageCatalog,
	gameRepo repo.GameRepo, gameEngineSvc service.MineSweeperGameSvc,
	authSvc services.AuthSvc, requestHelper common.RequestHelper, responseHelper common.ResponseHelper) GameHandler {

	return &GameHandlerSvc{
		log:            log,
		catalog:        catalog,
		gameRepo:       gameRepo,
		gameEngineSvc:  gameEngineSvc,
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
// @Router /v1/api/game [post]
// @Param X-API-KEY header string true "API Key"
// @Param gameInput body requests.GameCreateInput true "Game Input"
func (svc *GameHandlerSvc) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	var input requests.GameCreateInput
	err, status := svc.requestHelper.DecodeJSONBody(w, r, &input)
	if err != nil {
		svc.responseHelper.Error(w, r, status, err)

		return
	}
	game, err := svc.gameEngineSvc.CreateGame(input.Rows, input.Cols, input.Mines)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	var mf struct {
		Minefield [][]engine.Field `json:"mineField"`
	}
	mf.Minefield = game.MineField
	b, _ := json.Marshal(mf)
	gameStore := &models.Game{
		Rows:        1,
		Cols:        1,
		Mines:       1,
		CreatedByID: currentUser.ID,
	}
	err = json.Unmarshal(b, &gameStore.MineField)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	_, err = svc.gameRepo.UpsertGame(ctx, gameStore)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, game.ID)
}

// Read godoc
// @Summary Gets the information of a minesweeper game
// @Description Gets the information of a minesweeper game, fields and users
// @Tags game
// @Accept json
// @Produce json
// @Success 201 {object} engine.Game
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/api/game/{id} [get]
// @Param X-API-KEY header string true "API Key"
// @Param id path string true "Game ID"
func (svc *GameHandlerSvc) Read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	// if err != nil {
	// 	svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
	// 	return
	// }
	gameID := path.Base(r.URL.Path)
	list, err := svc.gameRepo.Read(ctx, gameID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, list)
}

// Click godoc
// @Summary Updates a game of minesweeper
// @Description Updates a game of minesweeper and returns a gameID
// @Tags game
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/api/game/{id} [put]
// @Param id path string true "Game ID"
// @Param X-API-KEY header string true "API Key"
// @Param gameInput body requests.GameInput true "Game Input"
func (svc *GameHandlerSvc) Click(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	var input requests.GameInput
	err, status := svc.requestHelper.DecodeJSONBody(w, r, &input)
	if err != nil {
		svc.responseHelper.Error(w, r, status, err)

		return
	}
	gameID := path.Base(r.URL.Path)
	err = svc.gameEngineSvc.Click(gameID, currentUser.Fullname, input.GetClickType(), input.Row, input.Col)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)

		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, nil)
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
