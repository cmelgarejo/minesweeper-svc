package games

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"

	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/cmelgarejo/minesweeper-svc/database/repo"
	"github.com/cmelgarejo/minesweeper-svc/utils"
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
	List(w http.ResponseWriter, r *http.Request)
	Start(w http.ResponseWriter, r *http.Request)
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
// @Router /v1/api/games [post]
// @Param X-API-KEY header string true "API Key" default(587fa65a9c375165828a6fbb5f9963a7)
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
	game, err := svc.gameEngineSvc.CreateGame(input.Rows, input.Cols, input.Mines, currentUser.Fullname)
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
		Rows:        input.Rows,
		Cols:        input.Cols,
		Mines:       input.Mines,
		CreatedByID: currentUser.Fullname,
	}
	err = json.Unmarshal(b, &gameStore.GameState)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	gameStore.ID = game.ID
	_, err = svc.gameRepo.UpsertGame(ctx, nil, gameStore)
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
// @Router /v1/api/games/{id} [get]
// @Param id path string true "Game ID" default(ef99fdfd88565827ad330d83aac5fbaa)
// @Param X-API-KEY header string true "API Key" default(587fa65a9c375165828a6fbb5f9963a7)
func (svc *GameHandlerSvc) Read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	gameID := path.Base(r.URL.Path)
	game, err := svc.gameEngineSvc.GetGame(gameID)
	if err != nil {
		// if the game is not cached in the game engine service, lets pick from db
		if errors.Is(err, service.ErrGameNotFound) {
			gameStore, err := svc.gameRepo.Read(ctx, gameID)
			if err == nil {
				game = gameStore.GetGameState()
			}
		}
		if err != nil {
			svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, game)
}

// Click godoc
// @Summary Clicks field on a game of minesweeper
// @Description Clicks field on a game of minesweeper and returns the mine field state
// @Tags game
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/api/games/{id} [patch]
// @Param id path string true "Game ID" default(ef99fdfd88565827ad330d83aac5fbaa)
// @Param X-API-KEY header string true "API Key" default(587fa65a9c375165828a6fbb5f9963a7)
// @Param gameInput body requests.GameInput true "Game Input"
func (svc *GameHandlerSvc) Click(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, err := svc.authSvc.GetCurrentUser(r.Context())
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	// Get player input
	var input requests.GameInput
	err, status := svc.requestHelper.DecodeJSONBody(w, r, &input)
	if err != nil {
		svc.responseHelper.Error(w, r, status, err)
		return
	}
	// Get game id
	gameID := path.Base(r.URL.Path)
	// Get game data from store and sync
	gameStore, err := svc.gameRepo.Read(ctx, gameID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	_ = svc.gameEngineSvc.UpdateGameState(gameID, gameStore.GetGameState())
	// Click in the game engine
	err = svc.gameEngineSvc.Click(gameID, currentUser.Fullname, input.GetClickType(), input.Row, input.Col)
	if err != nil {
		if errors.Is(err, engine.ErrDefeat) {
			gameStore.Status = engine.GameStatusDefeat
		} else {
			svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
	game, err := svc.gameEngineSvc.GetGame(gameID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	gameStore.UpdateGameState(game)
	_, err = svc.gameRepo.UpsertGame(ctx, &gameID, gameStore)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	svc.responseHelper.Send(w, r, http.StatusOK, game)
}

// List godoc
// @Summary Gets a list of games
// @Description Gets a list of games, only Admins can see it
// @Tags game
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/api/games [get]
// @Param X-API-KEY header string true "API Key" default(587fa65a9c375165828a6fbb5f9963a7)
func (svc *GameHandlerSvc) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	if err != nil || !currentUser.Admin {
		if !currentUser.Admin {
			err = svc.catalog.GetErrorWithCtx(ctx, 2)
		}
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	var list map[string]*engine.Game
	list, err = svc.gameEngineSvc.GetGameList()
	if len(list) < 1 {
		gamesStore, err := svc.gameRepo.List(ctx)
		if err == nil {
			for _, game := range gamesStore {
				list[game.ID] = &engine.Game{}
				b, _ := utils.ToJSONBytes(game.GameState)
				_ = utils.ToObject(b, list[game.ID])

			}
		}
	}
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	svc.responseHelper.Send(w, r, http.StatusOK, list)
}

// StartGame godoc
// @Summary Starts a game of minesweeper
// @Description Starts a game of minesweeper and returns the mine field state
// @Tags game
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.ResponseError
// @Failure 404 {object} responses.ResponseError
// @Failure 500 {object} responses.ResponseError
// @Router /v1/api/games/start/{id} [post]
// @Param id path string true "Game ID" default(ef99fdfd88565827ad330d83aac5fbaa)
// @Param X-API-KEY header string true "API Key" default(587fa65a9c375165828a6fbb5f9963a7)
func (svc *GameHandlerSvc) Start(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// FUTURE: I could store who started then game...
	// currentUser, err := svc.authSvc.GetCurrentUser(ctx)
	// if err != nil {
	// 	svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
	// 	return
	// }
	// Get game id
	gameID := path.Base(r.URL.Path)
	// Get game data from store and sync
	gameStore, err := svc.gameRepo.Read(ctx, gameID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	_ = svc.gameEngineSvc.UpdateGameState(gameID, gameStore.GetGameState())
	// Start the game
	err = svc.gameEngineSvc.StartGame(gameID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	game, err := svc.gameEngineSvc.GetGame(gameID)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	gameStore.UpdateGameState(game)
	_, err = svc.gameRepo.UpsertGame(ctx, &gameID, gameStore)
	if err != nil {
		svc.responseHelper.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	svc.responseHelper.Send(w, r, http.StatusOK, game)
}
