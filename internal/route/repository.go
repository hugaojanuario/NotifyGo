package route

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type RouteRepositoryMethods interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateRoute) (*Route, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]Route, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Route, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRoute) (*Route, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ToggleActive(ctx context.Context, id uuid.UUID) (*Route, error)
}

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) Create(ctx context.Context, userID uuid.UUID, req CreateRoute) (*Route, error) {
	query := `INSERT INTO routes (user_id, kafka_connection_id, name, topic)
			VALUES ($1, $2, $3, $4)
			RETURNING id, user_id, kafka_connection_id, name, topic, active, created_at`

	route := &Route{}
	err := r.db.QueryRowContext(ctx, query, userID, req.KafkaConnectionID, req.Name, req.Topic).
		Scan(&route.ID, &route.UserID, &route.KafkaConnectionID, &route.Name, &route.Topic, &route.Active, &route.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository - error create new route: %w", err)
	}

	return route, nil
}

func (r *RouteRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]Route, error) {
	query := `SELECT id, user_id, kafka_connection_id, name, topic, active, created_at
			FROM routes
			WHERE user_id = $1
			ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repository - error get all routes: %w", err)
	}
	defer rows.Close()

	var routes []Route
	for rows.Next() {
		var route Route
		err := rows.Scan(&route.ID, &route.UserID, &route.KafkaConnectionID, &route.Name, &route.Topic, &route.Active, &route.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("repository - error scanning route: %w", err)
		}
		routes = append(routes, route)
	}

	return routes, nil
}

func (r *RouteRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Route, error) {
	query := `SELECT id, user_id, kafka_connection_id, name, topic, active, created_at
			FROM routes
			WHERE id = $1 AND user_id = $2`

	route := &Route{}
	err := r.db.QueryRowContext(ctx, query, id, userID).
		Scan(&route.ID, &route.UserID, &route.KafkaConnectionID, &route.Name, &route.Topic, &route.Active, &route.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error get route by id: %w", err)
	}

	return route, nil
}

func (r *RouteRepository) Update(ctx context.Context, id uuid.UUID, req UpdateRoute) (*Route, error) {
	query := `UPDATE routes
			SET kafka_connection_id = $1, name = $2, topic = $3
			WHERE id = $4
			RETURNING id, user_id, kafka_connection_id, name, topic, active, created_at`

	route := &Route{}
	err := r.db.QueryRowContext(ctx, query, req.KafkaConnectionID, req.Name, req.Topic, id).
		Scan(&route.ID, &route.UserID, &route.KafkaConnectionID, &route.Name, &route.Topic, &route.Active, &route.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error update route: %w", err)
	}

	return route, nil
}

func (r *RouteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM routes WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository - error delete route: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *RouteRepository) ToggleActive(ctx context.Context, id uuid.UUID) (*Route, error) {
	query := `UPDATE routes
			SET active = NOT active
			WHERE id = $1
			RETURNING id, user_id, kafka_connection_id, name, topic, active, created_at`

	route := &Route{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&route.ID, &route.UserID, &route.KafkaConnectionID, &route.Name, &route.Topic, &route.Active, &route.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error toggle route active: %w", err)
	}

	return route, nil
}
