package kafka

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
	segkafka "github.com/segmentio/kafka-go"
	"github.com/hugaojanuario/NotifyGo/pkg/gateway"
)

type topicRoute struct {
	routeID uuid.UUID
	topic   string
	brokers []string
	groupID string
}

type ConsumerManager struct {
	db         *sql.DB
	dispatcher *gateway.Dispatcher
}

func NewConsumerManager(db *sql.DB, dispatcher *gateway.Dispatcher) *ConsumerManager {
	return &ConsumerManager{db: db, dispatcher: dispatcher}
}

func (m *ConsumerManager) Start(ctx context.Context, wg *sync.WaitGroup) error {
	routes, err := m.loadActiveRoutes(ctx)
	if err != nil {
		return err
	}

	for _, r := range routes {
		wg.Add(1)
		go func(r topicRoute) {
			defer wg.Done()
			m.consume(ctx, r)
		}(r)
	}

	return nil
}

func (m *ConsumerManager) consume(ctx context.Context, r topicRoute) {
	reader := segkafka.NewReader(segkafka.ReaderConfig{
		Brokers: r.brokers,
		Topic:   r.topic,
		GroupID: r.groupID,
	})
	defer reader.Close()

	log.Printf("kafka consumer started: topic=%s route=%s", r.topic, r.routeID)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("kafka consumer stopped: topic=%s", r.topic)
				return
			}
			log.Printf("kafka consumer read error: topic=%s err=%v", r.topic, err)
			continue
		}

		log.Printf("kafka event received: topic=%s partition=%d offset=%d", msg.Topic, msg.Partition, msg.Offset)

		m.dispatcher.Dispatch(ctx, r.routeID, msg.Topic, string(msg.Value))
	}
}

func (m *ConsumerManager) loadActiveRoutes(ctx context.Context) ([]topicRoute, error) {
	query := `SELECT r.id, r.topic, kc.brokers, kc.group_id
			FROM routes r
			JOIN kafka_connections kc ON r.kafka_connection_id = kc.id
			WHERE r.active = true AND kc.active = true`

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []topicRoute
	for rows.Next() {
		var r topicRoute
		var brokersStr string

		if err := rows.Scan(&r.routeID, &r.topic, &brokersStr, &r.groupID); err != nil {
			return nil, err
		}

		for _, b := range strings.Split(brokersStr, ",") {
			if b = strings.TrimSpace(b); b != "" {
				r.brokers = append(r.brokers, b)
			}
		}

		routes = append(routes, r)
	}

	return routes, nil
}
