package routes

import (
	"eve-wormhole-backend/go/controller/ESI"
	"eve-wormhole-backend/go/service/middleware"

	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	esiGroup := r.Group("esi")
	{
		esiGroup.GET("/auth", ESI.Auth)
		esiGroup.GET("/callback", ESI.Callback)
	}
	return r
}
