package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hugaojanuario/NotifyGo/internal/channel"
	"github.com/hugaojanuario/NotifyGo/internal/notification"
	"github.com/hugaojanuario/NotifyGo/pkg/config"
	"github.com/hugaojanuario/NotifyGo/pkg/database"
	"github.com/hugaojanuario/NotifyGo/pkg/gateway"
	"github.com/hugaojanuario/NotifyGo/pkg/kafka"
)

func main() {

	cfg := config.LoadDotEnv()

	db, err := database.DBConn(cfg)
	if err != nil {
		log.Fatalf("gateway - database connection: %v", err)
	}
	defer db.Close()

	// CHANNEL REGISTRY
	registry := channel.NewRegistry()
	registry.Register(channel.ChannelEmail, channel.NewEmailChannel(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Password,
		cfg.SMTP.From,
	))
	registry.Register(channel.ChannelWebhook, channel.NewWebhookChannel())
	registry.Register(channel.ChannelSlack, channel.NewSlackChannel())
	registry.Register(channel.ChannelSMS, channel.NewSMSChannel(
		cfg.Twilio.AccountSID,
		cfg.Twilio.AuthToken,
		cfg.Twilio.From,
	))

	// REPOSITORIES
	channelRepo := channel.NewChannelConfigRepository(db)
	logRepo := notification.NewNotificationLogRepository(db)

	// DISPATCHER
	dispatcher := gateway.NewDispatcher(registry, channelRepo, logRepo)

	// CONSUMER MANAGER
	manager := kafka.NewConsumerManager(db, dispatcher)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	if err := manager.Start(ctx, &wg); err != nil {
		log.Fatalf("gateway - start consumers: %v", err)
	}

	log.Println("gateway started — waiting for events")

	<-ctx.Done()
	log.Println("gateway shutting down")

	wg.Wait()
	log.Println("gateway stopped")
}
