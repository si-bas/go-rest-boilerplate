package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-boilerplate/config"
	"github.com/si-bas/go-rest-boilerplate/domain/repository"
	"github.com/si-bas/go-rest-boilerplate/pkg/gorm"
	"github.com/si-bas/go-rest-boilerplate/pkg/logger"
	"github.com/si-bas/go-rest-boilerplate/server/handler"
	"github.com/si-bas/go-rest-boilerplate/server/middleware"
	"github.com/si-bas/go-rest-boilerplate/service"
	"github.com/si-bas/go-rest-boilerplate/shared/constant"
)

type HTTPServer struct {
}

// New to instantiate HTTPServer
func New() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Start() {
	h := initHandler()

	if config.Config.App.Env == constant.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	if config.Config.App.Env != constant.EnvProduction {
		router.Use(middleware.CORS())
	}

	router.Use(middleware.InjectContext())
	router.GET("/healthcheck", h.HealthCheck)

	groupV1 := router.Group("/v1")
	groupV1.POST("/auth/token", h.GetToken)
	groupV1.POST("/auth/refresh", h.RefreshToken)

	groupV1.Use(middleware.AuthJwt())
	groupV1.GET("/auth/me", h.GetMe)

	groupV1.POST("/user", h.CreateUser)
	groupV1.GET("/user", h.ListUser)
	groupV1.GET("/user/:id", h.DetailUser)

	err := router.Run(fmt.Sprintf(":%d", config.Config.App.Port))
	if err != nil {
		logger.Error(context.Background(), "failed to run router", err)
	}
}

func initHandler() *handler.Handler {
	var err error
	config.TimeLocation, err = time.LoadLocation(config.Config.App.Timezone)
	if err != nil {
		panic("error set timezone, err=" + err.Error())
	}

	logger.InitLogger()

	// TODO: init DB
	db := gorm.ConnectDB()

	// TODO: init repositories
	userRepo := repository.NewUserRepository(db)

	// TODO: init pkgs

	// TODO: init services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	return handler.New(
		authService,
		userService,
	)
}
