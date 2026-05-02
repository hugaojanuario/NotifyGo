package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChannelConfig struct {
	ID              uuid.UUID   `json:"id"`
	RouteID         uuid.UUID   `json:"route_id"`
	ChannelType     ChannelType `json:"channel_type"`
	ToField         *string     `json:"to_field"`
	ToFixed         *string     `json:"to_fixed"`
	Subject         *string     `json:"subject"`
	TemplateID      *uuid.UUID  `json:"template_id"`
	MessageTemplate *string     `json:"message_template"`
	WebhookURL      *string     `json:"webhook_url"`
	WebhookSecret   *string     `json:"webhook_secret"`
	SlackChannel    *string     `json:"slack_channel"`
	Active          bool        `json:"active"`
	CreatedAt       time.Time   `json:"created_at"`
}

type CreateChannelConfig struct {
	ChannelType     ChannelType `json:"channel_type"`
	ToField         *string     `json:"to_field"`
	ToFixed         *string     `json:"to_fixed"`
	Subject         *string     `json:"subject"`
	TemplateID      *uuid.UUID  `json:"template_id"`
	MessageTemplate *string     `json:"message_template"`
	WebhookURL      *string     `json:"webhook_url"`
	WebhookSecret   *string     `json:"webhook_secret"`
	SlackChannel    *string     `json:"slack_channel"`
	Active          bool        `json:"active"`
}
