package transport

type CreateAccountDTO struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color"`
	Balance     float64 `json:"balance"`
}

type UpdateAccountDTO struct {
	Version     int64    `json:"version"`
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Icon        *string  `json:"icon"`
	Color       *string  `json:"color"`
	Balance     *float64 `json:"balance"`
}

type ResponseAccountDTO struct {
	ID          int64   `json:"id"`
	Version     int64   `json:"version"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Color       string  `json:"color"`
	Balance     float64 `json:"balance"`
}

type ResponseAccountsDTO struct {
	Accounts []ResponseAccountDTO `json:"accounts"`
}
