package service

import (
	"context"
	"fmt"

	"github.com/hugaojanuario/NotifyGo/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest, passwordHash string) (*domain.UserResponse, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar o hash da senha: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, req, string(passwordHash))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar o novo usuario: %w", err)
	}

	return user, nil
}
