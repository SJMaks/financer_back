package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/SJMaks/financer_back/internal/features/users/domain"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	UpdateByID(ctx context.Context, u *domain.PatchUser) (*domain.User, error)
	DeleteByID(ctx context.Context, id int64) error
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

func (s *Service) Register(ctx context.Context, u *domain.User, password string) (*domain.User, error) {
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

	u.Password = string(hash)
	return s.repo.Create(ctx, u)
}

func (s *Service) Login(ctx context.Context, u *domain.User) (string, error) {
	existingUser, err := s.repo.GetByLogin(ctx, u.Login)
	if errors.Is(err, UserNotFound) {
		return "", InvalidCredentials
	}
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(u.Password)); err != nil {
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

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, u *domain.PatchUser) (*domain.User, error) {
	if u.Fullname != nil && len(strings.TrimSpace(*u.Fullname)) < 1 {
		return nil, EmptyFullname
	}
	return s.repo.UpdateByID(ctx, u)
}

func (s *Service) DeleteByID(ctx context.Context, id int64) error {
	return s.repo.DeleteByID(ctx, id)
}
