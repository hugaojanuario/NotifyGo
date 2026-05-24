package notification

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hugaojanuario/NotifyGo/internal/channel"
)

type LogFilters struct {
	RouteID     *uuid.UUID
	Status      *NotificationStatus
	ChannelType *channel.ChannelType
	Page        int
	Limit       int
}

type CreateNotificationLog struct {
	RouteID         uuid.UUID
	ChannelConfigID uuid.UUID
	Topic           string
	Channel         channel.ChannelType
	Recipient       string
	Status          NotificationStatus
	Payload         string
	ErrorMessage    *string
	Attempts        int
	SentAt          *time.Time
}

type NotificationLogRepositoryMethods interface {
	Create(ctx context.Context, req CreateNotificationLog) (*NotificationLog, error)
	GetAll(ctx context.Context, userID uuid.UUID, filters LogFilters) ([]NotificationLog, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*NotificationLog, error)
	GetMetrics(ctx context.Context, userID uuid.UUID) (map[string]int, error)
}

type NotificationLogRepository struct {
	db *sql.DB
}

func NewNotificationLogRepository(db *sql.DB) *NotificationLogRepository {
	return &NotificationLogRepository{db: db}
}

func (r *NotificationLogRepository) Create(ctx context.Context, req CreateNotificationLog) (*NotificationLog, error) {
	query := `INSERT INTO notification_logs
			(route_id, channel_config_id, topic, channel, recipient, status, payload, error_message, attempts, sent_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, route_id, channel_config_id, topic, channel, recipient, status, payload, error_message, attempts, sent_at, created_at`

	log := &NotificationLog{}
	err := r.db.QueryRowContext(ctx, query,
		req.RouteID, req.ChannelConfigID, req.Topic, req.Channel, req.Recipient,
		req.Status, req.Payload, req.ErrorMessage, req.Attempts, req.SentAt,
	).Scan(
		&log.ID, &log.RouteID, &log.ChannelConfigID, &log.Topic, &log.Channel,
		&log.Recipient, &log.Status, &log.Payload, &log.ErrorMessage,
		&log.Attempts, &log.SentAt, &log.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("repository - error create notification log: %w", err)
	}

	return log, nil
}

func (r *NotificationLogRepository) GetAll(ctx context.Context, userID uuid.UUID, filters LogFilters) ([]NotificationLog, int, error) {
	if filters.Limit == 0 {
		filters.Limit = 20
	}
	if filters.Page == 0 {
		filters.Page = 1
	}
	offset := (filters.Page - 1) * filters.Limit

	args := []any{userID}
	where := `WHERE r.user_id = $1`

	if filters.RouteID != nil {
		args = append(args, *filters.RouteID)
		where += fmt.Sprintf(" AND nl.route_id = $%d", len(args))
	}
	if filters.Status != nil {
		args = append(args, *filters.Status)
		where += fmt.Sprintf(" AND nl.status = $%d", len(args))
	}
	if filters.ChannelType != nil {
		args = append(args, *filters.ChannelType)
		where += fmt.Sprintf(" AND nl.channel = $%d", len(args))
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM notification_logs nl
			JOIN routes r ON nl.route_id = r.id %s`, where)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("repository - error count notification logs: %w", err)
	}

	args = append(args, filters.Limit, offset)
	query := fmt.Sprintf(`SELECT nl.id, nl.route_id, nl.channel_config_id, nl.topic, nl.channel,
			nl.recipient, nl.status, nl.payload, nl.error_message, nl.attempts, nl.sent_at, nl.created_at
			FROM notification_logs nl
			JOIN routes r ON nl.route_id = r.id
			%s
			ORDER BY nl.created_at DESC
			LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository - error get all notification logs: %w", err)
	}
	defer rows.Close()

	var logs []NotificationLog
	for rows.Next() {
		var log NotificationLog
		err := rows.Scan(
			&log.ID, &log.RouteID, &log.ChannelConfigID, &log.Topic, &log.Channel,
			&log.Recipient, &log.Status, &log.Payload, &log.ErrorMessage,
			&log.Attempts, &log.SentAt, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("repository - error scanning notification log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

func (r *NotificationLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*NotificationLog, error) {
	query := `SELECT id, route_id, channel_config_id, topic, channel, recipient, status, payload, error_message, attempts, sent_at, created_at
			FROM notification_logs
			WHERE id = $1`

	log := &NotificationLog{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID, &log.RouteID, &log.ChannelConfigID, &log.Topic, &log.Channel,
		&log.Recipient, &log.Status, &log.Payload, &log.ErrorMessage,
		&log.Attempts, &log.SentAt, &log.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error get notification log by id: %w", err)
	}

	return log, nil
}

func (r *NotificationLogRepository) GetMetrics(ctx context.Context, userID uuid.UUID) (map[string]int, error) {
	query := `SELECT nl.status, COUNT(*) FROM notification_logs nl
			JOIN routes r ON nl.route_id = r.id
			WHERE r.user_id = $1
			GROUP BY nl.status`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repository - error get metrics: %w", err)
	}
	defer rows.Close()

	metrics := map[string]int{
		"total":    0,
		"success":  0,
		"failed":   0,
		"retrying": 0,
	}

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("repository - error scanning metrics: %w", err)
		}
		metrics["total"] += count
		switch NotificationStatus(status) {
		case StatusSuccess:
			metrics["success"] = count
		case StatusFailed:
			metrics["failed"] = count
		case StatusRetrying:
			metrics["retrying"] = count
		}
	}

	return metrics, nil
}
