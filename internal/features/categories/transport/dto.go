package transport

type CreateCategoryDTO struct {
	OperationTypeID int64   `json:"operation_type_id"`
	Title           string  `json:"title"`
	Icon            *string `json:"icon"`
	Color           *string `json:"color"`
}

type UpdateCategoryDTO struct {
	Version         int64   `json:"version"`
	OperationTypeID *int64  `json:"operation_type_id"`
	Title           *string `json:"title"`
	Icon            *string `json:"icon"`
	Color           *string `json:"color"`
}

type ResponseCategoryDTO struct {
	ID              int64  `json:"id"`
	Version         int64  `json:"version"`
	OperationTypeID int64  `json:"operation_type_id"`
	Title           string `json:"title"`
	Icon            string `json:"icon"`
	Color           string `json:"color"`
}

type ResponseCategoriesDTO struct {
	Categories []ResponseCategoryDTO `json:"categories"`
}
