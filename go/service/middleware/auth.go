package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		skipRoutes := map[string]bool{
			"/esi/auth":         true, // 示例路由
			"/esi/callback":     true, // 示例路由
			"/auth/submit-code": true,
		}

		// 检查当前请求的路径是否在跳过列表中
		if skipRoutes[c.Request.URL.Path] {
			c.Next() // 跳过权限检测，继续处理
			return
		}

		s := sessions.Default(c)

		if s.Get("login") != true {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort() // 停止后续处理
			return
		}
		c.Next()
	}
}
