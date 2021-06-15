package models

import (
	"time"
)

// Game contains the structure of the game
type Game struct {
	BaseModel
	Rows       int        `json:"rows"`
	Cols       int        `json:"cols"`
	Mines      int        `json:"mines"`
	Status     int        `json:"status"`
	MineField  [][]Field  `json:"mineField"`
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
	Players    []User     `json:"players"`
	CreatedBy  User       `json:"createdBy"` // who created this game
}

//Field represents a square unit in the MineField
type Field struct {
	Mine      bool     `json:"-"`
	Clicked   bool     `json:"clicked"`   // indicated whether the field was clicked
	Flagged   bool     `json:"flagged"`   // red flag in the field
	AdjCount  int      `json:"adjMines"`  // count of adjacent mines
	Position  Position `json:"position"`  // position in the minefield
	ClickedBy User     `json:"clickedBy"` // who clicked this field
}

// Position stores the position of the field in the board
type Position struct {
	Row int `json:"row"` // row of the field position
	Col int `json:"col"` // col of the field position
}
