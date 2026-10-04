package transport

import "github.com/SJMaks/financer_back/internal/features/operation_types/domain"

func DomainToResponseDTO(c *domain.OperationType) *ResponseOperationTypeDTO {
	return &ResponseOperationTypeDTO{
		ID:    c.ID,
		Title: c.Title,
	}
}

func DomainsToResponseDTO(accs []domain.OperationType) *ResponseOperationTypesDTO {
	dtos := make([]ResponseOperationTypeDTO, len(accs))
	for i, c := range accs {
		dtos[i] = *DomainToResponseDTO(&c)
	}
	return &ResponseOperationTypesDTO{
		OperationTypes: dtos,
	}
}
