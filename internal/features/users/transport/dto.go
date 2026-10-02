package transport

import "github.com/SJMaks/financer_back/internal/features/users/service"

type UserRegisterDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Fullname string `json:"full_name"`
}

func (u UserRegisterDTO) UserRegisterToDomain() service.User {
	return service.User{
		Login:    u.Login,
		Fullname: u.Fullname,
	}
}

type UserResponseDTO struct {
	Login    string `json:"login"`
	Fullname string `json:"full_name"`
}

func DomainToResponse(u service.User) UserResponseDTO {
	return UserResponseDTO{
		Login:    u.Login,
		Fullname: u.Fullname,
	}
}
