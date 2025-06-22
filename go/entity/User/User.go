package entity

import "gorm.io/gorm"

func (User) TableName() string {
	return "sys_user"
}

type User struct {
	gorm.Model
	Name     string        `gorm:"type:varchar(200)"`
	Accounts []UserAccount `gorm:"foreignKey:UserID;references:ID"`
}
