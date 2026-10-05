package service

import (
	"net/http"

	"github.com/SJMaks/financer_back/internal/core/apperrors"
)

var (
	UserAlreadyExists  = &apperrors.Error{Code: "user_exists", Msg: "user with this login already exists", Status: http.StatusConflict}
	UserNotFound       = &apperrors.Error{Code: "user_not_found", Msg: "user not found", Status: http.StatusNotFound}
	LoginRequired      = &apperrors.Error{Code: "login_required", Msg: "login is required", Status: http.StatusBadRequest}
	PasswordRequired   = &apperrors.Error{Code: "password_required", Msg: "password is required", Status: http.StatusBadRequest}
	FullnameRequired   = &apperrors.Error{Code: "fullname_required", Msg: "full name is required", Status: http.StatusBadRequest}
	InvalidLogin       = &apperrors.Error{Code: "invalid_login", Msg: "invalid login", Status: http.StatusBadRequest}
	InvalidPassword    = &apperrors.Error{Code: "invalid_password", Msg: "invalid password", Status: http.StatusBadRequest}
	InvalidFullname    = &apperrors.Error{Code: "invalid_fullname", Msg: "invalid fullname", Status: http.StatusBadRequest}
	InvalidCredentials = &apperrors.Error{Code: "invalid_credentials", Msg: "invalid credentials", Status: http.StatusUnauthorized}
)
