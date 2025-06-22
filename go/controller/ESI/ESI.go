package ESI

import (
	"eve-wormhole-backend/go/service/ESI"
	"eve-wormhole-backend/go/service/User"
	"eve-wormhole-backend/go/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Auth(c *gin.Context) {
	// Initialize the ESI service
	url, err := ESI.EveSSO(c)
	if err != nil {
		logrus.Debug("Error initiating SSO login:", err.Error())
		utils.Fail(c, "Failed to initiate SSO login")
		return
	}
	// Redirect the user to the SSO URL
	utils.OkWithData(c, "Redirecting to SSO", gin.H{"url": url})
}

func Callback(c *gin.Context) {
	// Handle the SSO callback
	url, err := User.Callback(c)
	if err != nil {
		logrus.Debug("Error during SSO callback:", err.Error())
		utils.Fail(c, "Failed to complete SSO login")
		return
	}
	utils.OkWithData(c, "SSO login successful", gin.H{"url": url})
}
