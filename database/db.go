package database

import (
	"github.com/cmelgarejo/minesweeper-svc/database/migrations"
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg DBConfig, log *logger.Logger) (*DB, error) {
	if cfg.DSN == "" {

	}
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.SendError(err)
	}
	if cfg.Automigrate {
		if err = migrations.RunMigrations(db); err != nil {
			return nil, err
		}
	}
	if cfg.Debug {
		db = db.Debug()
	}
	return &DB{
		Logger: log,
		DB:     db,
	}, nil
}
