// models/api_response.go
package models

type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(message string, data interface{}) *ApiResponse {
	return &ApiResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, status string) *ApiResponse {
	return &ApiResponse{
		Status:  status,
		Message: message,
	}
}
