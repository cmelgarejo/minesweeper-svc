package repo

import (
	"context"

	"github.com/cmelgarejo/minesweeper-svc/database"
	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"gorm.io/gorm/clause"
)

type GameRepo interface {
	UpsertGame(ctx context.Context, input *models.Game) (*models.Game, error)
	Read(ctx context.Context, gameID string) (game *models.Game, err error)
}

type GameRepoSvc struct {
	db *database.DB
}

func NewGameRepoSvc(db *database.DB) GameRepo {
	return &GameRepoSvc{
		db: db,
	}
}

func (svc *GameRepoSvc) Read(ctx context.Context, gameID string) (*models.Game, error) {
	rec := &models.Game{BaseModel: models.BaseModel{ID: gameID}}
	err := svc.db.Model(rec).Preload(clause.Associations).First(rec).Error

	return rec, err
}

func (svc *GameRepoSvc) UpsertGame(ctx context.Context, input *models.Game) (*models.Game, error) {
	err := svc.db.Model(input).FirstOrCreate(input, input).Save(input).Error
	if err != nil {
		return nil, err
	}

	return input, nil
}
