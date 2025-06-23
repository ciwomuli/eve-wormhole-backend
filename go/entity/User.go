package entity

import (
	"time"

	"gorm.io/gorm"
)

func (User) TableName() string {
	return "sys_user"
}

type User struct {
	gorm.Model
	Name       string `gorm:"type:varchar(200)"`
	ActiveTime time.Time
	Accounts   []UserAccount `gorm:"foreignKey:UserID;references:ID"`
}
