package middleware

import (
	"eve-wormhole-backend/go/service/user"
	"eve-wormhole-backend/go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		skipRoutes := map[string]bool{
			"/esi/auth":         true, // 示例路由
			"/esi/callback":     true, // 示例路由
			"/auth/submit-code": true,
			"/wormhole/ws":      true,
		}

		err := utils.DecodeJWT(c, c.Request.Header.Get("Authorization"))
		if err != nil {
			if skipRoutes[c.Request.URL.Path] {
				c.Next()
				return
			}
			logrus.Debug("JWT decode error:", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": "1", "message": "Unauthorized"})
			return
		}
		userId, _ := c.Get("userId")
		err = user.UpdateUserActiveTimebyID(userId.(uint))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": "2", "message": "Internal Server Error"})
			return
		}
		c.Next()
	}
}
