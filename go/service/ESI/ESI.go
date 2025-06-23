package esi

import (
	"crypto/rand"
	"encoding/base64"
	"eve-wormhole-backend/go/utils"
	"net/http"
	"time"

	"github.com/antihax/goesi"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gregjones/httpcache"
	httpcacheredis "github.com/gregjones/httpcache/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type ESI struct {
	ESI    *goesi.APIClient
	SSO    *goesi.SSOAuthenticator
	scopes []string
}

var E *ESI

func InitESI(conn redis.Conn, clientID, clientSecret, callbackURL string, scopes []string) {
	transport := httpcache.NewTransport(httpcacheredis.NewWithClient(conn))
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	//client := &http.Client{Transport: transport}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	E = &ESI{
		ESI:    goesi.NewAPIClient(client, "My App, contact 570727732@qq.com"),
		SSO:    goesi.NewSSOAuthenticatorV2(client, clientID, clientSecret, callbackURL, scopes),
		scopes: scopes,
	}
}

func EveSSO(c *gin.Context) (map[string]any, error) {
	userId, exist := c.Get("userId")
	if !exist {
		userId = (uint)(0)
	}
	// Generate a random state string
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	// Save the state on the session
	logrus.Debugf("Generated state: %s", state)
	logrus.Debugf("User ID: %v", userId)
	tokenString := utils.GenerateJWT(userId.(uint), state)
	// Generate the SSO URL with the state string
	url := E.SSO.AuthorizeURL(state, true, E.scopes)
	// Send the user to the URL
	return gin.H{"url": url, "token": tokenString}, nil
}

type CallbackData struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

func EveSSOCallback(c *gin.Context) (*oauth2.Token, error) {
	// Get the state from the session
	var callbackData CallbackData
	if err := c.ShouldBindJSON(&callbackData); err != nil {
		logrus.Debugf("Failed to bind JSON: %v", err)
		return nil, utils.WrapError(errors.Wrap(err, "invalid JSON"))
	}
	state := callbackData.State
	code := callbackData.Code
	state_token, exist := c.Get("state")
	if !exist || state == "" || state_token != state {
		logrus.Printf("State mismatch: expected %s, got %s", state_token, state)
		return nil, utils.WrapError(gin.Error{
			Err:  http.ErrNoCookie,
			Type: gin.ErrorTypePublic,
		})
	}
	if code == "" {
		return nil, utils.WrapError(gin.Error{
			Err:  http.ErrNoCookie,
			Type: gin.ErrorTypePublic,
		})
	}

	token, err := E.SSO.TokenExchange(code)
	if err != nil {
		return nil, utils.WrapError(errors.Wrap(err, "token exchange error"))
	}

	// Obtain a token source (automaticlly pulls refresh as needed)
	tokSrc := E.SSO.TokenSource(token)

	// Verify the client (returns clientID)
	_, err = E.SSO.Verify(tokSrc)
	if err != nil {
		return nil, utils.WrapError(errors.Wrap(err, "token verify error"))
	}

	token, err = tokSrc.Token()
	if err != nil {
		return nil, utils.WrapError(errors.Wrap(err, "token source error getting new token"))
	}
	return token, nil
}
