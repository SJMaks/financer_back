package repository

import (
	"context"
	"errors"

	"github.com/SJMaks/financer_back/internal/features/categories/domain"
	"github.com/SJMaks/financer_back/internal/features/categories/service"
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

func (r *PostgresRepository) Create(ctx context.Context, c *domain.Category) (*domain.Category, error) {
	qry := `INSERT INTO operation_categories (user_id, operation_type_id, title, icon, color) 
	VALUES ($1, $2, $3, $4, $5) RETURNING 
	id, version, user_id, operation_type_id, title, icon, color;`
	category := &domain.Category{}
	err := r.db.QueryRow(ctx, qry, c.UserID, c.OperationTypeID, c.Title, c.Icon, c.Color).Scan(
		&category.ID,
		&category.Version,
		&category.UserID,
		&category.OperationTypeID,
		&category.Title,
		&category.Icon,
		&category.Color,
	)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *PostgresRepository) GetAllOfUser(ctx context.Context, userID int64) ([]domain.Category, error) {
	qry := `SELECT id, version, user_id, operation_type_id, title, icon, color FROM operation_categories
	WHERE user_id = $1 AND title <> 'no_category' ORDER BY id;`
	rows, err := r.db.Query(ctx, qry, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Category])
}

func (r *PostgresRepository) GetByID(ctx context.Context, id, userID int64) (*domain.Category, error) {
	qry := `SELECT id, version, user_id, operation_type_id, title, icon, color FROM operation_categories
	WHERE id = $1 AND user_id = $2;`
	category := &domain.Category{}
	err := r.db.QueryRow(ctx, qry, id, userID).Scan(
		&category.ID,
		&category.Version,
		&category.UserID,
		&category.OperationTypeID,
		&category.Title,
		&category.Icon,
		&category.Color,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, service.CategoryNotFound
	}
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *PostgresRepository) UpdateByID(ctx context.Context, c *domain.PatchCategory) (*domain.Category, error) {
	qry := `UPDATE operation_categories SET 
	version = version + 1,
	operation_type_id = COALESCE($1, operation_type_id),
	title = COALESCE($2, title),
	icon = COALESCE($3, icon),
	color = COALESCE($4, color)
	WHERE id = $5 AND user_id = $6 AND version = $7 AND title <> 'no_category'
	RETURNING id, version, user_id, operation_type_id, title, icon, color;`
	category := &domain.Category{}
	err := r.db.QueryRow(ctx, qry,
		c.OperationTypeID,
		c.Title,
		c.Icon,
		c.Color,
		c.ID,
		c.UserID,
		c.Version).Scan(
		&category.ID,
		&category.Version,
		&category.UserID,
		&category.OperationTypeID,
		&category.Title,
		&category.Icon,
		&category.Color,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		var exists bool
		_ = r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM operation_categories 
		WHERE id = $1 AND user_id = $2);`, c.ID, c.UserID).Scan(&exists)

		if !exists {
			return nil, service.CategoryNotFound
		}
		return nil, service.UpdateVersionConflict
	}
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id, userID int64) error {
	qry := `DELETE FROM operation_categories WHERE id = $1 AND user_id = $2 AND title <> 'no_category';`
	tag, err := r.db.Exec(ctx, qry, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return service.CategoryNotFound
	}
	return err
}
