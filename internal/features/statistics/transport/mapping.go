package transport

import "github.com/SJMaks/financer_back/internal/features/statistics/domain"

func DomainToByUserDTO(s *domain.UserStatistics) *UserStatisticsDTO {
	return &UserStatisticsDTO{
		UserID:  s.UserID,
		Month:   s.Month,
		Income:  s.Income,
		Expense: s.Expense,
		Profit:  s.Profit,
	}
}

func DomainToByAccountDTO(s *domain.AccountStatistics) *AccountStatisticsDTO {
	return &AccountStatisticsDTO{
		AccountID: s.AccountID,
		UserID:    s.UserID,
		Month:     s.Month,
		Income:    s.Income,
		Expense:   s.Expense,
		Profit:    s.Profit,
	}
}

func ByUserDomainsToSliceDTO(stcs []domain.UserStatistics) *UserStatisticsSliceDTO {
	dtos := make([]UserStatisticsDTO, len(stcs))
	for i, s := range stcs {
		dtos[i] = *DomainToByUserDTO(&s)
	}
	return &UserStatisticsSliceDTO{
		Months: dtos,
	}
}

func ByAccountDomainsToSliceDTO(stcs []domain.AccountStatistics) *AccountStatisticsSliceDTO {
	dtos := make([]AccountStatisticsDTO, len(stcs))
	for i, s := range stcs {
		dtos[i] = *DomainToByAccountDTO(&s)
	}
	return &AccountStatisticsSliceDTO{
		Months: dtos,
	}
}
