package service

import (
	"math/rand"
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
	mineCount := 0
	mineField := make([][]Field, rows)
	for i := 0; i < rows; i++ {
		mineField[i] = make([]Field, cols)
		for j := 0; j < cols; j++ {
			mineField[i][j] = Field{}
		}
	}
	for mineCount < mines {
		seed := rand.NewSource(time.Now().UnixNano())
		row := rand.New(seed).Intn(rows)
		col := rand.New(seed).Intn(cols)
		if !mineField[row][col].Mine {
			mineField[row][col].Mine = true
			mineCount++
		}
	}
	newGame := Game{
		ID:        id,
		Rows:      rows,
		Cols:      cols,
		Mines:     mines,
		Status:    GameStatusCreated,
		MineField: mineField,
		CreatedAt: time.Now(),
	}
	ms.games[id] = &newGame
	return newGame.ID
}

func (ms *MineSweeperGameSvcImpl) StartGame(gameID uuid.UUID) (err error) {
	game := ms.games[gameID]
	return game.Start()
}

func (ms *MineSweeperGameSvcImpl) Click(gameID uuid.UUID, clickType ClickType, col, row int) (err error) {
	game := ms.games[gameID]
	return game.Click(clickType, row, col)
}
