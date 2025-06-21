package entity

import "gorm.io/gorm"

func (UserAccount) TableName() string {
	return "user_account"
}

type UserAccount struct {
	gorm.Model
	UserID uint
}
