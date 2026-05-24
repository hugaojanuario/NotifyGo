package notification

import "github.com/gin-gonic/gin"

func RegisterNotificationRoutes(rg *gin.RouterGroup, h *NotificationLogHandler) {
	rg.GET("/logs", h.GetAll)
	rg.GET("/logs/:id", h.GetByID)
	rg.GET("/metrics", h.GetMetrics)
}
