package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/si-bas/go-rest-boilerplate/config"
	"github.com/si-bas/go-rest-boilerplate/shared/constant"
)

func AuthJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(constant.AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "bad header value given"})
			c.Abort()
			return
		}

		splitAuthHeader := strings.Split(authHeader, " ")
		if len(splitAuthHeader) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "incorrectly formatted authorization header"})
			c.Abort()
			return
		}

		token, err := parseToken(splitAuthHeader[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unable to parse claims"})
			c.Abort()
		}

		userId := uint32(claims["sub"].(float64))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.UserID, strconv.FormatUint(uint64(userId), 10)))
		c.Next()
	}
}

func parseToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(config.Config.Jwt.Secret), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token")
	}

	return token, nil
}
