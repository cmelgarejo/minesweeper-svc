package migrations

import (
	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func initialMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "INITIAL_MIGRATION",
		Migrate: func(tx *gorm.DB) (err error) {
			return tx.AutoMigrate(
				models.User{},
				models.UserAPIKey{},
				models.Game{},
			)
		},
		Rollback: func(tx *gorm.DB) (err error) {
			if err := tx.Migrator().DropTable(tx.Model(&models.User{}).Name()); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(tx.Model(&models.UserAPIKey{}).Name()); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(tx.Model(&models.Game{}).Name()); err != nil {
				return err
			}
			return
		},
	}
}
