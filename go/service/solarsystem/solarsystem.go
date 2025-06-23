package solarsystem

import (
	"eve-wormhole-backend/go/dao"
	"eve-wormhole-backend/go/entity"

	"gorm.io/gorm"
)

func GetSolarSystemInfo(solarSystemID int) (*entity.MapSolarSystem, error) {
	solarSystem := &entity.MapSolarSystem{}
	if err := dao.SqlSession.Where("solarSystemID = ?", solarSystemID).
		First(solarSystem).Error; err != nil {
		return nil, err
	}
	return solarSystem, nil
}

func CheckSolarSystemJump(fromSolarSystemID, toSolarSystemID int) (bool, error) {
	var jump entity.MapSolarSystemJump
	if err := dao.SqlSession.Where("fromSolarSystemID = ? AND toSolarSystemID = ?", fromSolarSystemID, toSolarSystemID).
		First(&jump).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // No jump found
		}
		return false, err // Other error
	}
	return true, nil // Jump exists
}
