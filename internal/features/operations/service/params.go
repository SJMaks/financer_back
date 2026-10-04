package service

import "time"

type ListOperationsFilter struct {
	UserID          int64
	AccountID       int64
	OperationTypeID *int64
	CategoryID      *int64
	MinExecAt       *time.Time
	MaxExecAt       *time.Time
}

type ListOperationsParams struct {
	Filter  ListOperationsFilter
	Limit   int
	Offset  int
	SortBy  string
	SortOrd string
}
