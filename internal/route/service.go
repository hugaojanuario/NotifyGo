package route

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type RouteService struct {
	r RouteRepositoryMethods
}

func NewRouteService(r RouteRepositoryMethods) *RouteService {
	return &RouteService{r: r}
}

func (s *RouteService) CreateRoute(ctx context.Context, userID uuid.UUID, req CreateRoute) (*Route, error) {
	route, err := s.r.Create(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("service - error create new route: %w", err)
	}

	return route, nil
}

func (s *RouteService) GetAll(ctx context.Context, userID uuid.UUID) ([]Route, error) {
	routes, err := s.r.GetAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get all routes: %w", err)
	}

	return routes, nil
}

func (s *RouteService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Route, error) {
	route, err := s.r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get route by id: %w", err)
	}
	if route == nil {
		return nil, errors.New("route not found")
	}

	return route, nil
}

func (s *RouteService) UpdateRoute(ctx context.Context, id uuid.UUID, userID uuid.UUID, req UpdateRoute) (*Route, error) {
	_, err := s.r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get route before update: %w", err)
	}

	route, err := s.r.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("service - error update route: %w", err)
	}

	return route, nil
}

func (s *RouteService) DeleteRoute(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := s.r.GetByID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("service - error get route before delete: %w", err)
	}

	err = s.r.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service - error delete route: %w", err)
	}

	return nil
}

func (s *RouteService) ToggleActive(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Route, error) {
	_, err := s.r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get route before toggle: %w", err)
	}

	route, err := s.r.ToggleActive(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service - error toggle route: %w", err)
	}

	return route, nil
}
