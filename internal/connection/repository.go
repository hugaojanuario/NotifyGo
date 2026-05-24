package connection

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type KafkaRepositoryMethods interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateKafkaConnection) (*KafkaConnection, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]KafkaConnection, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*KafkaConnection, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateKafkaConnection) (*KafkaConnection, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type KafkaRepository struct {
	db *sql.DB
}

func NewKafkaRepository(db *sql.DB) *KafkaRepository {
	return &KafkaRepository{db: db}
}

func (r *KafkaRepository) Create(ctx context.Context, userID uuid.UUID, req CreateKafkaConnection) (*KafkaConnection, error) {
	query := `INSERT INTO kafka_connections (user_id, name, brokers, group_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id, user_id, name, brokers, group_id, active, created_at`

	conn := &KafkaConnection{}
	err := r.db.QueryRowContext(ctx, query, userID, req.Name, req.Brokers, req.GroupID).
		Scan(&conn.ID, &conn.UserID, &conn.Name, &conn.Brokers, &conn.GroupID, &conn.Active, &conn.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository - error create kafka connection: %w", err)
	}

	return conn, nil
}

func (r *KafkaRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]KafkaConnection, error) {
	query := `SELECT id, user_id, name, brokers, group_id, active, created_at
			FROM kafka_connections
			WHERE user_id = $1
			ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repository - error get all kafka connections: %w", err)
	}
	defer rows.Close()

	var connections []KafkaConnection
	for rows.Next() {
		var conn KafkaConnection
		err := rows.Scan(&conn.ID, &conn.UserID, &conn.Name, &conn.Brokers, &conn.GroupID, &conn.Active, &conn.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("repository - error scanning kafka connection: %w", err)
		}
		connections = append(connections, conn)
	}

	return connections, nil
}

func (r *KafkaRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*KafkaConnection, error) {
	query := `SELECT id, user_id, name, brokers, group_id, active, created_at
			FROM kafka_connections
			WHERE id = $1 AND user_id = $2`

	conn := &KafkaConnection{}
	err := r.db.QueryRowContext(ctx, query, id, userID).
		Scan(&conn.ID, &conn.UserID, &conn.Name, &conn.Brokers, &conn.GroupID, &conn.Active, &conn.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error get kafka connection by id: %w", err)
	}

	return conn, nil
}

func (r *KafkaRepository) Update(ctx context.Context, id uuid.UUID, req UpdateKafkaConnection) (*KafkaConnection, error) {
	query := `UPDATE kafka_connections
			SET name = $1, brokers = $2, group_id = $3
			WHERE id = $4
			RETURNING id, user_id, name, brokers, group_id, active, created_at`

	conn := &KafkaConnection{}
	err := r.db.QueryRowContext(ctx, query, req.Name, req.Brokers, req.GroupID, id).
		Scan(&conn.ID, &conn.UserID, &conn.Name, &conn.Brokers, &conn.GroupID, &conn.Active, &conn.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error update kafka connection: %w", err)
	}

	return conn, nil
}

func (r *KafkaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM kafka_connections WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository - error delete kafka connection: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
