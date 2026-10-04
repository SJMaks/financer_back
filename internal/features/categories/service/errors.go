package service

import "errors"

var (
	CategoryNotFound       = errors.New("category not found")
	EmptyTitle             = errors.New("title cannot be empty")
	MissingOperationTypeID = errors.New("missing operation type id")
	MissingID              = errors.New("missing category id")
	InvalidID              = errors.New("invalid category id")
	UpdateVersionConflict  = errors.New("update version conflict")
)
