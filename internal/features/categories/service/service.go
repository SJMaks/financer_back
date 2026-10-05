package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/SJMaks/financer_back/internal/features/categories/domain"
)

type Repository interface {
	Create(ctx context.Context, a *domain.Category) (*domain.Category, error)
	GetAllOfUser(ctx context.Context, userID int64) ([]domain.Category, error)
	GetByID(ctx context.Context, id, userID int64) (*domain.Category, error)
	UpdateByID(ctx context.Context, a *domain.PatchCategory) (*domain.Category, error)
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

func (s *Service) Create(ctx context.Context, c *domain.Category) (*domain.Category, error) {
	if c.OperationTypeID == 0 {
		return nil, OperationTypeIDRequired
	}

	title := strings.TrimSpace(c.Title)
	c.Title = title
	var icon, color string
	if c.Icon != nil {
		icon = strings.TrimSpace(*c.Icon)
		c.Icon = &icon
	}
	if c.Color != nil {
		color = strings.TrimSpace(*c.Color)
		c.Color = &color
	}

	if title == "" {
		return nil, TitleRequired
	}
	if utf8.RuneCountInString(title) > 50 {
		return nil, InvalidTitle
	}

	if utf8.RuneCountInString(icon) > 50 {
		return nil, InvalidIcon
	}
	if utf8.RuneCountInString(color) > 7 {
		return nil, InvalidColor
	}

	return s.repo.Create(ctx, c)
}

func (s *Service) GetAllOfUser(ctx context.Context, userID int64) ([]domain.Category, error) {
	return s.repo.GetAllOfUser(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id, userID int64) (*domain.Category, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *Service) UpdateByID(ctx context.Context, c *domain.PatchCategory) (*domain.Category, error) {
	if c.Version == nil {
		return nil, MissingVersion
	}
	var title, icon, color string
	if c.Title != nil {
		title = strings.TrimSpace(*c.Title)
		c.Title = &title
	}
	if c.Icon != nil {
		icon = strings.TrimSpace(*c.Icon)
		c.Icon = &icon
	}
	if c.Color != nil {
		color = strings.TrimSpace(*c.Color)
		c.Color = &color
	}

	if utf8.RuneCountInString(title) > 50 {
		return nil, InvalidTitle
	}
	if utf8.RuneCountInString(icon) > 50 {
		return nil, InvalidIcon
	}
	if utf8.RuneCountInString(color) > 7 {
		return nil, InvalidColor
	}

	return s.repo.UpdateByID(ctx, c)
}

func (s *Service) DeleteByID(ctx context.Context, id, userID int64) error {
	return s.repo.DeleteByID(ctx, id, userID)
}
