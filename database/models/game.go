package models

import (
	"time"
)

// Game contains the structure of the game
type Game struct {
	BaseModel
	Rows        int        `json:"rows"`
	Cols        int        `json:"cols"`
	Mines       int        `json:"mines"`
	Status      int        `json:"status"`
	MineField   JSONB      `json:"mineField" gorm:"type:jsonb"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`
	CreatedByID string     `json:"-"`         // who created this game - id needed by GORM
	CreatedBy   *User      `json:"createdBy"` // who created this game
	// Players    []User     `json:"players"`
}
