package route

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type RouteMethods interface {
	CreateRoute(ctx context.Context, request CreateRoute) (*Route, error)
	GetAll(ctx context.Context) ([]Route, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Route, error)
	Update(ctx context.Context, id uuid.UUID, request UpdateRoute) (*Route, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) RouteRepository {
	return RouteRepository{db: db}
}

func (r *RouteRepository) CreateRouter(ctx context.Context, request CreateRoute) (*Route, error) {
	query := `INSERT INTO routes (kafka_connection_id, name, topic)
			VALUES ($1, $2, $3)
			RETURNING id, user_id, kafka_connection_id, name, topic, active, created_at`

	route := &Route{}
	err := r.db.QueryRowContext(ctx, query, request.KafkaConnectionID, request.Name, request.Topic).
		Scan(&route.ID, &route.UserID, &route.KafkaConnectionID, &route.Name, &route.Topic, &route.Active, &route.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository - error create new route: %w", err)
	}

	return route, nil
}

func (r *RouteRepository) GetAll(ctx context.Context) ([]Route, error) {
	query := `SELECT id, name, topic, active, created_at, updated_at 
			FROM routes
			ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository - error get all routes: %w", err)
	}
	defer rows.Close()
	var routes []Route

	for rows.Next() {
		var route Route
		err := rows.Scan(&route.ID, &route.UserID, &route.KafkaConnectionID, &route.Name, &route.Topic, &route.Active, &route.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("repository - error get all routes: %w", err)
		}

		routes = append(routes, route)
	}

	return routes, nil
}

func (r *RouteRepository) GetByID(ctx context.Context, id uuid.UUID) (*Route, error) {
	query := `SELECT id, user_id, kafka_connection_id, name, topic, active, created_at
			FROM routes
			WHERE id = $1`

	route := &Route{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&route.ID, &route.UserID, &route.KafkaConnectionID, &route.Name, &route.Topic, &route.Active, &route.CreatedAt)
	if err != sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error get route by id: %w", err)
	}

	return route, nil
}

func (r *RouteRepository) Update(ctx context.Context, id uuid.UUID, request UpdateRoute) (*Route, error) {
	query := `UPDATE routes
			SET kafka_connection_id = $1,name = $2, topic = $3
			WHERE id = $4
			RETURNING id, user_id, kafka_connection_id, name, topic, active, created_at`

	route := &Route{}
	err := r.db.QueryRowContext(ctx, query, request.Name, request.Topic, id).
		Scan(&route.KafkaConnectionID, &route.ID, &route.UserID)
	if err != sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error update route: %w", err)
	}
	return route, nil
}

func (r *RouteRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET active = true WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository - error soft delete route: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
