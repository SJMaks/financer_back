package repository

import (
	"context"
	"errors"

	"github.com/SJMaks/financer_back/internal/features/users/domain"
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

var _ service.UserRepository = (*PostgresRepository)(nil)

func (r *PostgresRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	qry := `SELECT id, version, login, pass_hash, full_name, currency FROM users WHERE login = $1;`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, qry, login).Scan(&u.ID, &u.Version, &u.Login, &u.Password, &u.Fullname, &u.Currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, service.UserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	qry := `INSERT INTO users (login, pass_hash, full_name, currency) VALUES ($1, $2, $3, 'dollar')
	RETURNING id, version, login, pass_hash, full_name, currency;`
	err := r.db.QueryRow(ctx, qry, u.Login, u.Password, u.Fullname).Scan(
		&u.ID, &u.Version, &u.Login, &u.Password, &u.Fullname, &u.Currency)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, service.UserAlreadyExists
	}
	if err != nil {
		return nil, err
	}

	categoriesQry := `INSERT INTO operation_categories (user_id, operation_type_id, title, icon, color) 
	VALUES ($1, 1, 'no_category', 'Receipt', '#828282');`
	_, err = r.db.Exec(ctx, categoriesQry, u.ID)
	if err != nil {
		return nil, err
	}
	categoriesQry = `INSERT INTO operation_categories (user_id, operation_type_id, title, icon, color) 
	VALUES ($1, 2, 'no_category', 'Receipt', '#828282');`
	_, err = r.db.Exec(ctx, categoriesQry, u.ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	qry := `SELECT id, version, login, pass_hash, full_name, currency FROM users
	WHERE id = $1`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, qry, id).Scan(&u.ID, &u.Version, &u.Login, &u.Password, &u.Fullname, &u.Currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, service.UserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// TODO: если возникнет оптимистичной блокировки, реализовать UpdateByIDWithVersion
func (r *PostgresRepository) UpdateByID(ctx context.Context, u *domain.PatchUser) (*domain.User, error) {
	qry := `UPDATE users SET full_name = COALESCE($1, full_name),
	currency = COALESCE($2, currency),
	version = version + 1 WHERE id = $3 RETURNING
	id, version, login, pass_hash, full_name, currency;`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, qry, u.Fullname, u.Currency, u.ID).Scan(
		&user.ID, &user.Version, &user.Login, &user.Password, &user.Fullname, &user.Currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, service.UserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id int64) error {
	qry := `DELETE FROM users WHERE id = $1;`
	tag, err := r.db.Exec(ctx, qry, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return service.UserNotFound
	}
	return err
}
