package models

import (
	"time"

	"github.com/cmelgarejo/minesweeper-svc/utils"
	"gorm.io/gorm"
)

// BaseModel a basic Go struct which includes the following fields: ID (GUID), CreatedAt, UpdatedAt, DeletedAt
type BaseModel struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID, err = utils.GenerateGUID()

	return
}
