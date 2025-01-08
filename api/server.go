package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"users-service/config"
)

type Server struct {
	config *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	router := gin.Default()
	//db := database.NewPostgresDB(cfg.DatabaseURL)
	//userDb := database.NewUserDatabase(db)
	//andomUserClient := service.NewRandomUserClient(cfg.RandomUserAPI)

	return &Server{
		config: cfg,
		router: router,
	}
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.Port))
}
