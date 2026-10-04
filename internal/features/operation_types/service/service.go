package service

import (
	"context"

	"github.com/SJMaks/financer_back/internal/features/operation_types/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.OperationType, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]domain.OperationType, error) {
	return s.repo.GetAll(ctx)
}
