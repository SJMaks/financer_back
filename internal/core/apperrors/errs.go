package apperrors

import "errors"

var (
	EnvNotFound = errors.New("env var not found")
	EmptyEnv    = errors.New("env is empty")
)
