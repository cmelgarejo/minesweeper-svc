package service

import (
	"github.com/cmelgarejo/minesweeper-svc/web/game/engine"
)

// MineSweeperGame represents a minesweeper game service
type MineSweeperGameSvc interface {
	CreateGame(rows, cols, mines int) (gameID string)
	StartGame(gameID string) (err error)
	Click(gameID string, clickType engine.ClickType, row, col int) (err error)
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

func (ms *MineSweeperGameSvcImpl) CreateGame(rows, cols, mines int) (gameID string) {
	game := engine.NewGame(rows, cols, mines)
	ms.games[game.ID] = game
	return game.ID
}

func (ms *MineSweeperGameSvcImpl) StartGame(gameID string) (err error) {
	game := ms.games[gameID]
	return game.Start()
}

func (ms *MineSweeperGameSvcImpl) Click(gameID string, clickType engine.ClickType, col, row int) (err error) {
	game := ms.games[gameID]
	err = game.Click(clickType, row, col)
	return err
}
