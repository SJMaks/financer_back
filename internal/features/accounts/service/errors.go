package service

import "errors"

var (
	AccountNotFound       = errors.New("account not found")
	EmptyTitle            = errors.New("title cannot be empty")
	NegativeBalance       = errors.New("balance cannot be negative")
	MissingID             = errors.New("missing account id")
	InvalidID             = errors.New("invalid account id")
	UpdateVersionConflict = errors.New("update version conflict")
)
