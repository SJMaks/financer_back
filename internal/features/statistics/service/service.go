package service

import (
	"context"

	"github.com/SJMaks/financer_back/internal/features/statistics/domain"
)

type Repository interface {
	GetByUser(ctx context.Context, userID int64) ([]domain.UserStatistics, error)
	GetByAccount(ctx context.Context, accountID, userID int64) ([]domain.AccountStatistics, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByUser(ctx context.Context, userID int64) ([]domain.UserStatistics, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *Service) GetByAccount(ctx context.Context, accountID, userID int64) ([]domain.AccountStatistics, error) {
	return s.repo.GetByAccount(ctx, accountID, userID)
}
