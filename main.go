package main

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/cmelgarejo/minesweeper-svc/docs"
	"github.com/cmelgarejo/minesweeper-svc/service"
	"github.com/gofiber/fiber/v2"
)

// @title Minesweeper API
// @version 1.0.8
// @description A minesweeper game API
// @contact.name Christian Melgarejo
// @contact.email cmelgarejo.dev@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()

	app.Post("/create", CreateGame)

	app.Use("/swagger", swagger.Handler) // swagger

	_ = app.Listen(":3000")
}

// CreateGame godoc
// @Summary Creates a game of minesweeper
// @Description Creates a game of minesweeper and returns a gameID that yo uahve to store to keep tabs on it
// @Tags game
// @Accept json
// @Produce json
// @Success 201 {object} service.Game
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /create [post]
func CreateGame(c *fiber.Ctx) error {
	svc := service.MineSweeperGameSvcImpl{}
	gameSvc := svc.NewMineSweeperSvc()
	gameID := gameSvc.CreateGame(0, 0, 0)
	return c.SendString(gameID.String())
}

type HTTPError struct {
	status  string
	message string
}

// package main

// import (
// 	"fmt"

// 	"github.com/cmelgarejo/minesweeper-svc/service"
// )

// func main() {
// 	var err error
// 	svc := service.MineSweeperGameSvcImpl{}
// 	gameSvc := svc.NewMineSweeperSvc()
// 	gameID := gameSvc.CreateGame(5, 10, 5)
// 	// fmt.Printf("gameId: %s\n", gameID)
// 	err = gameSvc.StartGame(gameID)
// 	if err != nil {
// 		fmt.Printf("%q", err)
// 	}
// 	err = gameSvc.Click(gameID, service.GameClickTypeNormal, 1, 1)
// 	if err != nil {
// 		fmt.Printf("%q", err)
// 	}
// 	err = gameSvc.Click(gameID, service.GameClickTypeFlag, 1, 1)
// 	if err != nil {
// 		fmt.Printf("%q", err)
// 	}
// }
