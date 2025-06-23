package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userId uint, state string) string {
	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"state":  state,
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, _ := jwttoken.SignedString([]byte("your-secret-key"))
	return tokenString
}

func DecodeJWT(c *gin.Context, authHeader string) error {
	if authHeader == "" {
		return errors.New("authorization header is missing")
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return errors.New("authorization header format must be Bearer {token}")
	}

	// 解析并验证 Token
	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("your-secret-key"), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	// 将 Claims 存入上下文
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Set("state", claims["state"])
		if userIdFloat, ok := claims["userId"].(float64); ok {
			userId := uint(userIdFloat) // 将 float64 转换为 uint
			c.Set("userId", userId)
			if userId == 0 {
				return errors.New("invalid userId in token")
			}
		} else {
			c.Set("userId", uint(0)) // 如果转换失败，设置默认值
			return errors.New("invalid userId in token")
		}
	}
	return nil
}
