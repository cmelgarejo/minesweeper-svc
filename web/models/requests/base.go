package requests

import "github.com/cmelgarejo/minesweeper-svc/web/game/engine"

type Credentials struct {
	Username string `json:"username" example:"player1"`
	Password string `json:"password" example:"player1"`
}

type UserInput struct {
	Credentials
	Email    string `json:"email" example:"player1@player1.com"`
	Fullname string `json:"fullname" example:"Ready Player One"`
}

type GameInput struct {
	Row       int    `json:"row" example:"0"`
	Col       int    `json:"col" example:"0"`
	ClickType string `json:"clickType" enums:"click,flag"`
}

type GameCreateInput struct {
	Rows  int `json:"row" example:"5"`
	Cols  int `json:"col" example:"5"`
	Mines int `json:"mines" example:"5"`
}

func (gi *GameInput) GetClickType() engine.ClickType {
	switch gi.ClickType {
	case "flag":
		return engine.GameClickTypeFlag
	case "normal":
		fallthrough
	default:
		return engine.GameClickTypeNormal
	}
}
