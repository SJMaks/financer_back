package service

import (
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
)

var (
	AccountNotFound       = &apperrors.Error{Code: "account_not_found", Msg: "account not found", Status: http.StatusNotFound}
	TitleRequired         = &apperrors.Error{Code: "title_required", Msg: "title is required", Status: http.StatusBadRequest}
	InvalidTitle          = &apperrors.Error{Code: "invalid_title", Msg: "invalid title", Status: http.StatusBadRequest}
	InvalidDescription    = &apperrors.Error{Code: "invalid_description", Msg: "invalid description", Status: http.StatusBadRequest}
	InvalidIcon           = &apperrors.Error{Code: "invalid_icon", Msg: "invalid icon", Status: http.StatusBadRequest}
	InvalidColor          = &apperrors.Error{Code: "invalid_color", Msg: "invalid color", Status: http.StatusBadRequest}
	InvalidBalance        = &apperrors.Error{Code: "invalid_balance", Msg: "invalid balance", Status: http.StatusBadRequest}
	MissingID             = &apperrors.Error{Code: "missing_id", Msg: "missing account id", Status: http.StatusBadRequest}
	InvalidID             = &apperrors.Error{Code: "invalid_id", Msg: "invalid account id", Status: http.StatusBadRequest}
	MissingVersion        = &apperrors.Error{Code: "missing_version", Msg: "missing version", Status: http.StatusBadRequest}
	UpdateVersionConflict = &apperrors.Error{Code: "update_conflict", Msg: "update version conflict", Status: http.StatusConflict}
)
