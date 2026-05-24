package connection

import "github.com/gin-gonic/gin"

func RegisterKafkaConnectionRoutes(rg *gin.RouterGroup, h *KafkaConnectionHandler) {
	connections := rg.Group("/kafka-connections")

	connections.POST("/", h.Create)
	connections.GET("/", h.GetAll)
	connections.GET("/:id", h.GetByID)
	connections.PUT("/:id", h.Update)
	connections.DELETE("/:id", h.Delete)
	connections.POST("/:id/test", h.TestConnection)
}
