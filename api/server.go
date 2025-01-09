package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

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
