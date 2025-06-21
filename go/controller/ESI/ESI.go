package ESI

import (
	"eve-wormhole-backend/go/service/ESI"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	// Initialize the ESI service
	url, err := ESI.EveSSO(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate SSO login"})
		return
	}
	// Redirect the user to the SSO URL
	c.Redirect(http.StatusFound, url)
}

func Callback(c *gin.Context) {
	// Handle the SSO callback
	s := sessions.Default(c)
	err := ESI.EveSSOCallback(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete SSO login"})
		return
	}
	if s.Get("login") == nil {
		s.Set("login", true)
		err = s.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		c.Redirect(http.StatusFound, "/home")
	}
	c.Redirect(http.StatusFound, "/esi")
}
