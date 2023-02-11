package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	logCtx "github.com/si-bas/go-rest-boilerplate/pkg/logger/context"
	"github.com/si-bas/go-rest-boilerplate/shared/constant"
)

func InjectContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(constant.XRequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Request = c.Request.WithContext(logCtx.InjectRequestID(c.Request.Context(), requestID))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.XRequestIDHeader, requestID))

		c.Next()
	}
}
