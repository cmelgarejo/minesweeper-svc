package service

import (
	"time"

	"github.com/google/uuid"
)

// MineSweeperGame represents a minesweeper game service
type MineSweeperGameSvc interface {
	CreateGame(rows, cols, mines int) (gameID uuid.UUID)
	StartGame(gameID uuid.UUID) (err error)
	Click(gameID uuid.UUID, clickType ClickType, row, col int) (err error)
}

// MineSweeperGameSvcImpl implementing struct of a minesweeper game service
type MineSweeperGameSvcImpl struct {
	// Will hold persistance instance, logger, and/or any instances needed by the game
	games map[uuid.UUID]*Game
}

func (ms *MineSweeperGameSvcImpl) NewMineSweeperSvc() MineSweeperGameSvc {
	return &MineSweeperGameSvcImpl{
		games: make(map[uuid.UUID]*Game),
	}
}

func (ms *MineSweeperGameSvcImpl) CreateGame(rows, cols, mines int) (gameID uuid.UUID) {
	if rows < GameMinRows {
		rows = GameMinRows
	}
	if cols < GameMinCols {
		cols = GameMinCols
	}
	if mines < 1 || mines > rows*cols {
		mines = rows + cols // Make sure amount of mines is relative to a median of rows + cols
	}
	id, _ := uuid.NewUUID()
	newGame := Game{
		ID:        id,
		Rows:      rows,
		Cols:      cols,
		Mines:     mines,
		CreatedAt: time.Now(),
	}
	newGame.InitializeMinefield()

	ms.games[id] = &newGame
	return newGame.ID
}

func (ms *MineSweeperGameSvcImpl) StartGame(gameID uuid.UUID) (err error) {
	game := ms.games[gameID]
	return game.Start()
}

func (ms *MineSweeperGameSvcImpl) Click(gameID uuid.UUID, clickType ClickType, col, row int) (err error) {
	game := ms.games[gameID]
	err = game.Click(clickType, row, col)
	return err
}
