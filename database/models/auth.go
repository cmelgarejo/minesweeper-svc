package models

import (
	"time"

	"github.com/cmelgarejo/minesweeper-svc/utils"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Email     string       `json:"email" gorm:"uniqueIndex"`
	Username  string       `json:"username" gorm:"uniqueIndex"`
	Password  string       `json:"password"`
	Fullname  string       `json:"fullname"`
	APIKeys   []UserAPIKey `json:"apiKeys"`
	Admin     bool         `json:"-" gorm:"default:false"` // FUTURE: proper RBAC, roles, permissions, etc.
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
}

// BeforeCreate hook that runs before entity create
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	err = u.BaseModel.BeforeCreate(tx)
	if err != nil {
		return err
	}

	return u.hashPassword()
}

// BeforeUpdate hook to run before entity update
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	return u.hashPassword()
}

func (u *User) hashPassword() (err error) {
	if u.Password != "" {
		hash, err := utils.EncryptPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hash
	}

	return
}

type UserAPIKey struct {
	ID        uint      `json:"-" gorm:"primarykey"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	APIKey    string    `gorm:"unique" json:"apiKey"`
	UserID    string    `json:"-"`
	User      User      `json:"-"`
}

// BeforeCreate hook to run before entity create
func (uak *UserAPIKey) BeforeCreate(tx *gorm.DB) (err error) {
	uak.APIKey, err = utils.GenerateGUID()
	return
}
