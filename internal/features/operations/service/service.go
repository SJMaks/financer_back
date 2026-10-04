package service

import (
	"context"
	"strings"
	"time"

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
		return nil, MissingAccountID
	}
	if o.CategoryID == 0 {
		return nil, MissingCategoryID
	}
	if len(strings.TrimSpace(o.Title)) < 1 {
		return nil, EmptyTitle
	}
	if o.ExecutedAt.IsZero() {
		o.ExecutedAt = time.Now()
	}
	if o.Sum < 0 {
		return nil, NegativeSum
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
	if o.CategoryID != nil && *o.CategoryID == 0 {
		return nil, MissingCategoryID
	}
	if o.Title != nil && len(strings.TrimSpace(*o.Title)) < 1 {
		return nil, EmptyTitle
	}
	if o.ExecutedAt != nil && o.ExecutedAt.IsZero() {
		now := time.Now()
		o.ExecutedAt = &now
	}
	if o.Sum != nil && *o.Sum < 0 {
		return nil, NegativeSum
	}
	return s.repo.UpdateByID(ctx, o)
}

func (s *Service) DeleteByID(ctx context.Context, id, userID, accountID int64) error {
	return s.repo.DeleteByID(ctx, id, userID, accountID)
}
