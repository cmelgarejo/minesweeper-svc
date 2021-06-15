package service

import (
	"github.com/cmelgarejo/minesweeper-svc/web/game/engine"
)

// MineSweeperGame represents a minesweeper game service
type MineSweeperGameSvc interface {
	CreateGame(rows, cols, mines int) (game *engine.Game, err error)
	StartGame(gameID string) (err error)
	GetGame(gameID string) (game *engine.Game, err error)
	Click(gameID string, user string, clickType engine.ClickType, row, col int) (err error)
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

func (ms *MineSweeperGameSvcImpl) CreateGame(rows, cols, mines int) (game *engine.Game, err error) {
	game = engine.NewGame(rows, cols, mines)
	ms.games[game.ID] = game
	return ms.games[game.ID], err
}

func (ms *MineSweeperGameSvcImpl) StartGame(gameID string) (err error) {
	game := ms.games[gameID]
	return game.Start()
}

func (ms *MineSweeperGameSvcImpl) GetGame(gameID string) (game *engine.Game, err error) {
	return ms.games[gameID], nil
}

func (ms *MineSweeperGameSvcImpl) Click(gameID string, clickedBy string, clickType engine.ClickType, col, row int) (err error) {
	game := ms.games[gameID]
	err = game.Click(clickedBy,clickType, row, col)
	return err
}
