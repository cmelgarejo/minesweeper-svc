package migrations

import (
	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func firstUserMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "FIRST_USERS",
		Migrate: func(tx *gorm.DB) (err error) {
			adm := &models.User{
				Fullname: "Minesweeper Admin",
				Email:    "admin@minesweeper.svc",
				Username: "admin",
				Password: "test7eed0b835b8659dfb76b7956ce82ba6c",
				Admin:    true,
				APIKeys: []models.UserAPIKey{
					{}, // Forces to create a new API key
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

			return
		},
		Rollback: func(tx *gorm.DB) (err error) {
			return
		},
	}
}
