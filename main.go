package main

import (
	"fmt"

	"github.com/cmelgarejo/minesweeper-svc/service"
)

func main() {
	var err error
	svc := service.MineSweeperGameSvcImpl{}
	gameSvc := svc.NewMineSweeperSvc()
	gameID := gameSvc.CreateGame(0, 0, 0)
	// fmt.Printf("gameId: %v\n", gameID)
	err = gameSvc.StartGame(gameID)
	if err != nil {
		fmt.Printf("%q", err)
	}
	err = gameSvc.Click(gameID, service.GameClickTypeNormal, 1, 1)
	if err != nil {
		fmt.Printf("%q", err)
	}
}
