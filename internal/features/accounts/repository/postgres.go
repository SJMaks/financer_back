package repository

import (
	"context"
	"errors"

	"github.com/SJMaks/financer_back/internal/features/accounts/domain"
	"github.com/SJMaks/financer_back/internal/features/accounts/service"
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

func (r *PostgresRepository) Create(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	qry := `INSERT INTO accounts (user_id, title, description, icon, color, balance) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING 
	id, version, user_id, title, description, icon, color, balance;`
	account := &domain.Account{}
	err := r.db.QueryRow(ctx, qry, a.UserID, a.Title, a.Description, a.Icon, a.Color, a.Balance).Scan(
		&account.ID,
		&account.Version,
		&account.UserID,
		&account.Title,
		&account.Description,
		&account.Icon,
		&account.Color,
		&account.Balance,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *PostgresRepository) GetAllOfUser(ctx context.Context, userID int64) ([]domain.Account, error) {
	qry := `SELECT id, version, user_id, title, description, icon, color, balance FROM accounts
	WHERE user_id = $1;`
	rows, err := r.db.Query(ctx, qry, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Account])
}

func (r *PostgresRepository) GetByID(ctx context.Context, id, userID int64) (*domain.Account, error) {
	qry := `SELECT id, version, user_id, title, description, icon, color, balance FROM accounts
	WHERE id = $1 AND user_id = $2;`
	account := &domain.Account{}
	err := r.db.QueryRow(ctx, qry, id, userID).Scan(
		&account.ID,
		&account.Version,
		&account.UserID,
		&account.Title,
		&account.Description,
		&account.Icon,
		&account.Color,
		&account.Balance,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, service.AccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *PostgresRepository) UpdateByID(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	qry := `UPDATE accounts SET 
	version = version + 1,
	title = $1, 
	description = $2, 
	icon = $3, 
	color = $4, 
	balance = $5
	WHERE id = $6 AND user_id = $7 AND version = $8
	RETURNING id, version, user_id, title, description, icon, color, balance;`
	account := &domain.Account{}
	err := r.db.QueryRow(ctx, qry,
		a.Title,
		a.Description,
		a.Icon,
		a.Color,
		a.Balance,
		a.ID,
		a.UserID,
		a.Version).Scan(
		&account.ID,
		&account.Version,
		&account.UserID,
		&account.Title,
		&account.Description,
		&account.Icon,
		&account.Color,
		&account.Balance,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		var exists bool
		_ = r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM accounts 
		WHERE id = $1 AND user_id = $2);`, a.ID, a.UserID).Scan(&exists)

		if !exists {
			return nil, service.AccountNotFound
		}
		return nil, service.UpdateVersionConflict
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id, userID int64) error {
	qry := `DELETE FROM accounts WHERE id = $1 AND user_id = $2;`
	tag, err := r.db.Exec(ctx, qry, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return service.AccountNotFound
	}
	return err
}
