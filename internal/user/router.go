package user

import "github.com/gin-gonic/gin"

func RegisterProtectedUserRoutes(
	rg *gin.RouterGroup,
	h *UserHandler,
) {

	users := rg.Group("/users")

	users.GET("/me", h.GetMe)
	users.PUT("/me", h.Update)
	users.DELETE("/me", h.SoftDelete)
}
