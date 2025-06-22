package entity

import (
	"time"

	"gorm.io/gorm"
)

func (UserAccount) TableName() string {
	return "user_account"
}

type UserAccount struct {
	gorm.Model
	UserID        uint
	CharacterID   int32
	CharacterName string `gorm:"type:varchar(200)"`
	AccessToken   string `gorm:"type:text"`
	RefreshToken  string `gorm:"type:varchar(200)"`
	Expiry        time.Time
}
