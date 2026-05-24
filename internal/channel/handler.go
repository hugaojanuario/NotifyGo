package channel

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChannelConfigHandler struct {
	s *ChannelConfigService
}

func NewChannelConfigHandler(s *ChannelConfigService) *ChannelConfigHandler {
	return &ChannelConfigHandler{s: s}
}

func (h *ChannelConfigHandler) Create(c *gin.Context) {
	routeID, err := uuid.Parse(c.Param("routeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "routeId invalid"})
		return
	}

	var req CreateChannelConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body invalid"})
		return
	}

	cfg, err := h.s.Create(c.Request.Context(), routeID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cfg)
}

func (h *ChannelConfigHandler) GetAll(c *gin.Context) {
	routeID, err := uuid.Parse(c.Param("routeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "routeId invalid"})
		return
	}

	configs, err := h.s.GetAllByRouteID(c.Request.Context(), routeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func (h *ChannelConfigHandler) GetByID(c *gin.Context) {
	routeID, err := uuid.Parse(c.Param("routeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "routeId invalid"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalid"})
		return
	}

	cfg, err := h.s.GetByID(c.Request.Context(), id, routeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cfg)
}

func (h *ChannelConfigHandler) Update(c *gin.Context) {
	routeID, err := uuid.Parse(c.Param("routeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "routeId invalid"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalid"})
		return
	}

	var req CreateChannelConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body invalid"})
		return
	}

	cfg, err := h.s.Update(c.Request.Context(), id, routeID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cfg)
}

func (h *ChannelConfigHandler) Delete(c *gin.Context) {
	routeID, err := uuid.Parse(c.Param("routeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "routeId invalid"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalid"})
		return
	}

	if err := h.s.Delete(c.Request.Context(), id, routeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
