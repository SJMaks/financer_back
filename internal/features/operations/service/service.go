package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/SJMaks/financer_back/internal/features/operations/domain"
)

type Repository interface {
	Create(ctx context.Context, o *domain.Operation) (*domain.Operation, error)
	GetAllOfAccount(ctx context.Context, params ListOperationsParams) ([]domain.Operation, int64, error)
	GetByID(ctx context.Context, id, userID, accountID int64) (*domain.Operation, error)
	UpdateByID(ctx context.Context, o *domain.PatchOperation) (*domain.Operation, error)
	DeleteByID(ctx context.Context, id, userID, accountID int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, o *domain.Operation) (*domain.Operation, error) {
	if o.AccountID == 0 {
		return nil, AccountIDRequired
	}
	if o.CategoryID == 0 {
		return nil, CategoryIDRequired
	}

	title := strings.TrimSpace(o.Title)
	o.Title = title
	var description string
	if o.Description != nil {
		description = strings.TrimSpace(*o.Description)
		o.Description = &description
	}

	if title == "" {
		return nil, TitleRequired
	}
	if utf8.RuneCountInString(title) > 50 {
		return nil, InvalidTitle
	}

	if utf8.RuneCountInString(description) > 500 {
		return nil, InvalidDescription
	}

	if o.ExecutedAt.IsZero() {
		o.ExecutedAt = time.Now()
	}
	if o.Sum < 0 {
		return nil, InvalidSum
	}
	return s.repo.Create(ctx, o)
}

func (s *Service) GetAllOfAccount(ctx context.Context, params ListOperationsParams) ([]domain.Operation, int64, error) {
	return s.repo.GetAllOfAccount(ctx, params)
}

func (s *Service) GetByID(ctx context.Context, id, userID, accountID int64) (*domain.Operation, error) {
	return s.repo.GetByID(ctx, id, userID, accountID)
}

func (s *Service) UpdateByID(ctx context.Context, o *domain.PatchOperation) (*domain.Operation, error) {
	if o.Version == nil {
		return nil, MissingVersion
	}

	if o.CategoryID != nil && *o.CategoryID == 0 {
		return nil, CategoryIDRequired
	}

	var title, description string
	if o.Title != nil {
		title = strings.TrimSpace(*o.Title)
		o.Title = &title
	}
	if o.Description != nil {
		description = strings.TrimSpace(*o.Description)
		o.Description = &description
	}

	if utf8.RuneCountInString(title) > 50 {
		return nil, InvalidTitle
	}
	if utf8.RuneCountInString(description) > 500 {
		return nil, InvalidDescription
	}
	if o.Sum != nil && *o.Sum < 0 {
		return nil, InvalidSum
	}
	if o.ExecutedAt != nil && o.ExecutedAt.IsZero() {
		return nil, InvalidExecAt
	}
	return s.repo.UpdateByID(ctx, o)
}

func (s *Service) DeleteByID(ctx context.Context, id, userID, accountID int64) error {
	return s.repo.DeleteByID(ctx, id, userID, accountID)
}
