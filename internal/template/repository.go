package template

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type TemplateRepositoryMethods interface {
	CreateTemplate(ctx context.Context, userID uuid.UUID, req CreateTemplate) (*Template, error)
	GetAllTemplates(ctx context.Context, userID uuid.UUID) ([]Template, error)
	GetByIDTemplate(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Template, error)
	UpdateByIDTemplate(ctx context.Context, id uuid.UUID, req UpdateTemplate) (*Template, error)
	DeleteByIDTemplate(ctx context.Context, id uuid.UUID) error
}

type TemplateRepository struct {
	db *sql.DB
}

func NewTemplateRepository(db *sql.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) CreateTemplate(ctx context.Context, userID uuid.UUID, req CreateTemplate) (*Template, error) {
	query := `INSERT INTO templates (user_id, name, content)
			VALUES ($1, $2, $3)
			RETURNING id, user_id, name, content, created_at, updated_at`

	template := &Template{}
	err := r.db.QueryRowContext(ctx, query, userID, req.Name, req.Content).
		Scan(&template.ID, &template.UserID, &template.Name, &template.Content, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository - error create template: %w", err)
	}

	return template, nil
}

func (r *TemplateRepository) GetAllTemplates(ctx context.Context, userID uuid.UUID) ([]Template, error) {
	query := `SELECT id, user_id, name, content, created_at, updated_at
			FROM templates
			WHERE user_id = $1
			ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repository - error get all templates: %w", err)
	}
	defer rows.Close()

	var templates []Template
	for rows.Next() {
		var template Template
		err = rows.Scan(&template.ID, &template.UserID, &template.Name, &template.Content, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("repository - error scanning template: %w", err)
		}
		templates = append(templates, template)
	}

	return templates, nil
}

func (r *TemplateRepository) GetByIDTemplate(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Template, error) {
	query := `SELECT id, user_id, name, content, created_at, updated_at
			FROM templates
			WHERE id = $1 AND user_id = $2`

	template := &Template{}
	err := r.db.QueryRowContext(ctx, query, id, userID).
		Scan(&template.ID, &template.UserID, &template.Name, &template.Content, &template.CreatedAt, &template.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error get template by id: %w", err)
	}

	return template, nil
}

func (r *TemplateRepository) UpdateByIDTemplate(ctx context.Context, id uuid.UUID, req UpdateTemplate) (*Template, error) {
	query := `UPDATE templates
			SET name = $1, content = $2
			WHERE id = $3
			RETURNING id, user_id, name, content, created_at, updated_at`

	template := &Template{}
	err := r.db.QueryRowContext(ctx, query, req.Name, req.Content, id).
		Scan(&template.ID, &template.UserID, &template.Name, &template.Content, &template.CreatedAt, &template.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository - error update template: %w", err)
	}

	return template, nil
}

func (r *TemplateRepository) DeleteByIDTemplate(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM templates WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository - error delete template: %w", err)
	}

	return nil
}
