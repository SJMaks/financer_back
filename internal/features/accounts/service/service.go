package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/SJMaks/financer_back/internal/features/accounts/domain"
)

type Repository interface {
	Create(ctx context.Context, a *domain.Account) (*domain.Account, error)
	GetAllOfUser(ctx context.Context, userID int64) ([]domain.Account, error)
	GetByID(ctx context.Context, id, userID int64) (*domain.Account, error)
	UpdateByID(ctx context.Context, a *domain.PatchAccount) (*domain.Account, error)
	DeleteByID(ctx context.Context, id, userID int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	title := strings.TrimSpace(a.Title)
	a.Title = title
	var description, icon, color string
	if a.Description != nil {
		description = strings.TrimSpace(*a.Description)
		a.Description = &description
	}
	if a.Icon != nil {
		icon = strings.TrimSpace(*a.Icon)
		a.Icon = &icon
	}
	if a.Color != nil {
		color = strings.TrimSpace(*a.Color)
		a.Color = &color
	}

	if title == "" {
		return nil, TitleRequired
	}
	if utf8.RuneCountInString(title) > 50 {
		return nil, InvalidTitle
	}

	if utf8.RuneCountInString(description) > 500 {
		return nil, InvalidDescription
	}
	if utf8.RuneCountInString(icon) > 50 {
		return nil, InvalidIcon
	}
	if utf8.RuneCountInString(color) > 7 {
		return nil, InvalidColor
	}

	if a.Balance < 0 {
		return nil, InvalidBalance
	}

	return s.repo.Create(ctx, a)
}

func (s *Service) GetAllOfUser(ctx context.Context, userID int64) ([]domain.Account, error) {
	return s.repo.GetAllOfUser(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id, userID int64) (*domain.Account, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *Service) UpdateByID(ctx context.Context, a *domain.PatchAccount) (*domain.Account, error) {
	if a.Version == nil {
		return nil, MissingVersion
	}
	var title, description, icon, color string
	if a.Title != nil {
		title = strings.TrimSpace(*a.Title)
		a.Title = &title
	}
	if a.Description != nil {
		description = strings.TrimSpace(*a.Description)
		a.Description = &description
	}
	if a.Icon != nil {
		icon = strings.TrimSpace(*a.Icon)
		a.Icon = &icon
	}
	if a.Color != nil {
		color = strings.TrimSpace(*a.Color)
		a.Color = &color
	}

	if utf8.RuneCountInString(title) > 50 {
		return nil, InvalidTitle
	}
	if utf8.RuneCountInString(description) > 500 {
		return nil, InvalidDescription
	}
	if utf8.RuneCountInString(icon) > 50 {
		return nil, InvalidIcon
	}
	if utf8.RuneCountInString(color) > 7 {
		return nil, InvalidColor
	}
	if a.Balance != nil && *a.Balance < 0 {
		return nil, InvalidBalance
	}

	return s.repo.UpdateByID(ctx, a)
}

func (s *Service) DeleteByID(ctx context.Context, id, userID int64) error {
	return s.repo.DeleteByID(ctx, id, userID)
}
