package repository

import (
	"context"
	"errors"

	"github.com/SJMaks/financer_back/internal/features/statistics/domain"
	"github.com/SJMaks/financer_back/internal/features/statistics/service"
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

func (r *PostgresRepository) GetByUser(ctx context.Context, userID int64) ([]domain.UserStatistics, error) {
	qry := `SELECT user_id, month, income, expense, profit FROM v_user_monthly_stats
	WHERE user_id = $1 ORDER BY month ASC;`
	rows, err := r.db.Query(ctx, qry, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.UserStatistics])
}

func (r *PostgresRepository) GetByAccount(ctx context.Context, accountID, userID int64) ([]domain.AccountStatistics, error) {
	var exists int
	err := r.db.QueryRow(ctx,
		`SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2`,
		accountID, userID,
	).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.AccountNotFound
		}
		return nil, err
	}
	qry := `SELECT account_id, user_id, month, income, expense, profit FROM v_account_monthly_stats
	WHERE user_id = $1 AND account_id = $2 ORDER BY month ASC;`
	rows, err := r.db.Query(ctx, qry, userID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.AccountStatistics])
}
