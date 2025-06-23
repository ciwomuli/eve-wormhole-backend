package entity

// MapSolarSystemJump 对应 mapSolarSystemJumps 表
// 表示EVE游戏中星系之间的跳跃连接
type MapSolarSystemJump struct {
	FromRegionID        int `gorm:"column:fromRegionID"`
	FromConstellationID int `gorm:"column:fromConstellationID"`
	FromSolarSystemID   int `gorm:"column:fromSolarSystemID;primaryKey;autoIncrement:false"`
	ToSolarSystemID     int `gorm:"column:toSolarSystemID;primaryKey;autoIncrement:false"`
	ToConstellationID   int `gorm:"column:toConstellationID"`
	ToRegionID          int `gorm:"column:toRegionID"`
}

// TableName 指定表名
func (MapSolarSystemJump) TableName() string {
	return "mapSolarSystemJumps"
}
