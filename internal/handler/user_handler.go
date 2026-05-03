package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/NotifyGo/internal/domain"
	"github.com/hugaojanuario/NotifyGo/internal/service"
)

type Handler struct {
	s *service.UserService
}

func NewHandler(s *service.UserService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo de req invalido"})
		return
	}

	user, err := h.s.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
