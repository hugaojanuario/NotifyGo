package connection

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
)

type KafkaConnectionService struct {
	r KafkaRepositoryMethods
}

func NewKafkaConnectionService(r KafkaRepositoryMethods) *KafkaConnectionService {
	return &KafkaConnectionService{r: r}
}

func (s *KafkaConnectionService) Create(ctx context.Context, userID uuid.UUID, req CreateKafkaConnection) (*KafkaConnection, error) {
	conn, err := s.r.Create(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("service - error create kafka connection: %w", err)
	}

	return conn, nil
}

func (s *KafkaConnectionService) GetAll(ctx context.Context, userID uuid.UUID) ([]KafkaConnection, error) {
	connections, err := s.r.GetAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get all kafka connections: %w", err)
	}

	return connections, nil
}

func (s *KafkaConnectionService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*KafkaConnection, error) {
	conn, err := s.r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get kafka connection by id: %w", err)
	}
	if conn == nil {
		return nil, errors.New("kafka connection not found")
	}

	return conn, nil
}

func (s *KafkaConnectionService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req UpdateKafkaConnection) (*KafkaConnection, error) {
	_, err := s.r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get kafka connection before update: %w", err)
	}

	conn, err := s.r.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("service - error update kafka connection: %w", err)
	}

	return conn, nil
}

func (s *KafkaConnectionService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := s.r.GetByID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("service - error get kafka connection before delete: %w", err)
	}

	err = s.r.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service - error delete kafka connection: %w", err)
	}

	return nil
}

func (s *KafkaConnectionService) TestConnection(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*KafkaConnection, error) {
	conn, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	brokers := strings.Split(conn.Brokers, ",")
	for _, broker := range brokers {
		broker = strings.TrimSpace(broker)
		tcpConn, err := net.DialTimeout("tcp", broker, 5*time.Second)
		if err != nil {
			return nil, fmt.Errorf("service - could not connect to broker %s: %w", broker, err)
		}
		tcpConn.Close()
	}

	return conn, nil
}
