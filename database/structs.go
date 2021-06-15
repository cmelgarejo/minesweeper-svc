package database

import (
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	"gorm.io/gorm"
)

type DB struct {
	Logger *logger.Logger
	*gorm.DB
}

type DBConfig struct {
	Automigrate bool
	Debug       bool
	DSN         string
}
