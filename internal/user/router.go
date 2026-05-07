package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, h *UserHandler) {
	users := rg.Group("users")

	users.POST("/", h.CreateUser)
	users.GET("/", h.GetAll)
	users.GET("/:id", h.GetById)
	users.PUT("/:id", h.Update)
	users.DELETE("/:id", h.SoftDelete)
}
