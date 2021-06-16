package responses

import "time"

type ResponseBase struct {
	Code    int    `json:"code"  example:"12345"`
	Message string `json:"message"  example:"message"`
	Details string `json:"details"  example:"some more details"`
}

type Response struct {
	ResponseBase
	Error    []ResponseError `json:"errors" example:""`
	Response interface{}     `json:"response"  example:""`
}

type ResponseError struct {
	ResponseBase
}

//Field represents a square unit in the MineField
type Field struct {
	Mine      bool           `json:"mine,omitempty"`
	Clicked   bool           `json:"clicked"`   // indicated whether the field was clicked
	Flagged   bool           `json:"flagged"`   // red flag in the field
	AdjCount  int            `json:"adjMines"`  // count of adjacent mines
	Position  map[string]int `json:"position"`  // position in the minefield
	ClickedBy string         `json:"clickedBy"` // who clicked this field
}

// Game contains the structure of the game
type Game struct {
	ID         string     `json:"id"`
	Rows       int        `json:"rows"`
	Cols       int        `json:"cols"`
	Mines      int        `json:"mines"`
	Status     string     `json:"status"`
	MineField  [][]Field  `json:"mineField"`
	StartedAt  *time.Time `json:"startedAt,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	CreatedBy  string     `json:"createdBy"` // who created this game
}
