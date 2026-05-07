package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	r UserRepositoryMethods
}

func NewUserService(r UserRepositoryMethods) *UserService {
	return &UserService{r: r}
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("service - error generating password hash: %w", err)
	}

	user, err := s.r.CreateUser(ctx, req, string(passwordHash))
	if err != nil {
		return nil, fmt.Errorf("service - error create new user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]UserResponse, error) {
	return s.r.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service - error when listing the alert by ID: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserResponse, error) {
	_, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service - database error when fetching user: %w", err)
	}
	return s.r.Update(ctx, id, req)
}

func (s *UserService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := s.r.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository - database error when fetching user: %w", err)
	}

	return s.r.SoftDelete(ctx, id)
}
