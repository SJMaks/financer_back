package domain

type Account struct {
	ID          int64   `json:"id"`
	Version     int64   `json:"version"`
	UserID      int64   `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color"`
	Balance     float64 `json:"balance"`
}

type PatchAccount struct {
	ID          int64    `json:"id"`
	Version     *int64   `json:"version"`
	UserID      int64    `json:"user_id"`
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Icon        *string  `json:"icon"`
	Color       *string  `json:"color"`
	Balance     *float64 `json:"balance"`
}
