package transport

import "github.com/SJMaks/financer_back/internal/features/users/domain"

func RegisterDTOToDomain(dto UserRegisterDTO) *domain.User {
	return &domain.User{
		Login:    dto.Login,
		Password: dto.Password,
		Fullname: dto.Fullname,
	}
}

func DomainToResponse(u domain.User) *UserResponseDTO {
	return &UserResponseDTO{
		Login:    u.Login,
		Fullname: u.Fullname,
	}
}

func LoginDTOToDomain(dto LoginRequestDTO) *domain.User {
	return &domain.User{
		Login:    dto.Login,
		Password: dto.Password,
	}
}

func UpdateDTOToDomain(dto UpdateRequestDTO) *domain.PatchUser {
	return &domain.PatchUser{
		Fullname: dto.Fullname,
	}
}
