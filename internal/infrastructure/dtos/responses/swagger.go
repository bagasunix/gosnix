package responses

// BaseResponseCustomer digunakan khusus untuk dokumentasi Swagger

// 400 Bad Request
type ErrorBadRequestResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
	Error   string `json:"error" example:"invalid request body"`
}

// 401 Unauthorized
type ErrorUnauthorizedResponse struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"Unauthorized"`
	Error   string `json:"error" example:"missing or invalid token"`
}

// 403 Forbidden
type ErrorForbiddenResponse struct {
	Code    int    `json:"code" example:"403"`
	Message string `json:"message" example:"Forbidden"`
	Error   string `json:"error" example:"access denied"`
}

// 404 Not Found
type ErrorNotFoundResponse struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"Not Found"`
	Error   string `json:"error" example:"resource not found"`
}

// 409 Conflict
type ErrorConflictResponse struct {
	Code    int    `json:"code" example:"409"`
	Message string `json:"message" example:"Conflict"`
	Error   string `json:"error" example:"duplicate data"`
}

// 422 Unprocessable Entity (validation)
type ErrorUnprocessableEntityResponse struct {
	Code    int    `json:"code" example:"422"`
	Message string `json:"message" example:"Unprocessable Entity"`
	Error   string `json:"error" example:"validation failed"`
}

// 500 Internal Server Error
type ErrorInternalServerResponse struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
	Error   string `json:"error" example:"unexpected error occurred"`
}

// 200 OK (single data)
type SuccessOKResponse[T any] struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
	Data    *T     `json:"data"`
}

// 201 Created
type SuccessCreatedResponse[T any] struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"Created"`
	Data    *T     `json:"data"`
}

// 202 Accepted (async processing)
type SuccessAcceptedResponse struct {
	Code    int    `json:"code" example:"202"`
	Message string `json:"message" example:"Accepted"`
	Info    string `json:"info" example:"request is being processed"`
}

// 204 No Content (biasanya untuk delete, tidak ada data)
type SuccessNoContentResponse struct {
	Code    int    `json:"code" example:"204"`
	Message string `json:"message" example:"No Content"`
}

type CustomerRegisterResponseWrapper SuccessCreatedResponse[CustomerResponse]
