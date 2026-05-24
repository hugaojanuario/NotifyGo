package template

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	tmpl "html/template"

	"github.com/google/uuid"
)

type TemplateService struct {
	r TemplateRepositoryMethods
}

func NewTemplateService(r TemplateRepositoryMethods) *TemplateService {
	return &TemplateService{r: r}
}

func (s *TemplateService) Create(ctx context.Context, userID uuid.UUID, req CreateTemplate) (*Template, error) {
	template, err := s.r.CreateTemplate(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("service - error create template: %w", err)
	}

	return template, nil
}

func (s *TemplateService) GetAll(ctx context.Context, userID uuid.UUID) ([]Template, error) {
	templates, err := s.r.GetAllTemplates(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get all templates: %w", err)
	}

	return templates, nil
}

func (s *TemplateService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Template, error) {
	template, err := s.r.GetByIDTemplate(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get template by id: %w", err)
	}
	if template == nil {
		return nil, errors.New("template not found")
	}

	return template, nil
}

func (s *TemplateService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req UpdateTemplate) (*Template, error) {
	_, err := s.r.GetByIDTemplate(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service - error get template before update: %w", err)
	}

	template, err := s.r.UpdateByIDTemplate(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("service - error update template: %w", err)
	}

	return template, nil
}

func (s *TemplateService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := s.r.GetByIDTemplate(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("service - error get template before delete: %w", err)
	}

	err = s.r.DeleteByIDTemplate(ctx, id)
	if err != nil {
		return fmt.Errorf("service - error delete template: %w", err)
	}

	return nil
}

func (s *TemplateService) Preview(ctx context.Context, id uuid.UUID, userID uuid.UUID, data map[string]any) (string, error) {
	t, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return "", err
	}

	parsed, err := tmpl.New("preview").Parse(t.Content)
	if err != nil {
		return "", fmt.Errorf("service - error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := parsed.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("service - error rendering template: %w", err)
	}

	return buf.String(), nil
}
