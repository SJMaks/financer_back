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

func UpdateAccountDTOToDomain(dto *UpdateAccountDTO) *domain.PatchAccount {
	return &domain.PatchAccount{
		Version:     dto.Version,
		Title:       dto.Title,
		Description: dto.Description,
		Icon:        dto.Icon,
		Color:       dto.Color,
		Balance:     dto.Balance,
	}
}

func DomainToResponseDTO(a *domain.Account) *ResponseAccountDTO {
	var description, icon, color string
	if a.Description != nil {
		description = *a.Description
	}
	if a.Icon != nil {
		icon = *a.Icon
	}
	if a.Color != nil {
		color = *a.Color
	}
	return &ResponseAccountDTO{
		ID:          a.ID,
		Version:     a.Version,
		Title:       a.Title,
		Description: description,
		Icon:        icon,
		Color:       color,
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
