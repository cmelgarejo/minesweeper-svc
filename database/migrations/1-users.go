package migrations

import (
	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/cmelgarejo/minesweeper-svc/web/game/engine"
	"github.com/cmelgarejo/minesweeper-svc/web/game/service"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const (
	AdminApiKey = "587fa65a9c375165828a6fbb5f9963a7"
	TestGameID  = "ef99fdfd88565827ad330d83aac5fbaa"
)

func firstUserMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "FIRST_USERS",
		Migrate: func(tx *gorm.DB) (err error) {
			adm := &models.User{
				Fullname: "Minesweeper Admin",
				Email:    "admin@minesweeper.svc",
				Username: "mineadmin",
				Password: AdminApiKey,
				Admin:    true,
				APIKeys: []models.UserAPIKey{
					{
						APIKey: AdminApiKey,
					},
				},
			}
			if err = tx.Create(adm).Error; err != nil {
				return err
			}
			p1 := &models.User{
				Fullname: "Player One",
				Email:    "player.one@minesweeper.svc",
				Username: "player1",
				Password: "player1",
				APIKeys: []models.UserAPIKey{
					{}, // Forces to create a new API key
				},
			}
			if err = tx.Create(p1).Error; err != nil {
				return err
			}
			p2 := &models.User{
				Fullname: "Player Two",
				Email:    "player.two@minesweeper.svc",
				Username: "player2",
				Password: "player2",
				APIKeys: []models.UserAPIKey{
					{}, // Forces to create a new API key
				},
			}
			if err = tx.Create(p2).Error; err != nil {
				return err
			}
			gameService := service.MineSweeperGameSvcImpl{}
			game, _ := gameService.NewMineSweeperSvc().CreateGame(5, 5, 3, adm.Fullname)
			game.ID = TestGameID
			testGame := &models.Game{
				Rows: 5, Cols: 5, Mines: 3,
				Status:      engine.GameStatusCreated,
				CreatedByID: adm.ID,
			}
			testGame.ID = TestGameID
			testGame.UpdateGameState(game)
			if err = tx.Create(testGame).Error; err != nil {
				return err
			}

			return
		},
		Rollback: func(tx *gorm.DB) (err error) {
			return
		},
	}
}
