package notification

import (
	"time"

	"github.com/google/uuid"
	"github.com/hugaojanuario/NotifyGo/internal/domain"
)

type NotificationStatus string

const (
	StatusSuccess  NotificationStatus = "SUCCESS"
	StatusFailed   NotificationStatus = "FAILED"
	StatusRetrying NotificationStatus = "RETRYING"
)

type NotificationLog struct {
	ID              uuid.UUID          `json:"id"`
	RouteID         uuid.UUID          `json:"route_id"`
	ChannelConfigID uuid.UUID          `json:"channel_config_id"`
	Topic     string             `json:"topic"`
	Channel   domain.ChannelType `json:"channel_type"`
	Recipient string             `json:"recipient"`
	Status    NotificationStatus `json:"status"`
	Payload   string             `json:"payload"`
	ErrorMessage    *string            `json:"error_message"`
	Attempts        int                `json:"attempts"`
	SentAt          *time.Time         `json:"sent_at"`
	CreatedAt       time.Time          `json:"created_at"`
}
