package wormhole

import (
	"encoding/json"
	"eve-wormhole-backend/go/service/wormhole"
	"eve-wormhole-backend/go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有CORS请求，生产环境应当配置更严格的规则
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c *gin.Context) {
	logrus.Info("WebSocket connection request received")
	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Fail(c, "Failed to upgrade connection to WebSocket")
		return
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		utils.Fail(c, "Failed to read message from WebSocket")
		return
	}

	var initialMessage map[string]interface{}
	err = json.Unmarshal(message, &initialMessage)
	if err != nil {
		utils.Fail(c, "Invalid JSON format in initial message")
		return
	}

	err = utils.DecodeJWT(c, initialMessage["Authorization"].(string))
	if err != nil {
		utils.Fail(c, "Invalid JWT token")
		return
	}

	defer conn.Close()

	wormhole.StartWormholeTicker(c, conn)
	logrus.Info("WebSocket connection ended")
}
