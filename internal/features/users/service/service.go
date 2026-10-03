package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/SJMaks/financer_back/internal/core/auth"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetByLogin(ctx context.Context, login string) (*User, error)
	Create(ctx context.Context, u *User) (*User, error)
}

type Service struct {
	repo Repository
	tm   auth.TokenManager
}

func NewService(repo Repository, tm auth.TokenManager) *Service {
	return &Service{
		repo: repo,
		tm:   tm,
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

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.PasswordHash = string(hash)
	return s.repo.Create(ctx, &u)
}

func (s *Service) Login(ctx context.Context, u User, password string) (string, error) {
	existingUser, err := s.repo.GetByLogin(ctx, u.Login)
	if errors.Is(err, UserNotFound) {
		return "", InvalidCredentials
	}
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password)); err != nil {
		return "", InvalidCredentials
	}
	payload := auth.TokenPayload{
		ID:    existingUser.ID,
		Login: existingUser.Login,
	}
	token, err := s.tm.Generate(payload)
	if err != nil {
		return "", err
	}
	return token, nil
}
