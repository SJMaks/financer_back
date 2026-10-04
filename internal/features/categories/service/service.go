package service

import (
	"context"
	"strings"

	"github.com/SJMaks/financer_back/internal/features/categories/domain"
)

type Repository interface {
	Create(ctx context.Context, a *domain.Category) (*domain.Category, error)
	GetAllOfUser(ctx context.Context, userID int64) ([]domain.Category, error)
	GetByID(ctx context.Context, id, userID int64) (*domain.Category, error)
	UpdateByID(ctx context.Context, a *domain.PatchCategory) (*domain.Category, error)
	DeleteByID(ctx context.Context, id, userID int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, c *domain.Category) (*domain.Category, error) {
	if len(strings.TrimSpace(c.Title)) < 1 {
		return nil, EmptyTitle
	}
	if c.OperationTypeID == 0 {
		return nil, MissingOperationTypeID
	}
	return s.repo.Create(ctx, c)
}

func (s *Service) GetAllOfUser(ctx context.Context, userID int64) ([]domain.Category, error) {
	return s.repo.GetAllOfUser(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id, userID int64) (*domain.Category, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *Service) UpdateByID(ctx context.Context, c *domain.PatchCategory) (*domain.Category, error) {
	if c.Title != nil && len(strings.TrimSpace(*c.Title)) < 1 {
		return nil, EmptyTitle
	}
	if c.OperationTypeID != nil && *c.OperationTypeID == 0 {
		return nil, MissingOperationTypeID
	}
	return s.repo.UpdateByID(ctx, c)
}

func (s *Service) DeleteByID(ctx context.Context, id, userID int64) error {
	return s.repo.DeleteByID(ctx, id, userID)
}
