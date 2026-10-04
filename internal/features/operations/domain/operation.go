package domain

import "time"

type Operation struct {
	ID          int64     `db:"id"`
	Version     int64     `db:"version"`
	UserID      int64     `db:"user_id"`
	AccountID   int64     `db:"account_id"`
	CategoryID  int64     `db:"category_id"`
	Title       string    `db:"title"`
	Description *string   `db:"description"`
	ExecutedAt  time.Time `db:"exec_at"`
	Sum         float64   `db:"sum"`
}

type PatchOperation struct {
	ID          int64      `json:"id"`
	Version     int64      `json:"version"`
	UserID      int64      `json:"user_id"`
	AccountID   int64      `json:"account_id"`
	CategoryID  *int64     `json:"category_id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	ExecutedAt  *time.Time `json:"exec_at"`
	Sum         *float64   `json:"sum"`
}
