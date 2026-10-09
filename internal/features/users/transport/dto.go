package transport

type UserRegisterDTO struct {
	Login    string `json:"login" example:"alice"`
	Password string `json:"password" example:"SECret123"`
	Fullname string `json:"full_name" example:"Alice White"`
}

type UserResponseDTO struct {
	Login    string `json:"login" example:"alice"`
	Fullname string `json:"full_name" example:"Alice White"`
	Currency string `json:"currency" example:"dollar"`
}

type LoginRequestDTO struct {
	Login    string `json:"login" example:"alice"`
	Password string `json:"password" example:"SECret123"`
}

type UpdateRequestDTO struct {
	Fullname *string `json:"full_name" example:"Alice Pink"`
	Currency *string `json:"currency" example:"euro"`
}

type LoginResponseDTO struct {
	Token string `json:"jwt"`
}
