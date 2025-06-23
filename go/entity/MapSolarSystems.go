package entity

// MapSolarSystem 对应 mapSolarSystems 表
type MapSolarSystem struct {
	RegionID        int     `gorm:"column:regionID;index:ix_mapSolarSystems_regionID"`
	ConstellationID int     `gorm:"column:constellationID;index:ix_mapSolarSystems_constellationID"`
	SolarSystemID   int     `gorm:"column:solarSystemID;primaryKey"`
	SolarSystemName string  `gorm:"column:solarSystemName;type:varchar(100)"`
	X               float64 `gorm:"column:x"`
	Y               float64 `gorm:"column:y"`
	Z               float64 `gorm:"column:z"`
	XMin            float64 `gorm:"column:xMin"`
	XMax            float64 `gorm:"column:xMax"`
	YMin            float64 `gorm:"column:yMin"`
	YMax            float64 `gorm:"column:yMax"`
	ZMin            float64 `gorm:"column:zMin"`
	ZMax            float64 `gorm:"column:zMax"`
	Luminosity      float64 `gorm:"column:luminosity"`
	Border          bool    `gorm:"column:border;check:mapss_border,border in (0,1)"`
	Fringe          bool    `gorm:"column:fringe;check:mapss_fringe,fringe in (0,1)"`
	Corridor        bool    `gorm:"column:corridor;check:mapss_corridor,corridor in (0,1)"`
	Hub             bool    `gorm:"column:hub;check:mapss_hub,hub in (0,1)"`
	International   bool    `gorm:"column:international;check:mapss_internat,international in (0,1)"`
	Regional        bool    `gorm:"column:regional;check:mapss_regional,regional in (0,1)"`
	Constellation   bool    `gorm:"column:constellation;check:mapss_constel,constellation in (0,1)"`
	Security        float64 `gorm:"column:security;index:ix_mapSolarSystems_security"`
	FactionID       int     `gorm:"column:factionID"`
	Radius          float64 `gorm:"column:radius"`
	SunTypeID       int     `gorm:"column:sunTypeID"`
	SecurityClass   string  `gorm:"column:securityClass;type:varchar(2)"`
}

// TableName 指定表名
func (MapSolarSystem) TableName() string {
	return "mapSolarSystems"
}
