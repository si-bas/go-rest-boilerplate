package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-boilerplate/domain/model"
	"github.com/si-bas/go-rest-boilerplate/pkg/logger"
	"github.com/si-bas/go-rest-boilerplate/pkg/logger/tag"
	"github.com/si-bas/go-rest-boilerplate/shared/helper/pagination"
	"github.com/si-bas/go-rest-boilerplate/shared/helper/response"
	"gorm.io/gorm"
)

func (h *Handler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse()

	var payload model.CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Warn(ctx, "failed to bindJSON", tag.Err(err))
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}

	if emailIsUsed, err := h.userService.EmailIsUsed(ctx, payload.Email); emailIsUsed || err != nil {
		if err != nil {
			logger.Error(ctx, "error user email check", err)
			c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
			return
		}

		c.JSON(result.APIStatusConflict().StatusCode, result.SetError(response.ErrConflict, errors.New("email already used").Error()))
		return
	}

	user, err := h.userService.Create(ctx, model.CreateUser(payload))
	if err != nil {
		logger.Warn(ctx, "failed to create user", tag.Err(err))
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}

	c.JSON(result.APIStatusCreated().StatusCode, result.SetData(user))
}

func (h *Handler) ListUser(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse()

	var query model.UserListRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Warn(ctx, "failed to bindQuery", tag.Err(err))
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}

	var sortBys []pagination.ParamSort
	if len(query.Sort) > 0 {
		for k, v := range query.Sort {
			sortBys = append(sortBys, pagination.ParamSort{
				Column: k,
				Order:  v,
			})
		}
	}

	users, meta, err := h.userService.ListPaginate(ctx, model.UserFilter{
		Keyword: query.Keyword,
		Name:    query.Name,
		Email:   query.Email,
	}, pagination.Param{
		Limit: query.Limit,
		Page:  query.Page,
		Sort:  sortBys,
	})
	if err != nil {
		logger.Warn(ctx, "failed to get users with pagination", tag.Err(err))
		c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
		return
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(users).SetMeta(meta))
}

func (h *Handler) DetailUser(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse()

	var payload model.UserFind
	if err := c.BindUri(&payload); err != nil {
		logger.Warn(ctx, "failed to bindURI", tag.Err(err))
		c.JSON(result.APIStatusBadRequest().StatusCode, result.SetError(response.ErrBadRequest, err.Error()))
		return
	}

	user, err := h.userService.Detail(ctx, payload.ID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Warn(ctx, "failed to get user detail", tag.Err(err))
			c.JSON(result.APIInternalServerError().StatusCode, result.SetError(response.ErrInternalServerError, err.Error()))
			return
		}

		c.JSON(result.APIStatusNotFound().StatusCode, result.SetError(response.ErrNotFound, errors.New("user not found").Error()))
		return
	}

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData(user))
}
