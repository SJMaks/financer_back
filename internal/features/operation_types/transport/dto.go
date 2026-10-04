package transport

type ResponseOperationTypeDTO struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type ResponseOperationTypesDTO struct {
	OperationTypes []ResponseOperationTypeDTO `json:"operation_types"`
}
