package transport

type UserRegisterDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Fullname string `json:"full_name"`
}

type UserResponseDTO struct {
	Login    string `json:"login"`
	Fullname string `json:"full_name"`
}

type LoginRequestDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateRequestDTO struct {
	Fullname *string `json:"full_name"`
}

type LoginResponseDTO struct {
	Token string `json:"jwt"`
}
