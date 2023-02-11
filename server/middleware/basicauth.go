package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-boilerplate/config"
	"github.com/si-bas/go-rest-boilerplate/shared/constant"
)

func BasicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "forbidden"})
			c.Abort()
			return
		}

		if username != config.Config.Credential.Username || password != config.Config.Credential.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "forbidden"})
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.User, username))

		c.Next()
	}

}
