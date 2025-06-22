package ESI

import (
	"eve-wormhole-backend/go/service/ESI"
	"eve-wormhole-backend/go/service/User"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Auth(c *gin.Context) {
	// Initialize the ESI service
	url, err := ESI.EveSSO(c)
	if err != nil {
		logrus.Debug("Error initiating SSO login:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate SSO login"})
		return
	}
	// Redirect the user to the SSO URL
	c.Redirect(http.StatusFound, url)
}

func Callback(c *gin.Context) {
	// Handle the SSO callback
	url, err := User.Callback(c)
	if err != nil {
		logrus.Debug("Error during SSO callback:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete SSO login"})
		return
	}
	c.Redirect(http.StatusFound, url)
}
