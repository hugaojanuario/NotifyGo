package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/NotifyGo/internal/user"
)

func SetupRouter(userHandler *user.UserHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("api/v1")

	user.RegisterUserRoutes(api, userHandler)

	return r
}
