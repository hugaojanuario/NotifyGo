package channel

import "github.com/gin-gonic/gin"

func RegisterChannelConfigRoutes(rg *gin.RouterGroup, h *ChannelConfigHandler) {
	channels := rg.Group("/routes/:routeId/channels")

	channels.POST("/", h.Create)
	channels.GET("/", h.GetAll)
	channels.GET("/:id", h.GetByID)
	channels.PUT("/:id", h.Update)
	channels.DELETE("/:id", h.Delete)
}
