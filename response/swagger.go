package response

// SuccessResponse represents a successful API response for Swagger documentation
type SuccessResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response for Swagger documentation
type ErrorResponse struct {
	Status int    `json:"status" example:"400"`
	Error  string `json:"error" example:"Invalid request parameters"`
}
