package transport

import "time"

type CreateOperationDTO struct {
	AccountID   int64     `json:"account_id"`
	CategoryID  int64     `json:"category_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	ExecutedAt  time.Time `json:"exec_at"`
	Sum         float64   `json:"sum"`
}

type UpdateOperationDTO struct {
	Version     int64      `json:"version"`
	CategoryID  *int64     `json:"category_id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	ExecutedAt  *time.Time `json:"exec_at"`
	Sum         *float64   `json:"sum"`
}

type ResponseOperationDTO struct {
	ID          int64     `json:"id"`
	Version     int64     `json:"version"`
	AccountID   int64     `json:"account_id"`
	CategoryID  int64     `json:"category_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	ExecutedAt  time.Time `json:"exec_at"`
	Sum         float64   `json:"sum"`
}

type ResponseOperationsDTO struct {
	Items  []ResponseOperationDTO `json:"items"`
	Total  int64                  `json:"total,omitempty"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}
