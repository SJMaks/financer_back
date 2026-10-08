package transport

type CreateAccountDTO struct {
	Title       string  `json:"title" example:"Visa"`
	Description *string `json:"description" example:"An account for groceries"`
	Icon        *string `json:"icon" example:"dollar"`
	Color       *string `json:"color" example:"#197b15"`
	Balance     float64 `json:"balance" example:"15.55"`
}

type UpdateAccountDTO struct {
	Version     *int64   `json:"version" example:"0"`
	Title       *string  `json:"title" example:"Visa 1234"`
	Description *string  `json:"description" example:"An account for clothes"`
	Icon        *string  `json:"icon" example:"t-shirt"`
	Color       *string  `json:"color" example:"#15397b"`
	Balance     *float64 `json:"balance" example:"50.35"`
}

type ResponseAccountDTO struct {
	ID          int64   `json:"id" example:"1"`
	Version     int64   `json:"version" example:"0"`
	Title       string  `json:"title" example:"Visa"`
	Description string  `json:"description" example:"An account for groceries"`
	Icon        string  `json:"icon" example:"dollar"`
	Color       string  `json:"color" example:"#197b15"`
	Balance     float64 `json:"balance" example:"15.55"`
}

type ResponseAccountsDTO struct {
	Accounts []ResponseAccountDTO `json:"accounts"`
}
