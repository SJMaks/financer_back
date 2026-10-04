package service

import "errors"

var (
	UserAlreadyExists   = errors.New("user with this login already exists")
	UserNotFound        = errors.New("user not found")
	EmptyLogin          = errors.New("login cannot be empty")
	WrongPasswordFormat = errors.New("password should contain only latin characters, at least one capital letter, at least one digit and be 5 to 30 characters length")
	EmptyFullname       = errors.New("full name cannot be empty")
	InvalidCredentials  = errors.New("invalid credentials")
)
