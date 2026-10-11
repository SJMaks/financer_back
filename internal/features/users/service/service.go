package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/SJMaks/financer_back/internal/core/auth"
	"github.com/SJMaks/financer_back/internal/features/users/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	UpdateByID(ctx context.Context, u *domain.PatchUser) (*domain.User, error)
	DeleteByID(ctx context.Context, id int64) error
}

type Service struct {
	repo UserRepository
	tm   auth.TokenManager
}

func NewService(repo UserRepository, tm auth.TokenManager) *Service {
	return &Service{
		repo: repo,
		tm:   tm,
	}
}

var (
	reLoginAllowed    = regexp.MustCompile(`^[A-Za-z0-9]{5,30}$`)
	reNameAllowed     = regexp.MustCompile(`^[\p{L}]+(?: [\p{L}]+)*$`)
	rePassAllowed     = regexp.MustCompile(`^[A-Za-z0-9]{5,50}$`)
	reCurrencyAllowed = regexp.MustCompile(`^[A-Za-z]{1,50}$`)
	rePassUpper       = regexp.MustCompile(`[A-Z]`)
	rePassDigit       = regexp.MustCompile(`[0-9]`)
)

func validateLogin(l string) bool {
	return reLoginAllowed.MatchString(l)
}

func validatePassword(p string) bool {
	return rePassAllowed.MatchString(p) &&
		rePassUpper.MatchString(p) &&
		rePassDigit.MatchString(p)
}

func validateFullname(f string) bool {
	return reNameAllowed.MatchString(f)
}

func validateCurrency(c string) bool {
	return reCurrencyAllowed.MatchString(c)
}

func (s *Service) Register(ctx context.Context, u *domain.User, password string) (*domain.User, error) {
	u.Login = strings.TrimSpace(u.Login)
	login := u.Login
	fullname := strings.TrimSpace(u.Fullname)

	switch {
	case login == "":
		return nil, LoginRequired
	case password == "":
		return nil, PasswordRequired
	case fullname == "":
		return nil, FullnameRequired
	}

	if loginOk := validateLogin(login); !loginOk {
		return nil, InvalidLogin
	}
	if utf8.RuneCountInString(login) > 30 {
		return nil, InvalidLogin
	}

	if passOk := validatePassword(password); !passOk {
		return nil, InvalidPassword
	}
	if utf8.RuneCountInString(password) > 50 {
		return nil, InvalidPassword
	}

	if nameOk := validateFullname(fullname); !nameOk {
		return nil, InvalidFullname
	}
	if utf8.RuneCountInString(fullname) > 100 {
		return nil, InvalidFullname
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.Password = string(hash)
	return s.repo.Create(ctx, u)
}

func (s *Service) Login(ctx context.Context, u *domain.User) (string, error) {
	switch {
	case strings.TrimSpace(u.Login) == "":
		return "", LoginRequired
	case strings.TrimSpace(u.Password) == "":
		return "", PasswordRequired
	}

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
	if u.Fullname != nil {
		fullname := strings.TrimSpace(*u.Fullname)
		if fullname == "" {
			return nil, FullnameRequired
		}
		if nameOk := validateFullname(fullname); !nameOk {
			return nil, InvalidFullname
		}
		if utf8.RuneCountInString(fullname) > 100 {
			return nil, InvalidFullname
		}
	}
	if u.Currency != nil {
		currency := strings.TrimSpace(*u.Currency)
		if currency == "" {
			return nil, CurrencyRequired
		}
		if currencyOk := validateCurrency(currency); !currencyOk {
			return nil, InvalidCurrency
		}
		if utf8.RuneCountInString(currency) > 50 {
			return nil, InvalidCurrency
		}
	}

	return s.repo.UpdateByID(ctx, u)
}

func (s *Service) DeleteByID(ctx context.Context, id int64) error {
	return s.repo.DeleteByID(ctx, id)
}
