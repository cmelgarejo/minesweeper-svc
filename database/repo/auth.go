package repo

import (
	"context"
	"errors"

	"github.com/cmelgarejo/minesweeper-svc/database"
	"github.com/cmelgarejo/minesweeper-svc/database/models"
	"github.com/cmelgarejo/minesweeper-svc/resources/messages/codes"
	"github.com/cmelgarejo/minesweeper-svc/utils"
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
	SignIn(ctx context.Context, username, password string) (userApiKey *models.UserAPIKey, err error)
	CheckUserAPIKey(ctx context.Context, apiKey string) (user *models.User, err error)
	Read(ctx context.Context, userID string) (*models.User, error)
	UpsertUser(ctx context.Context, userID *string, input *models.User) (*models.User, error)
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

func (svc *AuthRepoSvc) SignIn(ctx context.Context, username, password string) (userApiKey *models.UserAPIKey, err error) {
	user := &models.User{Username: username}
	err = svc.db.Model(user).Preload(clause.Associations).Preload("APIKeys").Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, svc.catalog.WrapErrorWithCtx(ctx, gorm.ErrRecordNotFound, codes.MsgCodeUnauthorized, TblUsers)
	} else if err != nil {
		return nil, svc.catalog.WrapErrorWithCtx(ctx, err, codes.MsgCodeDBUnexpectedErr, err.Error())
	}
	err = utils.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return nil, svc.catalog.GetErrorWithCtx(ctx, codes.MsgCodeUnauthorized, TblUsers)
	}

	return &user.APIKeys[0], nil
}

func (svc *AuthRepoSvc) CheckUserAPIKey(ctx context.Context, apiKey string) (user *models.User, err error) {
	uak := &models.UserAPIKey{APIKey: apiKey}
	err = svc.db.Model(uak).Preload(clause.Associations).Preload("User").Where("api_key = ?", apiKey).First(&uak).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, svc.catalog.WrapErrorWithCtx(ctx, gorm.ErrRecordNotFound, codes.MsgCodeUnauthorized, TblUserAPIKeys)
	} else if err != nil {
		return nil, svc.catalog.WrapErrorWithCtx(ctx, err, codes.MsgCodeDBUnexpectedErr, err.Error())
	}
	user = &uak.User

	return
}

func (svc *AuthRepoSvc) Read(ctx context.Context, userID string) (*models.User, error) {
	rec := &models.User{BaseModel: models.BaseModel{ID: userID}}
	err := svc.db.Model(rec).Preload(clause.Associations).First(rec).Error

	return rec, err
}

func (svc *AuthRepoSvc) UpsertUser(ctx context.Context, userID *string, input *models.User) (*models.User, error) {
	var err error
	if userID == nil {
		input.APIKeys = make([]models.UserAPIKey, 1)
		err = svc.db.Model(input).FirstOrCreate(input, input).Error
	} else {
		input.ID = *userID
		err = svc.db.Model(input).Save(input).Error
	}
	if err != nil {
		return nil, err
	}

	return input, nil
}
