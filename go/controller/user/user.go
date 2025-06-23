package user

import (
	"eve-wormhole-backend/go/service/user"
	"eve-wormhole-backend/go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Info(c *gin.Context) {
	userId, exist := c.Get("userId")
	if !exist {
		logrus.Error("User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "2",
			"message": "Failed to retrieve user information",
		})
		return
	}
	logrus.Debugf("User ID from context: %v", userId)
	user, err := user.GetUserByID(userId.(uint))
	if err != nil {
		logrus.Error("Error retrieving user information:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "2",
			"message": "Failed to retrieve user information",
		})
		return
	}
	utils.OkWithData(c, "User information retrieved successfully", gin.H{
		"roles":    []string{"admin"},
		"realName": user.Name,
	})
}
