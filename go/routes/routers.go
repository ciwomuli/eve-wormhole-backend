package routes

import (
	"eve-wormhole-backend/go/controller/auth"
	"eve-wormhole-backend/go/controller/esi"
	"eve-wormhole-backend/go/controller/user"
	"eve-wormhole-backend/go/controller/wormhole"
	"eve-wormhole-backend/go/service/middleware"

	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	esiGroup := r.Group("esi")
	{
		esiGroup.GET("/auth", esi.Auth)
		esiGroup.GET("/callback", esi.Callback)
	}
	authGroup := r.Group("auth")
	{
		authGroup.POST("/submit-code", auth.SubmitCode)
	}
	userGroup := r.Group("user")
	{
		userGroup.GET("/info", user.Info)
	}
	wormholeGroup := r.Group("wormhole")
	{
		wormholeGroup.GET("/ws", wormhole.HandleWebSocket)
	}
	return r
}
