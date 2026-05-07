package route

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type RouteService struct {
	r RouteMethods
}

func NewRouteService(r RouteMethods) *RouteService {
	return &RouteService{r: r}
}

func (s *RouteService) CreateRoute(ctx context.Context, request CreateRoute) (*Route, error) {
	route, err := s.r.CreateRoute(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("ssrvice - error create new route: %w", err)
	}
	return route, nil
}

func (s *RouteService) GetAll(ctx context.Context) ([]Route, error) {
	return s.r.GetAll(ctx)
}

func (s *RouteService) GetByID(ctx context.Context, id uuid.UUID) (*Route, error) {
	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ssrvice - error find route: %w", err)
	}
	if user == nil {
		return nil, errors.New("route not found")
	}

	return user, nil
}

func (s *RouteService) UpdateRoute(ctx context.Context, id uuid.UUID, request UpdateRoute) (*Route, error) {
	_, err := s.r.Update(ctx, id, request)
	if err != nil {
		return nil, fmt.Errorf("ssrvice - error update route: %w", err)
	}

	return s.r.Update(ctx, id, request)
}

func (s *RouteService) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	err := s.r.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("ssrvice - error delete route: %w", err)
	}
	return nil
}
