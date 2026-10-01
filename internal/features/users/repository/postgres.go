package repository

import (
	"github.com/SJMaks/financer_back/internal/features/users/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgres(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

var _ service.Repository = (*PostgresRepository)(nil)

func (r *PostgresRepository) Ping() (string, error) {
	return "Hello, World!", nil
}
