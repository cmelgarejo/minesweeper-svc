package service

import (
	"errors"

	"github.com/cmelgarejo/minesweeper-svc/web/game/engine"
)

var (
	ErrGameNotFound = errors.New("Game not found")
)

// MineSweeperGame represents a minesweeper game service
type MineSweeperGameSvc interface {
	CreateGame(rows, cols, mines int, createdBy string) (game *engine.Game, err error)
	StartGame(gameID string) (err error)
	GetGame(gameID string) (game *engine.Game, err error)
	Click(gameID string, user string, clickType engine.ClickType, row, col int) (err error)
	GetGameList() (games map[string]*engine.Game, err error)
	UpdateGameState(gameID string, game *engine.Game) (err error)
}

// MineSweeperGameSvcImpl implementing struct of a minesweeper game service
type MineSweeperGameSvcImpl struct {
	// Will hold persistance instance, logger, and/or any instances needed by the game
	games map[string]*engine.Game
}

func (ms *MineSweeperGameSvcImpl) NewMineSweeperSvc() MineSweeperGameSvc {
	return &MineSweeperGameSvcImpl{
		games: make(map[string]*engine.Game),
	}
}

func (ms *MineSweeperGameSvcImpl) CreateGame(rows, cols, mines int, createdBy string) (game *engine.Game, err error) {
	game = engine.NewGame(rows, cols, mines, createdBy)
	ms.games[game.ID] = game
	return ms.games[game.ID], err
}

func (ms *MineSweeperGameSvcImpl) StartGame(gameID string) (err error) {
	game, err := ms.GetGame(gameID)
	if err != nil {
		return err
	}

	return game.Start()
}

func (ms *MineSweeperGameSvcImpl) GetGame(gameID string) (game *engine.Game, err error) {
	if _, found := ms.games[gameID]; found {
		return ms.games[gameID], nil
	}
	return nil, ErrGameNotFound
}

func (ms *MineSweeperGameSvcImpl) Click(gameID string, clickedBy string, clickType engine.ClickType, row, col int) (err error) {
	game, err := ms.GetGame(gameID)
	if err != nil {
		return err
	}
	return game.Click(clickedBy, clickType, row, col)
}

func (ms *MineSweeperGameSvcImpl) GetGameList() (games map[string]*engine.Game, err error) {
	return ms.games, nil
}

func (ms *MineSweeperGameSvcImpl) UpdateGameState(gameID string, game *engine.Game) (err error) {
	ms.games[gameID] = game

	return nil
}
