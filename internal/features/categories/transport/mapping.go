package transport

import "github.com/SJMaks/financer_back/internal/features/categories/domain"

func CreateCategoryDTOToDomain(dto *CreateCategoryDTO) *domain.Category {
	return &domain.Category{
		OperationTypeID: dto.OperationTypeID,
		Title:           dto.Title,
		Icon:            dto.Icon,
		Color:           dto.Color,
	}
}

func UpdateCategoryDTOToDomain(dto *UpdateCategoryDTO) *domain.PatchCategory {
	return &domain.PatchCategory{
		Version:         dto.Version,
		OperationTypeID: dto.OperationTypeID,
		Title:           dto.Title,
		Icon:            dto.Icon,
		Color:           dto.Color,
	}
}

func DomainToResponseDTO(c *domain.Category) *ResponseCategoryDTO {
	return &ResponseCategoryDTO{
		ID:              c.ID,
		Version:         c.Version,
		OperationTypeID: c.OperationTypeID,
		Title:           c.Title,
		Icon:            *c.Icon,
		Color:           *c.Color,
	}
}

func DomainsToResponseDTO(ctgs []domain.Category) *ResponseCategoriesDTO {
	dtos := make([]ResponseCategoryDTO, len(ctgs))
	for i, c := range ctgs {
		dtos[i] = *DomainToResponseDTO(&c)
	}
	return &ResponseCategoriesDTO{
		Categories: dtos,
	}
}
