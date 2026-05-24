package channel

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ChannelConfigService struct {
	r ChannelConfigRepositoryMethods
}

func NewChannelConfigService(r ChannelConfigRepositoryMethods) *ChannelConfigService {
	return &ChannelConfigService{r: r}
}

func (s *ChannelConfigService) Create(ctx context.Context, routeID uuid.UUID, req CreateChannelConfig) (*ChannelConfig, error) {
	if err := validateChannelConfig(req); err != nil {
		return nil, err
	}

	cfg, err := s.r.Create(ctx, routeID, req)
	if err != nil {
		return nil, fmt.Errorf("service - error create channel config: %w", err)
	}

	return cfg, nil
}

func (s *ChannelConfigService) GetAllByRouteID(ctx context.Context, routeID uuid.UUID) ([]ChannelConfig, error) {
	configs, err := s.r.GetAllByRouteID(ctx, routeID)
	if err != nil {
		return nil, fmt.Errorf("service - error get all channel configs: %w", err)
	}

	return configs, nil
}

func (s *ChannelConfigService) GetByID(ctx context.Context, id uuid.UUID, routeID uuid.UUID) (*ChannelConfig, error) {
	cfg, err := s.r.GetByID(ctx, id, routeID)
	if err != nil {
		return nil, fmt.Errorf("service - error get channel config by id: %w", err)
	}
	if cfg == nil {
		return nil, errors.New("channel config not found")
	}

	return cfg, nil
}

func (s *ChannelConfigService) Update(ctx context.Context, id uuid.UUID, routeID uuid.UUID, req CreateChannelConfig) (*ChannelConfig, error) {
	_, err := s.r.GetByID(ctx, id, routeID)
	if err != nil {
		return nil, fmt.Errorf("service - error get channel config before update: %w", err)
	}

	if err := validateChannelConfig(req); err != nil {
		return nil, err
	}

	cfg, err := s.r.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("service - error update channel config: %w", err)
	}

	return cfg, nil
}

func (s *ChannelConfigService) Delete(ctx context.Context, id uuid.UUID, routeID uuid.UUID) error {
	_, err := s.r.GetByID(ctx, id, routeID)
	if err != nil {
		return fmt.Errorf("service - error get channel config before delete: %w", err)
	}

	err = s.r.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service - error delete channel config: %w", err)
	}

	return nil
}

func validateChannelConfig(req CreateChannelConfig) error {
	switch req.ChannelType {
	case ChannelEmail:
		if req.ToField == nil && req.ToFixed == nil {
			return errors.New("email channel requires to_field or to_fixed")
		}
		if req.Subject == nil {
			return errors.New("email channel requires subject")
		}
	case ChannelWebhook:
		if req.WebhookURL == nil {
			return errors.New("webhook channel requires webhook_url")
		}
	case ChannelSlack:
		if req.SlackChannel == nil {
			return errors.New("slack channel requires slack_channel")
		}
		if req.WebhookURL == nil {
			return errors.New("slack channel requires webhook_url")
		}
	case ChannelSMS:
		if req.ToField == nil && req.ToFixed == nil {
			return errors.New("sms channel requires to_field or to_fixed")
		}
	default:
		return errors.New("unknown channel type")
	}

	return nil
}
