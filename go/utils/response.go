package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context, msg string) {
	if msg == "" { // 如果参数为空，设置默认值
		msg = "ok"
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": msg})
}

func Fail(c *gin.Context, msg string) {
	if msg == "" { // 如果参数为空，设置默认值
		msg = "fail"
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": msg})
}

func OkWithData(c *gin.Context, msg string, data any) {
	if msg == "" { // 如果参数为空，设置默认值
		msg = "ok"
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": msg, "data": data})
}
