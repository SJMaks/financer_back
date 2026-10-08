package transport

type ResponseOperationTypeDTO struct {
	ID    int64  `json:"id" example:"1"`
	Title string `json:"title" example:"income"`
}

type ResponseOperationTypesDTO struct {
	OperationTypes []ResponseOperationTypeDTO `json:"operation_types"`
}
