package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChannelType string

const (
	Email   ChannelType = "EMAIL"
	Slack   ChannelType = "SLACK"
	Webhook ChannelType = "WEBHOOK"
	Sms     ChannelType = "SMS"
)

type ChannelConfig struct {
	ID              uuid.UUID   `json:"id"`
	RouteID         uuid.UUID   `json:"route_id"`
	ChannelType     ChannelType `json:"channel_type"`
	ToField         *string     `json:"to_field"`
	ToFixed         *string     `json:"to_fixed"`
	Subject         *string     `json:"subject"`
	TemplateID      *uuid.UUID  `json:"template_id"`
	MessageTemplate *string     `json:"mensage_template"`
	WebhookURL      *string     `json:"webhook_url"`
	WebhookSecret   *string     `json:"webhook_secret"`
	SlackChannel    *string     `json:"slack_channel"`
	Active          bool        `json:"active"`
	CreatedAt       time.Time   `json:""created_at`
}
