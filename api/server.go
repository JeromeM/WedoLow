package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"users-service/config"
	"users-service/database"
	"users-service/handler"
	"users-service/service"
)

type Server struct {
	config *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	router := gin.Default()
	db := database.NewPostgresDB(cfg.DatabaseURL)
	userDb := database.NewUserDatabase(db)
	randomUserClient := service.NewRandomUserClient(cfg.RandomUserAPI)
	userService := service.NewUserService(userDb.(*database.UserDatabase), randomUserClient)
	handler := handler.NewUserHandler(userService)

	// Add OpenTelemetry middleware
	router.Use(otelgin.Middleware("WedoLow"))

	router.POST("/user", handler.CreateUsers)
	router.GET("/user/:id", handler.GetUser)
	router.GET("/users", handler.ListUsers)

	return &Server{
		config: cfg,
		router: router,
	}
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.Port))
}
