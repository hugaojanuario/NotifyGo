package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/NotifyGo/internal/handler"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("api/v1")

	RegisterUserRoutes(api, userHandler)

	return r
}
