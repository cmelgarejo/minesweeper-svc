package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrGameOver       = errors.New("Game Over")
	ErrGameNotActive  = errors.New("Game is not active")
	ErrAlreadyClicked = errors.New("Field already clicked")
)

// Some default game parameters, if the user does not provide those.
const (
	GameMinRows         = 3
	GameMinCols         = 3
	GameMaxRows         = 50
	GameMaxCols         = 50
	GameStatusCreated   = "created"
	GameStatusStarted   = "started"
	GameStatusOver      = "over"
	GameStatusWin       = "win"
	GameClickTypeNormal = 1
	GameClickTypeFlag   = 2
)

type GameStatus string
type ClickType int

//Field represents a square unit in the MineField
type Field struct {
	Mine     bool `json:"-"`
	Clicked  bool `json:"clicked"`  // indicated whether the field was clicked
	Flagged  bool `json:"flagged"`  // red flag in the field
	AdjCount int  `json:"adjMines"` // count of adjacent mines
}

// ClickAction represents an action on a given field
type ClickAction struct {
	Row       int       `json:"row"`
	Col       int       `json:"col"`
	ClickType ClickType `json:"clickType"`
}

// Game contains the structure of the game
type Game struct {
	ID         uuid.UUID  `json:"id"`
	Rows       int        `json:"rows"`
	Cols       int        `json:"cols"`
	Mines      int        `json:"mines"`
	Status     GameStatus `json:"status"`
	MineField  [][]Field  `json:"mineField"`
	CreatedAt  time.Time  `json:"createdAt"`
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

func (g *Game) Start() error {
	if g.Status == GameStatusCreated {
		g.Status = GameStatusStarted
		return nil
	}
	return fmt.Errorf("Cannot start a game that is in status: %s", g.Status)
}

func (g *Game) Click(clickType ClickType, row, col int) error {
	if !g.IsActive() {
		return ErrGameNotActive
	}
	if row > g.Rows || col > g.Cols {
		return fmt.Errorf("Field [%d, %d] out of bounds", row, col)
	}
	if !g.MineField[row][col].Clicked {
		g.MineField[row][col].Clicked = true
		switch clickType {
		case GameClickTypeFlag:
			g.MineField[row][col].Flagged = true
		case GameClickTypeNormal:
			if g.MineField[row][col].Mine {
				g.Status = GameStatusOver
				return ErrGameOver
			}
		}
	} else {
		return ErrAlreadyClicked
	}
	return nil
}

func (g *Game) GetStatusStr() string {
	switch g.Status {
	case GameStatusCreated:
		return "created"
	case GameStatusStarted:
		return "started"
	case GameStatusOver, GameStatusWin:
		return "finished"
	default:
		return ""
	}
}

func (g *Game) IsActive() bool {
	return g.Status == GameStatusStarted
}
