package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-boilerplate/domain/model"
	"github.com/si-bas/go-rest-boilerplate/pkg/logger"
	"github.com/si-bas/go-rest-boilerplate/pkg/logger/tag"
	"github.com/si-bas/go-rest-boilerplate/shared"
	"github.com/si-bas/go-rest-boilerplate/shared/constant"
	"github.com/si-bas/go-rest-boilerplate/shared/helper/response"
)

func (h *Handler) GetToken(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse()

	var payload model.AuthTokenRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Warn(ctx, "failed to bindJSON", tag.Err(err))
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}

	user, err := h.authService.ValidateUser(ctx, model.ValidateUser(payload))
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("invalid password user: %s", payload.Email), tag.Err(err))
		c.JSON(result.APIStatusInvalidAuthentication().StatusCode, result.SetError(response.ErrUnauthorized, err.Error()))
	}

	jwtToken, err := h.authService.GenerateToken(ctx, user)
	if err != nil {
		logger.Warn(ctx, "failed generate jwt token", tag.Err(err))
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(jwtToken))
}

func (h *Handler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse()

	var payload model.AuthRefreshTokenRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Warn(ctx, "failed to bindJSON", tag.Err(err))
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}

	token, err := h.authService.ParseToken(ctx, payload.RefreshToken)
	if err != nil {
		logger.Warn(ctx, err.Error(), tag.Err(err))
		c.JSON(result.APIStatusInvalidAuthentication().StatusCode, result.SetError(response.ErrUnauthorized, err.Error()))
		return
	}

	claims, err := h.authService.GetClaims(ctx, token)
	if err != nil {
		logger.Warn(ctx, err.Error(), tag.Err(err))
		c.JSON(result.APIStatusInvalidAuthentication().StatusCode, result.SetError(response.ErrUnauthorized, err.Error()))
		return
	}

	user, err := h.authService.GetUser(ctx, uint32(claims["sub"].(float64)))
	if err != nil {
		logger.Warn(ctx, "failed to get user", tag.Err(err))
		c.JSON(result.APIStatusInvalidAuthentication().StatusCode, result.SetError(response.ErrUnauthorized, err.Error()))
		return
	}

	jwtToken, err := h.authService.GenerateToken(ctx, user)
	if err != nil {
		logger.Warn(ctx, "failed generate jwt token", tag.Err(err))
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(jwtToken))
}

func (h *Handler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse()

	userId, err := strconv.ParseUint(shared.GetContextValueAsString(ctx, constant.UserID), 10, 32)
	if err != nil {
		logger.Warn(ctx, "failed get user from context", tag.Err(err))
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}

	user, err := h.authService.GetUser(ctx, uint32(userId))
	if err != nil {
		logger.Warn(ctx, "failed to get user", tag.Err(err))
		c.JSON(result.APIStatusInvalidAuthentication().StatusCode, result.SetError(response.ErrUnauthorized, err.Error()))
		return
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(user))
}
