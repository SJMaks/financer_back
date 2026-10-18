package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SJMaks/financer_back/internal/features/statistics/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Repository interface {
	GetByUser(ctx context.Context, userID int64) ([]domain.UserStatistics, error)
	GetByAccount(ctx context.Context, accountID, userID int64) ([]domain.AccountStatistics, error)
}

type Service struct {
	repo Repository
	rds  *redis.Client
	log  *zap.SugaredLogger
}

func NewService(repo Repository, rds *redis.Client, log *zap.SugaredLogger) *Service {
	return &Service{
		repo: repo,
		rds:  rds,
		log:  log,
	}
}

func (s *Service) GetByUser(ctx context.Context, userID int64) ([]domain.UserStatistics, error) {
	rdsCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rdsKey := fmt.Sprintf("financer:v1:user:%d:statistics", userID)

	if res, err := s.rds.Get(rdsCtx, rdsKey).Result(); err == nil {
		var stats []domain.UserStatistics
		if err := json.Unmarshal([]byte(res), &stats); err == nil {
			return stats, nil
		}
	}

	stats, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if statsJSON, err := json.Marshal(stats); err == nil {
		if err := s.rds.Set(rdsCtx, rdsKey, statsJSON, time.Minute).Err(); err != nil {
			s.log.Warnf("failed to cache user %d statistics: %v", userID, err)
		}
	}

	return stats, nil
}

func (s *Service) GetByAccount(ctx context.Context, accountID, userID int64) ([]domain.AccountStatistics, error) {
	rdsCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rdsKey := fmt.Sprintf("financer:v1:user:%d:account:%d:statistics", userID, accountID)

	if res, err := s.rds.Get(rdsCtx, rdsKey).Result(); err == nil {
		var stats []domain.AccountStatistics
		if err := json.Unmarshal([]byte(res), &stats); err == nil {
			return stats, nil
		}
	}

	stats, err := s.repo.GetByAccount(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}

	if statsJSON, err := json.Marshal(stats); err == nil {
		if err := s.rds.Set(rdsCtx, rdsKey, statsJSON, time.Minute).Err(); err != nil {
			s.log.Warnf("failed to cache user %d account %d statistics: %v", userID, accountID, err)
		}
	}

	return stats, nil
}
