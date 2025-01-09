package handler

import (
	"net/http"
	"strconv"
	"users-service/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Handler for creating users
func (h *UserHandler) CreateUsers(c *gin.Context) {
	// Tracing
	tracer := otel.Tracer("handler/CreateUser")
	ctx, span := tracer.Start(c.Request.Context(), "CreateUserHandler")
	defer span.End()

	count, _ := strconv.Atoi(c.DefaultQuery("count", "1"))
	gender := c.DefaultQuery("gender", "any")

	response, err := h.service.CreateUsers(ctx, count, gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Handler for getting a user
func (h *UserHandler) GetUser(c *gin.Context) {
	tracer := otel.Tracer("handler/GetUser")
	ctx, span := tracer.Start(c.Request.Context(), "GetUserHandler")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Handler for listing users
func (h *UserHandler) ListUsers(c *gin.Context) {
	tracer := otel.Tracer("handler/ListUsers")
	ctx, span := tracer.Start(c.Request.Context(), "ListUsersHandler")
	defer span.End()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	nameFilter := c.DefaultQuery("name", "")

	users, err := h.service.ListUsers(ctx, page, limit, nameFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
