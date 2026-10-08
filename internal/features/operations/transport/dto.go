package transport

import "time"

type CreateOperationDTO struct {
	CategoryID  int64     `json:"category_id" example:"1"`
	Title       string    `json:"title" example:"Apples"`
	Description *string   `json:"description" example:"Bought some apples"`
	ExecutedAt  time.Time `json:"exec_at" example:"2026-08-09T00:00:00Z"`
	Sum         float64   `json:"sum" example:"21.55"`
}

type UpdateOperationDTO struct {
	Version     *int64     `json:"version" example:"0"`
	CategoryID  *int64     `json:"category_id" example:"2"`
	Title       *string    `json:"title" example:"Bananas"`
	Description *string    `json:"description" example:"Bought some bananas"`
	ExecutedAt  *time.Time `json:"exec_at" example:"2026-08-10T00:00:00Z"`
	Sum         *float64   `json:"sum" example:"14.99"`
}

type ResponseOperationDTO struct {
	ID          int64     `json:"id" example:"1"`
	Version     int64     `json:"version" example:"0"`
	AccountID   int64     `json:"account_id" example:"1"`
	CategoryID  int64     `json:"category_id" example:"1"`
	Title       string    `json:"title" example:"Apples"`
	Description *string   `json:"description" example:"Bought some apples"`
	ExecutedAt  time.Time `json:"exec_at" example:"2026-08-09T00:00:00Z"`
	Sum         float64   `json:"sum" example:"21.55"`
}

type ResponseOperationsDTO struct {
	Items  []ResponseOperationDTO `json:"items"`
	Total  int64                  `json:"total,omitempty"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}
