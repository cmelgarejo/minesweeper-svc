package repo

import (
	"context"

	"github.com/cmelgarejo/minesweeper-svc/database"
	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"gorm.io/gorm/clause"
)

type GameRepo interface {
	UpsertGame(ctx context.Context, gameID *string, input *models.Game) (*models.Game, error)
	Read(ctx context.Context, gameID string) (game *models.Game, err error)
	List(ctx context.Context) (games []*models.Game, err error)
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

func (svc *GameRepoSvc) List(ctx context.Context) (games []*models.Game, err error) {
	err = svc.db.Model(games).Preload(clause.Associations).Find(&games).Error

	return
}

func (svc *GameRepoSvc) UpsertGame(ctx context.Context, gameID *string, input *models.Game) (*models.Game, error) {
	var err error
	if gameID == nil {
		err = svc.db.Model(input).FirstOrCreate(input, input).Error
	} else {
		input.ID = *gameID
		err = svc.db.Unscoped().Model(input).Save(input).Error
	}
	if err != nil {
		return nil, err
	}

	return input, nil
}
