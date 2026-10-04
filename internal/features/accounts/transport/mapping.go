package transport

import "github.com/SJMaks/financer_back/internal/features/accounts/domain"

func CreateAccountDTOToDomain(dto *CreateAccountDTO) *domain.Account {
	return &domain.Account{
		Title:       dto.Title,
		Description: dto.Description,
		Icon:        dto.Icon,
		Color:       dto.Color,
		Balance:     dto.Balance,
	}
}

func UpdateAccountDTOToDomain(dto *UpdateAccountDTO) *domain.Account {
	return &domain.Account{
		Version:     dto.Version,
		Title:       dto.Title,
		Description: dto.Description,
		Icon:        dto.Icon,
		Color:       dto.Color,
		Balance:     dto.Balance,
	}
}

func DomainToResponseDTO(a *domain.Account) *ResponseAccountDTO {
	return &ResponseAccountDTO{
		ID:          a.ID,
		Version:     a.Version,
		Title:       a.Title,
		Description: a.Description,
		Icon:        a.Icon,
		Color:       a.Color,
		Balance:     a.Balance,
	}
}

func DomainsToResponseDTO(accs []domain.Account) *ResponseAccountsDTO {
	dtos := make([]ResponseAccountDTO, len(accs))
	for i, a := range accs {
		dtos[i] = *DomainToResponseDTO(&a)
	}
	return &ResponseAccountsDTO{
		Accounts: dtos,
	}
}
