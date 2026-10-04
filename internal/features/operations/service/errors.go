package service

import "errors"

var (
	OperationNotFound     = errors.New("operation not found")
	EmptyTitle            = errors.New("title cannot be empty")
	MissingAccountID      = errors.New("missing account id")
	MissingCategoryID     = errors.New("missing category id")
	NegativeSum           = errors.New("sum cannot be negative")
	MissingID             = errors.New("missing operation id")
	InvalidID             = errors.New("invalid operation id")
	UpdateVersionConflict = errors.New("update version conflict")
)
