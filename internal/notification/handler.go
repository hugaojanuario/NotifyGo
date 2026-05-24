package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hugaojanuario/NotifyGo/internal/authctx"
	"github.com/hugaojanuario/NotifyGo/internal/channel"
)

type NotificationLogHandler struct {
	r NotificationLogRepositoryMethods
}

func NewNotificationLogHandler(r NotificationLogRepositoryMethods) *NotificationLogHandler {
	return &NotificationLogHandler{r: r}
}

func (h *NotificationLogHandler) GetAll(c *gin.Context) {
	userID, err := authctx.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filters := LogFilters{
		Page:  1,
		Limit: 20,
	}

	if routeIDStr := c.Query("route_id"); routeIDStr != "" {
		routeID, err := uuid.Parse(routeIDStr)
		if err == nil {
			filters.RouteID = &routeID
		}
	}

	if status := c.Query("status"); status != "" {
		s := NotificationStatus(status)
		filters.Status = &s
	}

	if ct := c.Query("channel"); ct != "" {
		c_ := channel.ChannelType(ct)
		filters.ChannelType = &c_
	}

	logs, total, err := h.r.GetAll(c.Request.Context(), userID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  filters.Page,
		"limit": filters.Limit,
	})
}

func (h *NotificationLogHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalid"})
		return
	}

	log, err := h.r.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if log == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

func (h *NotificationLogHandler) GetMetrics(c *gin.Context) {
	userID, err := authctx.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	metrics, err := h.r.GetMetrics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
