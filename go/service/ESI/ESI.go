package ESI

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/antihax/goesi"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gregjones/httpcache"
	httpcacheredis "github.com/gregjones/httpcache/redis"
	"github.com/pkg/errors"
)

type ESI struct {
	esi    *goesi.APIClient
	sso    *goesi.SSOAuthenticator
	scopes []string
}

var e *ESI

func InitESI(conn redis.Conn, clientID, clientSecret, callbackURL string, scopes []string) {
	transport := httpcache.NewTransport(httpcacheredis.NewWithClient(conn))
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}
	e = &ESI{
		esi:    goesi.NewAPIClient(client, "My App, contact someone@nowhere"),
		sso:    goesi.NewSSOAuthenticatorV2(client, clientID, clientSecret, callbackURL, scopes),
		scopes: scopes,
	}
}

func EveSSO(c *gin.Context) (string, error) {
	s := sessions.Default(c)
	// Generate a random state string
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	// Save the state on the session
	s.Set("state", state)
	err := s.Save()
	if err != nil {
		return "", err
	}
	// Generate the SSO URL with the state string
	url := e.sso.AuthorizeURL(state, true, e.scopes)
	// Send the user to the URL
	return url, nil
}

func EveSSOCallback(c *gin.Context) error {
	s := sessions.Default(c)
	// Get the state from the session
	state := c.Query("state")
	if state != "" && s.Get("state") != state {
		return gin.Error{
			Err:  http.ErrNoCookie,
			Type: gin.ErrorTypePublic,
		}
	}
	// Get the code from the query parameters
	code := c.Query("code")
	if code == "" {
		return gin.Error{
			Err:  http.ErrNoCookie,
			Type: gin.ErrorTypePublic,
		}
	}

	token, err := e.sso.TokenExchange(code)
	if err != nil {
		return errors.Wrap(err, "token exchange error")
	}

	// Obtain a token source (automaticlly pulls refresh as needed)
	tokSrc := e.sso.TokenSource(token)

	// Verify the client (returns clientID)
	v, err := e.sso.Verify(tokSrc)
	if err != nil {
		return errors.Wrap(err, "token verify error")
	}

	token, err = tokSrc.Token()
	if err != nil {
		return errors.Wrap(err, "token source error getting new token")
	}
	// Save token.
	s.Set("token", *token)

	// Save the verification structure on the session for quick access.
	s.Set("character", v)
	err = s.Save()
	if err != nil {
		return errors.Wrap(err, "unable to save session")
	}
	return nil
}
