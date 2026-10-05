package service

import (
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
)

var (
	CategoryNotFound        = &apperrors.Error{Code: "category_not_found", Msg: "category not found", Status: http.StatusNotFound}
	OperationTypeIDRequired = &apperrors.Error{Code: "operation_type_id_required", Msg: "operation type id is required", Status: http.StatusBadRequest}
	TitleRequired           = &apperrors.Error{Code: "title_required", Msg: "title is required", Status: http.StatusBadRequest}
	InvalidTitle            = &apperrors.Error{Code: "invalid_title", Msg: "invalid title", Status: http.StatusBadRequest}
	InvalidIcon             = &apperrors.Error{Code: "invalid_icon", Msg: "invalid icon", Status: http.StatusBadRequest}
	InvalidColor            = &apperrors.Error{Code: "invalid_color", Msg: "invalid color", Status: http.StatusBadRequest}
	MissingID               = &apperrors.Error{Code: "missing_id", Msg: "missing account id", Status: http.StatusBadRequest}
	InvalidID               = &apperrors.Error{Code: "invalid_id", Msg: "invalid account id", Status: http.StatusBadRequest}
	MissingVersion          = &apperrors.Error{Code: "missing_version", Msg: "missing version", Status: http.StatusBadRequest}
	UpdateVersionConflict   = &apperrors.Error{Code: "update_conflict", Msg: "update version conflict", Status: http.StatusConflict}
)
