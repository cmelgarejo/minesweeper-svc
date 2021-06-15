package repo

//go:generate mockgen -source=$GOFILE -package repo_test -destination=../../test/mock/repo/$GOFILE

import (
	"context"
	"errors"

	"github.com/cmelgarejo/minesweeper-svc/database"
	"github.com/cmelgarejo/minesweeper-svc/database/models"
	msgcodes "github.com/cmelgarejo/minesweeper-svc/resources/messages/codes"
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	"github.com/loopcontext/msgcat"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	TblUsers       = "users"
	TblUserAPIKeys = "user_api_keys"
)

type AuthRepo interface {
	CheckUserAPIKey(ctx context.Context, apiKey string) (user *models.User, err error)
}

type AuthRepoSvc struct {
	db      *database.DB
	log     logger.Logger
	catalog msgcat.MessageCatalog
}

func NewAuthRepoSvc(db *database.DB, log logger.Logger, catalog msgcat.MessageCatalog) AuthRepo {
	return &AuthRepoSvc{
		db:      db,
		log:     log,
		catalog: catalog,
	}
}

func (svc *AuthRepoSvc) CheckUserAPIKey(ctx context.Context, apiKey string) (user *models.User, err error) {
	uak := &models.UserAPIKey{APIKey: apiKey}
	err = svc.db.Model(uak).Preload(clause.Associations).Preload("User").Where("api_key = ?", apiKey).First(&uak).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, svc.catalog.WrapErrorWithCtx(ctx, gorm.ErrRecordNotFound, msgcodes.MsgCodeUnauthorized, TblUserAPIKeys)
	} else if err != nil {
		return nil, svc.catalog.WrapErrorWithCtx(ctx, err, msgcodes.MsgCodeDBUnexpectedErr, err.Error())
	}
	user = &uak.User

	return
}
