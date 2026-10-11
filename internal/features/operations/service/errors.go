package service

import (
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
)

var (
	OperationNotFound     = &apperrors.Error{Code: "operation_not_found", Msg: "operation not found", Status: http.StatusNotFound}
	AccountNotFound       = &apperrors.Error{Code: "account_not_found", Msg: "account not found", Status: http.StatusNotFound}
	AccountIDRequired     = &apperrors.Error{Code: "account_id_required", Msg: "account id is required", Status: http.StatusBadRequest}
	CategoryIDRequired    = &apperrors.Error{Code: "category_id_required", Msg: "category id is required", Status: http.StatusBadRequest}
	TitleRequired         = &apperrors.Error{Code: "title_required", Msg: "title is required", Status: http.StatusBadRequest}
	InvalidTitle          = &apperrors.Error{Code: "invalid_title", Msg: "invalid title", Status: http.StatusBadRequest}
	InvalidDescription    = &apperrors.Error{Code: "invalid_description", Msg: "invalid description", Status: http.StatusBadRequest}
	InvalidExecAt         = &apperrors.Error{Code: "invalid_exec_at", Msg: "invalid executed at date", Status: http.StatusBadRequest}
	MissingAccountID      = &apperrors.Error{Code: "missing_id", Msg: "missing account id", Status: http.StatusBadRequest}
	InvalidAccountID      = &apperrors.Error{Code: "invalid_id", Msg: "invalid account id", Status: http.StatusBadRequest}
	MissingOperationID    = &apperrors.Error{Code: "missing_id", Msg: "missing operation id", Status: http.StatusBadRequest}
	InvalidOperationID    = &apperrors.Error{Code: "invalid_id", Msg: "invalid operation id", Status: http.StatusBadRequest}
	MissingVersion        = &apperrors.Error{Code: "missing_version", Msg: "missing version", Status: http.StatusBadRequest}
	UpdateVersionConflict = &apperrors.Error{Code: "update_conflict", Msg: "update version conflict", Status: http.StatusConflict}
)
