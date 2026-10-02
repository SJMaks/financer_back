package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetByLogin(ctx context.Context, login string) (*User, error)
	Create(ctx context.Context, u *User) (*User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

var (
	reAllowed = regexp.MustCompile(`^[A-Za-z0-9]{5,30}$`)
	reUpper   = regexp.MustCompile(`[A-Z]`)
	reDigit   = regexp.MustCompile(`[0-9]`)
)

func validatePassword(p string) bool {
	return reAllowed.MatchString(p) &&
		reUpper.MatchString(p) &&
		reDigit.MatchString(p)
}

func (s *Service) Register(ctx context.Context, u User, password string) (*User, error) {
	if len(strings.TrimSpace(u.Login)) < 1 {
		return nil, EmptyLogin
	}
	if passOk := validatePassword(password); !passOk {
		return nil, WrongPasswordFormat
	}
	if len(strings.TrimSpace(u.Fullname)) < 1 {
		return nil, EmptyFullname
	}
	existingUser, err := s.repo.GetByLogin(ctx, u.Login)
	if err != nil && !errors.Is(err, UserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, UserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.PasswordHash = string(hash)
	return s.repo.Create(ctx, &u)
}
