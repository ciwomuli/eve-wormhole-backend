package Auth

import (
	"eve-wormhole-backend/go/service/User"
	"eve-wormhole-backend/go/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SubmitCode(c *gin.Context) {
	err := User.Callback(c)
	if err != nil {
		logrus.Debug("Error during user callback:", err.Error())
		utils.Fail(c, "Failed to process user callback")
		return
	}
	utils.Ok(c, "Login successful")
}
