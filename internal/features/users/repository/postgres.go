package repository

import (
	"context"
	"errors"

	"github.com/SJMaks/financer_back/internal/features/users/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *PostgresRepository) GetByLogin(ctx context.Context, login string) (*service.User, error) {
	qry := `SELECT id, version, login, pass_hash, full_name FROM users WHERE login = $1;`
	u := &service.User{}
	err := r.db.QueryRow(ctx, qry, login).Scan(&u.ID, &u.Version, &u.Login, &u.PasswordHash, &u.Fullname)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, service.UserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepository) Create(ctx context.Context, u *service.User) (*service.User, error) {
	qry := `INSERT INTO users (login, pass_hash, full_name) VALUES ($1, $2, $3)
	RETURNING id, version, login, pass_hash, full_name;`
	err := r.db.QueryRow(ctx, qry, u.Login, u.PasswordHash, u.Fullname).Scan(
		&u.ID, &u.Version, &u.Login, &u.PasswordHash, &u.Fullname)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, service.UserAlreadyExists
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
