package kafka

import (
	"time"

	"github.com/google/uuid"
)

type KafkaConnection struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Brokers   string    `json:"brokers"`
	GroupID   string    `json:"group_id"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateKafkaConnection struct {
	Name    string `json:"name"`
	Brokers string `json:"brokers"`
	GroupID string `json:"group_id"`
}
