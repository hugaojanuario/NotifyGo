package template

import "github.com/gin-gonic/gin"

func RegisterTemplateRoutes(rg *gin.RouterGroup, h *TemplateHandler) {
	templates := rg.Group("/templates")

	templates.POST("/", h.Create)
	templates.GET("/", h.GetAll)
	templates.GET("/:id", h.GetByID)
	templates.PUT("/:id", h.Update)
	templates.DELETE("/:id", h.Delete)
	templates.POST("/:id/preview", h.Preview)
}
