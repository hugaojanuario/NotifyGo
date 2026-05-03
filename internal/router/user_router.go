package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/NotifyGo/internal/handler"
)

func RegisterUserRoutes(rg *gin.RouterGroup, h *handler.UserHandler) {
	users := rg.Group("users")

	users.POST("/", h.CreateUser)
	users.GET("/", h.GetAll)
	users.GET("/:id", h.GetById)
	users.PUT("/:id", h.Update)
	users.DELETE("/:id", h.Update)
}
