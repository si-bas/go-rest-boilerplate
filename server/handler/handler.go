package handler

import (
	"github.com/si-bas/go-rest-boilerplate/service"
)

type Handler struct {
	userService service.UserService
	authService service.AuthService
}

func New(
	authService service.AuthService,
	userService service.UserService) *Handler {
	return &Handler{
		userService: userService,
		authService: authService,
	}
}
