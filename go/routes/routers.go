package routes

import (
	"eve-wormhole-backend/go/controller/Auth"
	"eve-wormhole-backend/go/controller/ESI"
	"eve-wormhole-backend/go/service/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()
	store, err := redis.NewStore(
		16,               // 最大空闲链接数量，过大会浪费，过小将来会触发性能瓶颈
		"tcp",            // 指定与Redis服务器通信的网络类型，通常为"tcp"
		"localhost:6379", // Redis服务器的地址，格式为"host:port"
		"",
		"", // Redis服务器的密码，如果没有密码可以为空字符串
		[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), // authentication key
		[]byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"), // encryption key
	)
	if err != nil {
		panic("Failed to create Redis store: " + err.Error())
	}

	r.Use(sessions.Sessions("test", store))
	r.Use(middleware.AuthMiddleware())
	esiGroup := r.Group("esi")
	{
		esiGroup.GET("/auth", ESI.Auth)
		esiGroup.GET("/callback", ESI.Callback)
	}
	authGroup := r.Group("auth")
	{
		authGroup.GET("/submit-code", Auth.SubmitCode)
	}
	return r
}
