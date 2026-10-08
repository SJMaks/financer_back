package transport

import (
	"github.com/SJMaks/financer_back/internal/features/operations/domain"
)

func CreateOperationDTOToDomain(dto *CreateOperationDTO) *domain.Operation {
	return &domain.Operation{
		CategoryID:  dto.CategoryID,
		Title:       dto.Title,
		Description: dto.Description,
		ExecutedAt:  dto.ExecutedAt,
		Sum:         dto.Sum,
	}
}

func UpdateOperationDTOToDomain(dto *UpdateOperationDTO) *domain.PatchOperation {
	return &domain.PatchOperation{
		Version:     dto.Version,
		CategoryID:  dto.CategoryID,
		Title:       dto.Title,
		Description: dto.Description,
		ExecutedAt:  dto.ExecutedAt,
		Sum:         dto.Sum,
	}
}

func DomainToResponseDTO(o *domain.Operation) *ResponseOperationDTO {
	return &ResponseOperationDTO{
		ID:          o.ID,
		Version:     o.Version,
		AccountID:   o.AccountID,
		CategoryID:  o.CategoryID,
		Title:       o.Title,
		Description: o.Description,
		ExecutedAt:  o.ExecutedAt,
		Sum:         o.Sum,
	}
}

func DomainsToResponseDTO(ops []domain.Operation, total int64, limit, offset int) *ResponseOperationsDTO {
	dtos := make([]ResponseOperationDTO, len(ops))
	for i, o := range ops {
		dtos[i] = *DomainToResponseDTO(&o)
	}
	return &ResponseOperationsDTO{
		Items:  dtos,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}
