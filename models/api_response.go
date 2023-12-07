// models/api_response.go
package models

type ApiSuccesResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ApiErrorResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	ErrorDetails any    `json:"error_details,omitempty"`
}

func NewSuccessResponse(message string, data any) *ApiSuccesResponse {
	return &ApiSuccesResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, status string, error any) *ApiErrorResponse {
	return &ApiErrorResponse{
		Status:       status,
		Message:      message,
		ErrorDetails: error,
	}
}

type RegisterResponse struct {
	Id       string
	Username string
	Email    string
}

type LoginResponse struct {
	Id       string
	Username string
	Email    string
	Token    string
}
