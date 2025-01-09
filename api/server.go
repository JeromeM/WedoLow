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
	// Create a new Gin router
	router := gin.Default()

	// Initialize the database and services
	db := database.NewPostgresDB(cfg.DatabaseURL)
	userDb := database.NewUserDatabase(db)

	// Initialize the RandomUser client
	randomUserClient := service.NewRandomUserClient(cfg.RandomUserAPI)

	// Initialize the user service and handler
	userService := service.NewUserService(userDb.(*database.UserDatabase), randomUserClient)
	handler := handler.NewUserHandler(userService)

	// Add OpenTelemetry middleware
	router.Use(otelgin.Middleware("WedoLow"))

	// Define the routes
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
