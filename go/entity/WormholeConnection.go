package entity

import (
	"time"

	"gorm.io/gorm"
)

type WormholeConnection struct {
	gorm.Model
	UserId            uint      `gorm:"not null"`
	Type              string    `gorm:"type:varchar(255);not null"`
	FromSolarSystemID int       `gorm:"not null"`
	FromSignal        string    `gorm:"type:varchar(255);not null"`
	ToSolarSystemID   int       `gorm:"not null"`
	ToSignal          string    `gorm:"type:varchar(255);not null"`
	MassStatus        int       `gorm:"not null"`
	TimeStatus        int       `gorm:"not null"`
	StatusUpdateTime  time.Time `gorm:"not null"`
	UseCount          int       `gorm:"not null"`
	Alive             bool      `gorm:"not null"`
}

func (WormholeConnection) TableName() string {
	return "wormhole_connection"
}
