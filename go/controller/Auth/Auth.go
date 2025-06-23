package auth

import (
	"eve-wormhole-backend/go/service/user"
	"eve-wormhole-backend/go/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SubmitCode(c *gin.Context) {
	data, err := user.Callback(c)
	if err != nil {
		logrus.Debug("Error during user callback:", err.Error())
		utils.Fail(c, "Failed to process user callback")
		return
	}
	utils.OkWithData(c, "User callback processed successfully", data)
}
