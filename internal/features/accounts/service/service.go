package service

import (
	"context"
	"strings"

	"github.com/SJMaks/financer_back/internal/features/accounts/domain"
)

type Repository interface {
	Create(ctx context.Context, a *domain.Account) (*domain.Account, error)
	GetAllOfUser(ctx context.Context, userID int64) ([]domain.Account, error)
	GetByID(ctx context.Context, id, userID int64) (*domain.Account, error)
	UpdateByID(ctx context.Context, a *domain.Account) (*domain.Account, error)
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

func (s *Service) Create(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	if len(strings.TrimSpace(a.Title)) < 1 {
		return nil, EmptyTitle
	}
	if a.Balance < 0 {
		return nil, NegativeBalance
	}
	return s.repo.Create(ctx, a)
}

func (s *Service) GetAllOfUser(ctx context.Context, userID int64) ([]domain.Account, error) {
	return s.repo.GetAllOfUser(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id, userID int64) (*domain.Account, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *Service) UpdateByID(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	if len(strings.TrimSpace(a.Title)) < 1 {
		return nil, EmptyTitle
	}
	if a.Balance < 0 {
		return nil, NegativeBalance
	}
	return s.repo.UpdateByID(ctx, a)
}

func (s *Service) DeleteByID(ctx context.Context, id, userID int64) error {
	return s.repo.DeleteByID(ctx, id, userID)
}
