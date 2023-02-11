package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-boilerplate/config"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		//origin header exists
		if origin := c.Request.Header.Get("origin"); len(origin) > 1 {
			if strings.HasSuffix(origin, config.Config.App.Url) || config.Config.App.Env != "production" {
				corsHeaders(c, origin)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		} else {
			corsHeaders(c, "*")
			c.Next()
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func corsHeaders(c *gin.Context, origin string) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
}
