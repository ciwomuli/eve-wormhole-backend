package ESI

import (
	"eve-wormhole-backend/go/service/ESI"
	"eve-wormhole-backend/go/utils"
	"net/http"

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
	c.Redirect(http.StatusFound, "http://127.0.0.1:5777/auth/wait?code="+c.Query("code")+"&state="+c.Query("state"))
	//c.Redirect(http.StatusFound, "http://baidu.com")
}

