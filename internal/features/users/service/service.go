package service

type Repository interface {
	Ping() (string, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Ping() (string, error) {
	return s.repo.Ping()
}
