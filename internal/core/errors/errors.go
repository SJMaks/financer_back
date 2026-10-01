package errors

import "errors"

var (
	EnvNotFound = errors.New("env var not found")
)
