package models

import (
	"database/sql/driver"
	"encoding/json"
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

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}
