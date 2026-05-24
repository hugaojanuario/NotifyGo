package route

import "github.com/gin-gonic/gin"

func RegisterRouteRoutes(rg *gin.RouterGroup, h *RouteHandler) {
	routes := rg.Group("/routes")

	routes.POST("/", h.CreateRoute)
	routes.GET("/", h.GetAll)
	routes.GET("/:id", h.GetById)
	routes.PUT("/:id", h.UpdateRoute)
	routes.DELETE("/:id", h.DeleteRoute)
	routes.PATCH("/:id/toggle", h.ToggleActive)
}
