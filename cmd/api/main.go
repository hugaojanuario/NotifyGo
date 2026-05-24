package main

import (
	"log"

	"github.com/hugaojanuario/NotifyGo/internal/auth"
	"github.com/hugaojanuario/NotifyGo/internal/channel"
	"github.com/hugaojanuario/NotifyGo/internal/connection"
	"github.com/hugaojanuario/NotifyGo/internal/notification"
	"github.com/hugaojanuario/NotifyGo/internal/route"
	"github.com/hugaojanuario/NotifyGo/internal/server"
	tmpl "github.com/hugaojanuario/NotifyGo/internal/template"
	"github.com/hugaojanuario/NotifyGo/internal/user"
	"github.com/hugaojanuario/NotifyGo/pkg/config"
	"github.com/hugaojanuario/NotifyGo/pkg/database"
)

func main() {

	cfg := config.LoadDotEnv()

	db, err := database.DBConn(cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()

	// USER
	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	// AUTH
	authHandler := auth.NewHandler(userService)

	// ROUTE
	routeRepository := route.NewRouteRepository(db)
	routeService := route.NewRouteService(routeRepository)
	routeHandler := route.NewRouteHandler(routeService)

	// KAFKA CONNECTION
	connectionRepository := connection.NewKafkaRepository(db)
	connectionService := connection.NewKafkaConnectionService(connectionRepository)
	connectionHandler := connection.NewKafkaConnectionHandler(connectionService)

	// CHANNEL CONFIG
	channelRepository := channel.NewChannelConfigRepository(db)
	channelService := channel.NewChannelConfigService(channelRepository)
	channelHandler := channel.NewChannelConfigHandler(channelService)

	// TEMPLATE
	templateRepository := tmpl.NewTemplateRepository(db)
	templateService := tmpl.NewTemplateService(templateRepository)
	templateHandler := tmpl.NewTemplateHandler(templateService)

	// NOTIFICATION LOG
	notificationRepository := notification.NewNotificationLogRepository(db)
	notificationHandler := notification.NewNotificationLogHandler(notificationRepository)

	// SERVER
	router := server.SetupRouter(
		authHandler,
		userHandler,
		routeHandler,
		connectionHandler,
		channelHandler,
		templateHandler,
		notificationHandler,
	)

	router.Run(":9292")
}
