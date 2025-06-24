package wormhole

import (
	"context"
	"errors"
	"eve-wormhole-backend/go/dao"
	"eve-wormhole-backend/go/entity"
	"eve-wormhole-backend/go/service/esi"
	"eve-wormhole-backend/go/service/solarsystem"
	"eve-wormhole-backend/go/service/user"
	"math"
	"regexp"
	"time"

	"github.com/antihax/goesi"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func StartWormholeTicker(c *gin.Context, conn *websocket.Conn) {
	// 创建一个定时器，每 10 秒触发一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop() // 确保定时器在退出时停止
	userId, exist := c.Get("userId")
	if !exist {
		logrus.Error("User ID not found in context")
		return
	}
	//创建一个字典记录每个账号的位置
	var accountLocations = make(map[int32]int)

	for now := range ticker.C {
		//发送一个心跳包 如果发送失败就关闭定时器
		if err := conn.WriteJSON(gin.H{"type": "heartbeat", "time": now.Unix()}); err != nil {
			return
		}
		currentUser, err := user.GetUserByIDWithAccounts(userId.(uint))
		if err != nil {
			logrus.Errorf("Failed to get user by ID %d: %v", userId, err)
			continue
		}
		for _, account := range currentUser.Accounts {
			token := user.Token(&account)
			ctx := context.WithValue(context.Background(), goesi.ContextOAuth2, token)
			location, response, err := esi.E.ESI.ESI.LocationApi.GetCharactersCharacterIdLocation(ctx, account.CharacterID, nil)
			if err != nil || response.StatusCode != 200 {
				logrus.Errorf("Failed to check online status for character %d: %v, response: %v", account.CharacterID, err, response)
				continue
			}
			if accountLocations[account.CharacterID] != 0 && int(location.SolarSystemId) != accountLocations[account.CharacterID] {
				fromSystem, err := solarsystem.GetSolarSystemInfo(accountLocations[account.CharacterID])
				if err != nil {
					logrus.Errorf("Failed to get solar system info for character %d: %v", account.CharacterID, err)
					continue
				}
				toSystem, err := solarsystem.GetSolarSystemInfo(int(location.SolarSystemId))
				if err != nil {
					logrus.Errorf("Failed to get solar system info for character %d: %v", account.CharacterID, err)
					continue
				}
				// 检查是否存在跳跃
				jumpExists, err := solarsystem.CheckSolarSystemJump(fromSystem.SolarSystemID, toSystem.SolarSystemID)
				if err != nil {
					logrus.Errorf("Failed to check solar system jump for character %d: %v", account.CharacterID, err)
					continue
				}
				if !jumpExists {
					err := conn.WriteJSON(
						gin.H{
							"type":           "wormhole",
							"fromSystemId":   accountLocations[account.CharacterID],
							"fromSystemName": fromSystem.SolarSystemName,
							"toSystemId":     int(location.SolarSystemId),
							"toSystemName":   toSystem.SolarSystemName},
					)
					if err != nil {
						break
					}
				}
			}
			accountLocations[account.CharacterID] = int(location.SolarSystemId)
		}
	}
}

func GetWormholeType(wormholeType string) (*entity.WormholeType, error) {
	wormholeTypeEntity := &entity.WormholeType{}
	err := dao.SqlSession.Where("type = ?", wormholeType).First(wormholeTypeEntity).Error
	if err != nil {
		logrus.Errorf("Failed to get wormhole type %s: %v", wormholeType, err)
		return nil, err
	}
	return wormholeTypeEntity, nil
}

type WormholeConnectionData struct {
	Type         string `json:"type"`
	FromSystemId int    `json:"from_system_id"`
	FromSignal   string `json:"from_signal"`
	ToSystemId   int    `json:"to_system_id"`
	ToSignal     string `json:"to_signal"`
	MassStatus   int    `json:"mass_status"`
	TimeStatus   int    `json:"time_status"`
}

func AddWormholeConnection(userId uint, wormholeConnectionData *WormholeConnectionData) error {
	_, err := GetWormholeType(wormholeConnectionData.Type)
	if err != nil {
		logrus.Errorf("Failed to get wormhole type %s: %v", wormholeConnectionData.Type, err)
		return err
	}
	_, err = solarsystem.GetSolarSystemInfo(wormholeConnectionData.FromSystemId)
	if err != nil {
		logrus.Errorf("Failed to get solar system info for FromSystemId %d: %v", wormholeConnectionData.FromSystemId, err)
		return err
	}
	_, err = solarsystem.GetSolarSystemInfo(wormholeConnectionData.ToSystemId)
	if err != nil {
		logrus.Errorf("Failed to get solar system info for ToSystemId %d: %v", wormholeConnectionData.ToSystemId, err)
		return err
	}

	//使用正则表达式检查Signal是否是由三个大写字母-三个数字构成的
	signal_pattern := `^[A-Z]{3}-\d{3}$`
	re, _ := regexp.Compile(signal_pattern)
	if !re.MatchString(wormholeConnectionData.FromSignal) || !re.MatchString(wormholeConnectionData.ToSignal) {
		logrus.Errorf("Invalid signal format: FromSignal %s, ToSignal %s", wormholeConnectionData.FromSignal, wormholeConnectionData.ToSignal)
		return errors.New("invalid signal format")
	}

	wormConnection := &entity.WormholeConnection{
		UserId:            userId,
		Type:              wormholeConnectionData.Type,
		FromSolarSystemID: wormholeConnectionData.FromSystemId,
		FromSignal:        wormholeConnectionData.FromSignal,
		ToSolarSystemID:   wormholeConnectionData.ToSystemId,
		ToSignal:          wormholeConnectionData.ToSignal,
		MassStatus:        wormholeConnectionData.MassStatus,
		TimeStatus:        wormholeConnectionData.TimeStatus,
		StatusUpdateTime:  time.Now(),
		UseCount:          0,
		Alive:             true,
	}
	err = dao.SqlSession.Create(wormConnection).Error

	if err != nil {
		logrus.Errorf("Failed to add wormhole connection: %v", err)
		return err
	}

	return nil
}

var WormholeTimeStatusMin = []int{1440, 240, 0}
var WormholeTimeStatusMax = []int{2880, 1440, 240}

var WormholeMassStatusMin = []float32{0.5, 0.1, 0}
var WormholeMassStatusMax = []float32{1.0, 0.5, 0.1}

// 返回一个json数组
func ListWormholeConnections(c *gin.Context, all bool) ([]*gin.H, error) {
	userId, exist := c.Get("userId")
	if !exist {
		logrus.Error("User ID not found in context")
		return nil, errors.New("user ID not found")
	}

	var wormholeConnections []*entity.WormholeConnection

	query := dao.SqlSession
	if !all {
		query = query.Where("user_id = ?", userId.(uint))
	}

	err := query.Find(&wormholeConnections).Error
	if err != nil {
		logrus.Errorf("Failed to list wormhole connections for user %d: %v", userId, err)
		return nil, err
	}
	ret := []*gin.H{}
	for _, connection := range wormholeConnections {
		wormholeType, err := GetWormholeType(connection.Type)
		if err != nil {
			logrus.Errorf("Failed to get wormhole type for connection %d: %v", connection.ID, err)
			continue
		}
		fromSystem, err := solarsystem.GetSolarSystemInfo(connection.FromSolarSystemID)
		if err != nil {
			logrus.Errorf("Failed to get from solar system info for connection %d: %v", connection.ID, err)
			continue
		}
		toSystem, err := solarsystem.GetSolarSystemInfo(connection.ToSolarSystemID)
		if err != nil {
			logrus.Errorf("Failed to get to solar system info for connection %d: %v", connection.ID, err)
			continue
		}
		if connection.TimeStatus < 0 || connection.TimeStatus >= len(WormholeTimeStatusMin) {
			logrus.Errorf("Invalid time status %d for connection %d", connection.TimeStatus, connection.ID)
			continue
		}
		if connection.MassStatus < 0 || connection.MassStatus >= len(WormholeMassStatusMin) {
			logrus.Errorf("Invalid mass status %d for connection %d", connection.MassStatus, connection.ID)
			continue
		}
		if connection.StatusUpdateTime.Add(time.Duration(math.Min(float64(WormholeTimeStatusMax[connection.TimeStatus]), float64(wormholeType.MaxStableTime))) * time.Minute).Before(time.Now()) {
			connection.Alive = false
			SaveWormholeConnection(connection)
		}

		ret = append(ret, &gin.H{
			"type":               connection.Type,
			"from_system_id":     connection.FromSolarSystemID,
			"from_system_name":   fromSystem.SolarSystemName,
			"from_signal":        connection.FromSignal,
			"to_system_id":       connection.ToSolarSystemID,
			"to_system_name":     toSystem.SolarSystemName,
			"to_signal":          connection.ToSignal,
			"mass_status":        connection.MassStatus,
			"time_status":        connection.TimeStatus,
			"status_update_time": connection.StatusUpdateTime.Unix(),
			"stable_time_min":    connection.StatusUpdateTime.Add(time.Duration(WormholeTimeStatusMin[connection.TimeStatus]) * time.Minute).Unix(),
			"stable_time_max":    connection.StatusUpdateTime.Add(time.Duration(math.Min(float64(WormholeTimeStatusMax[connection.TimeStatus]), float64(wormholeType.MaxStableTime))) * time.Minute).Unix(),
			"mass_min":           float32(wormholeType.TotalMass) * WormholeMassStatusMin[connection.MassStatus],
			"mass_max":           float32(wormholeType.TotalMass) * WormholeMassStatusMax[connection.MassStatus],
		})
	}

	return ret, nil
}

func SaveWormholeConnection(connection *entity.WormholeConnection) error {
	// 更新连接状态
	err := dao.SqlSession.Save(connection).Error
	if err != nil {
		logrus.Errorf("Failed to update wormhole connection status for connection %d: %v", connection.ID, err)
		return err
	}
	return nil
}
