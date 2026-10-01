package database

import (
	"context"
	"time"

	"github.com/SJMaks/financer_back/internal/core/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(cfg.DBConnString)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = 20
	conf.MinConns = 5
	conf.MaxConnLifetime = time.Hour
	conf.MaxConnIdleTime = 30 * time.Minute
	conf.HealthCheckPeriod = time.Minute

	return pgxpool.NewWithConfig(ctx, conf)
}
