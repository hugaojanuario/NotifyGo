package channel

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ChannelConfigRepositoryMethods interface {
	Create(ctx context.Context, routeID uuid.UUID, req CreateChannelConfig) (*ChannelConfig, error)
	GetAllByRouteID(ctx context.Context, routeID uuid.UUID) ([]ChannelConfig, error)
	GetByID(ctx context.Context, id uuid.UUID, routeID uuid.UUID) (*ChannelConfig, error)
	Update(ctx context.Context, id uuid.UUID, req CreateChannelConfig) (*ChannelConfig, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChannelConfigRepository struct {
	db *sql.DB
}

func NewChannelConfigRepository(db *sql.DB) *ChannelConfigRepository {
	return &ChannelConfigRepository{db: db}
}

func (r *ChannelConfigRepository) Create(ctx context.Context, routeID uuid.UUID, req CreateChannelConfig) (*ChannelConfig, error) {
	query := `INSERT INTO channel_configs
			(route_id, channel_type, to_field, to_fixed, subject, template_id, message_template, webhook_url, webhook_secret, slack_channel, active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, route_id, channel_type, to_field, to_fixed, subject, template_id, message_template, webhook_url, webhook_secret, slack_channel, active, created_at`

	cfg := &ChannelConfig{}
	err := r.db.QueryRowContext(ctx, query,
		routeID, req.ChannelType, req.ToField, req.ToFixed, req.Subject,
		req.TemplateID, req.MessageTemplate, req.WebhookURL, req.WebhookSecret,
		req.SlackChannel, req.Active,
	).Scan(
		&cfg.ID, &cfg.RouteID, &cfg.ChannelType, &cfg.ToField, &cfg.ToFixed,
		&cfg.Subject, &cfg.TemplateID, &cfg.MessageTemplate, &cfg.WebhookURL,
		&cfg.WebhookSecret, &cfg.SlackChannel, &cfg.Active, &cfg.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("repository - error create channel config: %w", err)
	}

	return cfg, nil
}

func (r *ChannelConfigRepository) GetAllByRouteID(ctx context.Context, routeID uuid.UUID) ([]ChannelConfig, error) {
	query := `SELECT id, route_id, channel_type, to_field, to_fixed, subject, template_id, message_template, webhook_url, webhook_secret, slack_channel, active, created_at
			FROM channel_configs
			WHERE route_id = $1
			ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, routeID)
	if err != nil {
		return nil, fmt.Errorf("repository - error get all channel configs: %w", err)
	}
	defer rows.Close()

	var configs []ChannelConfig
	for rows.Next() {
		var cfg ChannelConfig
		err := rows.Scan(
			&cfg.ID, &cfg.RouteID, &cfg.ChannelType, &cfg.ToField, &cfg.ToFixed,
			&cfg.Subject, &cfg.TemplateID, &cfg.MessageTemplate, &cfg.WebhookURL,
			&cfg.WebhookSecret, &cfg.SlackChannel, &cfg.Active, &cfg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository - error scanning channel config: %w", err)
		}
		configs = append(configs, cfg)
	}

	return configs, nil
}

func (r *ChannelConfigRepository) GetByID(ctx context.Context, id uuid.UUID, routeID uuid.UUID) (*ChannelConfig, error) {
	query := `SELECT id, route_id, channel_type, to_field, to_fixed, subject, template_id, message_template, webhook_url, webhook_secret, slack_channel, active, created_at
			FROM channel_configs
			WHERE id = $1 AND route_id = $2`

	cfg := &ChannelConfig{}
	err := r.db.QueryRowContext(ctx, query, id, routeID).Scan(
		&cfg.ID, &cfg.RouteID, &cfg.ChannelType, &cfg.ToField, &cfg.ToFixed,
		&cfg.Subject, &cfg.TemplateID, &cfg.MessageTemplate, &cfg.WebhookURL,
		&cfg.WebhookSecret, &cfg.SlackChannel, &cfg.Active, &cfg.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error get channel config by id: %w", err)
	}

	return cfg, nil
}

func (r *ChannelConfigRepository) Update(ctx context.Context, id uuid.UUID, req CreateChannelConfig) (*ChannelConfig, error) {
	query := `UPDATE channel_configs
			SET channel_type = $1, to_field = $2, to_fixed = $3, subject = $4, template_id = $5,
			    message_template = $6, webhook_url = $7, webhook_secret = $8, slack_channel = $9, active = $10
			WHERE id = $11
			RETURNING id, route_id, channel_type, to_field, to_fixed, subject, template_id, message_template, webhook_url, webhook_secret, slack_channel, active, created_at`

	cfg := &ChannelConfig{}
	err := r.db.QueryRowContext(ctx, query,
		req.ChannelType, req.ToField, req.ToFixed, req.Subject, req.TemplateID,
		req.MessageTemplate, req.WebhookURL, req.WebhookSecret, req.SlackChannel, req.Active,
		id,
	).Scan(
		&cfg.ID, &cfg.RouteID, &cfg.ChannelType, &cfg.ToField, &cfg.ToFixed,
		&cfg.Subject, &cfg.TemplateID, &cfg.MessageTemplate, &cfg.WebhookURL,
		&cfg.WebhookSecret, &cfg.SlackChannel, &cfg.Active, &cfg.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error update channel config: %w", err)
	}

	return cfg, nil
}

func (r *ChannelConfigRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM channel_configs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository - error delete channel config: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
