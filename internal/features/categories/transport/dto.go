package transport

type CreateCategoryDTO struct {
	OperationTypeID int64   `json:"operation_type_id" example:"1"`
	Title           string  `json:"title" example:"Salary"`
	Icon            *string `json:"icon" example:"dollar"`
	Color           *string `json:"color" example:"#26441e"`
}

type UpdateCategoryDTO struct {
	Version         *int64  `json:"version" example:"0"`
	OperationTypeID *int64  `json:"operation_type_id" example:"2"`
	Title           *string `json:"title" example:"Groceries"`
	Icon            *string `json:"icon" example:"apple"`
	Color           *string `json:"color" example:"#960606"`
}

type ResponseCategoryDTO struct {
	ID              int64  `json:"id" example:"1"`
	Version         int64  `json:"version" example:"0"`
	OperationTypeID int64  `json:"operation_type_id" example:"1"`
	Title           string `json:"title" example:"Salary"`
	Icon            string `json:"icon" example:"dollar"`
	Color           string `json:"color" example:"#26441e"`
}

type ResponseCategoriesDTO struct {
	Categories []ResponseCategoryDTO `json:"categories"`
}
