package models

import (
	"time"

	"github.com/cmelgarejo/minesweeper-svc/utils"
	"github.com/cmelgarejo/minesweeper-svc/web/game/engine"
)

// Game contains the structure of the game
type Game struct {
	BaseModel
	Rows        int        `json:"rows"`
	Cols        int        `json:"cols"`
	Mines       int        `json:"mines"`
	Status      string     `json:"status"`
	GameState   JSONB      `json:"gameState" gorm:"type:jsonb"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`
	CreatedByID string     `json:"-"`         // who created this game - id needed by GORM
	CreatedBy   *User      `json:"createdBy"` // who created this game
}

func (g *Game) UpdateGameState(game *engine.Game) {
	b, _ := utils.ToJSONBytes(game)
	gameState := JSONB{}
	_ = utils.ToObject(b, &gameState)
	g.GameState = gameState
}

func (g *Game) GetGameState() (game *engine.Game) {
	b, _ := utils.ToJSONBytes(g.GameState)
	_ = utils.ToObject(b, &game)

	return
}
