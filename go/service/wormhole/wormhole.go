package wormhole

import (
	"context"
	"eve-wormhole-backend/go/service/esi"
	"eve-wormhole-backend/go/service/solarsystem"
	"eve-wormhole-backend/go/service/user"
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
