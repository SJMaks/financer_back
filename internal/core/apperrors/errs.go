package apperrors

import "net/http"

type Error struct {
	Code   string
	Msg    string
	Status int
	cause  error
}

func (e *Error) Error() string { return e.Msg }
func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.Code == e.Code
}

var (
	EnvNotFound  = &Error{Code: "env_not_found", Msg: "env var not found"}
	EmptyEnv     = &Error{Code: "env_empty", Msg: "env var is empty"}
	BadJSON      = &Error{Code: "bad_json", Status: http.StatusBadRequest}
	Unauthorized = &Error{Code: "unauthorized", Msg: "unauthorized", Status: http.StatusUnauthorized}
	IvalidQuery  = &Error{Code: "invalid_query", Status: http.StatusBadRequest}
)
