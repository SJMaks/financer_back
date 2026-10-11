package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SJMaks/financer_back/internal/features/operations/domain"
	"github.com/SJMaks/financer_back/internal/features/operations/service"
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

func (r *PostgresRepository) Create(ctx context.Context, o *domain.Operation) (*domain.Operation, error) {
	qry := `INSERT INTO operations (user_id, account_id, category_id, title, description, exec_at, sum) 
	VALUES ($1, $2, $3, $4, $5, $6, ABS($7)) RETURNING 
	id, version, user_id, account_id, category_id, title, description, exec_at, sum;`
	operation := &domain.Operation{}
	err := r.db.QueryRow(ctx, qry,
		o.UserID,
		o.AccountID,
		o.CategoryID,
		o.Title,
		o.Description,
		o.ExecutedAt,
		o.Sum).Scan(
		&operation.ID,
		&operation.Version,
		&operation.UserID,
		&operation.AccountID,
		&operation.CategoryID,
		&operation.Title,
		&operation.Description,
		&operation.ExecutedAt,
		&operation.Sum,
	)
	if err != nil {
		return nil, err
	}
	return operation, nil
}

func (r *PostgresRepository) GetAllOfAccount(
	ctx context.Context, params service.ListOperationsParams,
) ([]domain.Operation, int64, error) {

	var accountExists bool
	err := r.db.QueryRow(ctx, `
	SELECT EXISTS (SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2);
	`, params.Filter.AccountID, params.Filter.UserID).Scan(&accountExists)
	if err != nil {
		return nil, 0, err
	}
	if !accountExists {
		return nil, 0, service.AccountNotFound
	}

	where := []string{"o.user_id = $1", "o.account_id = $2"}
	args := []any{params.Filter.UserID, params.Filter.AccountID}
	i := 3

	if params.Filter.OperationTypeID != nil {
		where = append(where, fmt.Sprintf(
			`EXISTS (SELECT 1 FROM operation_categories c
			         WHERE c.id = o.category_id AND c.operation_type_id = $%d)`, i))
		args = append(args, *params.Filter.OperationTypeID)
		i++
	}
	if params.Filter.CategoryID != nil {
		where = append(where, fmt.Sprintf("o.category_id = $%d", i))
		args = append(args, *params.Filter.CategoryID)
		i++
	}
	if params.Filter.MinExecAt != nil {
		where = append(where, fmt.Sprintf("o.exec_at >= $%d", i))
		args = append(args, *params.Filter.MinExecAt)
		i++
	}
	if params.Filter.MaxExecAt != nil {
		where = append(where, fmt.Sprintf("o.exec_at <= $%d", i))
		args = append(args, *params.Filter.MaxExecAt)
		i++
	}

	whereQry := "WHERE " + strings.Join(where, " AND ")

	col, dir := r.normalizeSort(params.SortBy, params.SortOrd)
	orderQry := fmt.Sprintf("ORDER BY o.%s %s, o.id ASC", col, dir)

	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	// count — с теми же args, что и фильтр (без limit/offset)
	countQry := `SELECT COUNT(*) FROM operations o ` + whereQry
	var total int64
	if err := r.db.QueryRow(ctx, countQry, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// data — отдельная копия args + limit/offset
	dataArgs := append(append([]any{}, args...), limit, offset)
	dataQry := fmt.Sprintf(`
		SELECT o.id, o.version, o.user_id, o.account_id, o.category_id,
		       o.title, o.description, o.exec_at, o.sum
		FROM operations o
		%s %s
		LIMIT $%d OFFSET $%d`,
		whereQry, orderQry, i, i+1)

	rows, err := r.db.Query(ctx, dataQry, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Operation])
	if err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id, userID, accountID int64) (*domain.Operation, error) {
	qry := `SELECT id, version, user_id, account_id, category_id, title, description, exec_at, sum
	FROM operations
	WHERE id = $1 AND user_id = $2 AND account_id = $3;`
	operation := &domain.Operation{}
	err := r.db.QueryRow(ctx, qry, id, userID, accountID).Scan(
		&operation.ID,
		&operation.Version,
		&operation.UserID,
		&operation.AccountID,
		&operation.CategoryID,
		&operation.Title,
		&operation.Description,
		&operation.ExecutedAt,
		&operation.Sum,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, service.OperationNotFound
	}
	if err != nil {
		return nil, err
	}
	return operation, nil
}

func (r *PostgresRepository) UpdateByID(ctx context.Context, c *domain.PatchOperation) (*domain.Operation, error) {
	qry := `UPDATE operations SET 
	version = version + 1,
	category_id = COALESCE($1, category_id),
	title = COALESCE($2, title),
	description = COALESCE($3, description),
	exec_at = COALESCE($4, exec_at),
	sum = COALESCE(ABS($5), sum)
	WHERE id = $6 AND user_id = $7 AND version = $8
	RETURNING id, version, user_id, account_id, category_id, title, description, exec_at, sum;`
	operation := &domain.Operation{}
	err := r.db.QueryRow(ctx, qry,
		c.CategoryID,
		c.Title,
		c.Description,
		c.ExecutedAt,
		c.Sum,
		c.ID,
		c.UserID,
		c.Version).Scan(
		&operation.ID,
		&operation.Version,
		&operation.UserID,
		&operation.AccountID,
		&operation.CategoryID,
		&operation.Title,
		&operation.Description,
		&operation.ExecutedAt,
		&operation.Sum,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		var exists bool
		_ = r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM operations 
		WHERE id = $1 AND user_id = $2);`, c.ID, c.UserID).Scan(&exists)

		if !exists {
			return nil, service.OperationNotFound
		}
		return nil, service.UpdateVersionConflict
	}
	if err != nil {
		return nil, err
	}
	return operation, nil
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id, userID, accountID int64) error {
	qry := `DELETE FROM operations WHERE id = $1 AND user_id = $2 AND account_id = $3;`
	tag, err := r.db.Exec(ctx, qry, id, userID, accountID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return service.OperationNotFound
	}
	return err
}

var allowedSort = map[string]string{
	"exec_at": "exec_at",
	"sum":     "sum",
}

var allowedSortOrd = map[string]string{
	"asc":  "ASC",
	"desc": "DESC",
}

func (r *PostgresRepository) normalizeSort(sortBy, sortOrd string) (string, string) {
	col, ok := allowedSort[sortBy]
	if !ok {
		col = "id"
	}
	ord, ok := allowedSortOrd[sortOrd]
	if !ok {
		ord = "ASC"
	}
	return col, ord
}
