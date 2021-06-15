package migrations

import (
	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func firstUserMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "FIRST_REALM_AND_USER",
		Migrate: func(tx *gorm.DB) (err error) {
			u := &models.User{
				Fullname: "Admin",
				Email:    "admin@minesweeper.svc",
				Username: "admin",
				Password: "test7eed0b835b8659dfb76b7956ce82ba6c",
			}
			if err = tx.FirstOrCreate(u).Error; err != nil {
				return err
			}

			return
		},
		Rollback: func(tx *gorm.DB) (err error) {
			return
		},
	}
}
