package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/NotifyGo/internal/auth"
	"github.com/hugaojanuario/NotifyGo/internal/authctx"
	"github.com/hugaojanuario/NotifyGo/internal/channel"
	"github.com/hugaojanuario/NotifyGo/internal/connection"
	"github.com/hugaojanuario/NotifyGo/internal/notification"
	"github.com/hugaojanuario/NotifyGo/internal/route"
	tmpl "github.com/hugaojanuario/NotifyGo/internal/template"
	"github.com/hugaojanuario/NotifyGo/internal/user"
)

func SetupRouter(
	authHandler *auth.Handler,
	userHandler *user.UserHandler,
	routeHandler *route.RouteHandler,
	connectionHandler *connection.KafkaConnectionHandler,
	channelHandler *channel.ChannelConfigHandler,
	templateHandler *tmpl.TemplateHandler,
	notificationHandler *notification.NotificationLogHandler,
) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api/v1")

	// PUBLIC ROUTES
	public := api.Group("/")

	public.POST("/auth/login", authHandler.Login)
	public.POST("/auth/register", userHandler.CreateUser)
	public.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// PROTECTED ROUTES
	protected := api.Group("/")
	protected.Use(auth.AuthMiddleware())

	user.RegisterProtectedUserRoutes(protected, userHandler)
	route.RegisterRouteRoutes(protected, routeHandler)
	connection.RegisterKafkaConnectionRoutes(protected, connectionHandler)
	channel.RegisterChannelConfigRoutes(protected, channelHandler)
	tmpl.RegisterTemplateRoutes(protected, templateHandler)
	notification.RegisterNotificationRoutes(protected, notificationHandler)

	// ADMIN ROUTES
	admin := protected.Group("/admin")
	admin.Use(authctx.RoleMiddleware("ADMIN"))
	_ = admin

	return r
}
