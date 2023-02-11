package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-boilerplate/shared/helper/response"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(response.NewJSONResponse().APIStatusSuccess().StatusCode, response.NewJSONResponse().SetData("OK"))
}
