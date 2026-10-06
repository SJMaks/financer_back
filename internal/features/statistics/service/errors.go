package service

import (
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
)

var (
	AccountNotFound = &apperrors.Error{Code: "account_not_found", Msg: "account not found", Status: http.StatusNotFound}
	MissingID       = &apperrors.Error{Code: "missing_id", Msg: "missing account id", Status: http.StatusBadRequest}
	InvalidID       = &apperrors.Error{Code: "invalid_id", Msg: "invalid account id", Status: http.StatusBadRequest}
)
