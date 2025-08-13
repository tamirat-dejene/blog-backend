package domain

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
	Code  int    `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
