package gateway

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hugaojanuario/NotifyGo/internal/channel"
	"github.com/hugaojanuario/NotifyGo/internal/notification"
)

type Dispatcher struct {
	registry    *channel.Registry
	channelRepo channel.ChannelConfigRepositoryMethods
	logRepo     notification.NotificationLogRepositoryMethods
}

func NewDispatcher(
	registry *channel.Registry,
	channelRepo channel.ChannelConfigRepositoryMethods,
	logRepo notification.NotificationLogRepositoryMethods,
) *Dispatcher {
	return &Dispatcher{
		registry:    registry,
		channelRepo: channelRepo,
		logRepo:     logRepo,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, routeID uuid.UUID, topic string, payload string) {
	configs, err := d.channelRepo.GetAllByRouteID(ctx, routeID)
	if err != nil {
		log.Printf("dispatcher - get channel configs for route %s: %v", routeID, err)
		return
	}

	var wg sync.WaitGroup

	for _, cfg := range configs {
		if !cfg.Active {
			continue
		}

		wg.Add(1)
		go func(cfg channel.ChannelConfig) {
			defer wg.Done()

			ch, ok := d.registry.Get(cfg.ChannelType)
			if !ok {
				log.Printf("dispatcher - no handler for channel type %s", cfg.ChannelType)
				return
			}

			event := channel.Event{Topic: topic, Payload: payload}
			err := ch.Send(ctx, event, cfg)

			status := notification.StatusSuccess
			var errMsg *string
			var sentAt *time.Time

			if err != nil {
				status = notification.StatusFailed
				msg := err.Error()
				errMsg = &msg
				log.Printf("dispatcher - send via %s failed: %v", cfg.ChannelType, err)
			} else {
				now := time.Now()
				sentAt = &now
			}

			recipient, _ := resolveRecipient(payload, cfg.ToField, cfg.ToFixed)

			_, logErr := d.logRepo.Create(ctx, notification.CreateNotificationLog{
				RouteID:         routeID,
				ChannelConfigID: cfg.ID,
				Topic:           topic,
				Channel:         cfg.ChannelType,
				Recipient:       recipient,
				Status:          status,
				Payload:         payload,
				ErrorMessage:    errMsg,
				Attempts:        1,
				SentAt:          sentAt,
			})
			if logErr != nil {
				log.Printf("dispatcher - save notification log: %v", logErr)
			}
		}(cfg)
	}

	wg.Wait()
}

func resolveRecipient(payload string, toField *string, toFixed *string) (string, error) {
	return channel.ExtractRecipient(payload, toField, toFixed)
}
