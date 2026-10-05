package domain

type Category struct {
	ID              int64   `json:"id"`
	Version         int64   `json:"version"`
	UserID          int64   `json:"user_id"`
	OperationTypeID int64   `json:"operation_type_id"`
	Title           string  `json:"title"`
	Icon            *string `json:"icon"`
	Color           *string `json:"color"`
}

type PatchCategory struct {
	ID              int64   `json:"id"`
	Version         *int64  `json:"version"`
	UserID          int64   `json:"user_id"`
	OperationTypeID *int64  `json:"operation_type_id"`
	Title           *string `json:"title"`
	Icon            *string `json:"icon"`
	Color           *string `json:"color"`
}
