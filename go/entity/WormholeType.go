package entity

type WormholeType struct {
	Type          string `gorm:"type:varchar(255);not null;primaryKey"`
	TotalMass     int    `gorm:"not null"`
	MaxMass       int    `gorm:"not null"`
	MaxStableTime int    `gorm:"not null"`
}

func (WormholeType) TableName() string {
	return "wormhole_type"
}
