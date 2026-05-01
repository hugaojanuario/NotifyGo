package domain

import (
	"time"

	"github.com/google/uuid"
)

type Route struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	User              User      `json:"user"`
	KafkaConnectionID uuid.UUID `json:"kafka_connection_id"`
	Name              string    `json:"name"`
	Topic             string    `json:"topic"`
	Active            bool      `json:"active"`
	CreatedAt         time.Time `json:"created_at"`
}

type CreateRoute struct {
	UserID            uuid.UUID `json:"user_id"`
	User              User      `json:"user"`
	KafkaConnectionID uuid.UUID `json:"kafka_connection_id"`
	Name              string    `json:"name"`
	Topic             string    `json:"topic"`
}
