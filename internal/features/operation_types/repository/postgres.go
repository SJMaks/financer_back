package repository

import (
	"context"

	"github.com/SJMaks/financer_back/internal/features/operation_types/domain"
	"github.com/SJMaks/financer_back/internal/features/operation_types/service"
	"github.com/jackc/pgx/v5"
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

func (r *PostgresRepository) GetAll(ctx context.Context) ([]domain.OperationType, error) {
	qry := `SELECT id, title FROM operation_types;`
	rows, err := r.db.Query(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.OperationType])
}
