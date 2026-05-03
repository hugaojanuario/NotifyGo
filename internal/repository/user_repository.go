package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hugaojanuario/NotifyGo/internal/domain"
)

type UserRepositoryMethods interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest, passwordHash string) (*domain.UserResponse, error)
	GetAll(ctx context.Context) ([]domain.UserResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, req domain.CreateUserRequest, passwordHash string) (*domain.UserResponse, error) {
	query := `INSERT INTO users (name, email, password_hash)
			VALUES ($1, $2, $3)
			RETURNING id, name, email, active, created_at, updated_at`

	user := &domain.UserResponse{}
	err := r.db.QueryRowContext(ctx, query, req.Name, req.Email, passwordHash).
		Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository - error create new user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]domain.UserResponse, error) {
	query := `SELECT id, name, email, active, created_at, updated_at
			FROM users
			WHERE active = true
			ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository - error this listing the users: %w", err)
	}
	defer rows.Close()

	var users []domain.UserResponse

	for rows.Next() {
		var user domain.UserResponse

		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("repository - error when listing alerts: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	query := `SELECT id, name, email, active, created_at, updated_at
			FROM users
			WHERE ID = $1`

	user := &domain.UserResponse{}

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error when listing the alert by ID: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	query := `UPDATE users
			SET name = $1, email = $2, password_hash = $3
			WHERE id = $4
			RETURNING id, name, email, active, created_at, updated_at`

	user := &domain.UserResponse{}

	err := r.db.QueryRowContext(ctx, query, req.Name, req.Email, req.Password, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - repository: error when updating user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET active = true WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository - error soft delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // não encontrado
	}

	return nil
}
