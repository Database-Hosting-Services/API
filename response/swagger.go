package response

// SuccessResponse represents a successful API response for Swagger documentation
type SuccessResponse struct {
	Message string      `json:"message,omitempty" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response for Swagger documentation
type ErrorResponse struct {
	Error  string `json:"error" example:"Invalid request parameters"`
}

// ErrorResponse represents an error API response for Swagger documentation
type ErrorResponse400 struct {
	Error  string `json:"error" example:"Invalid request parameters"`
}

type ErrorResponse401 struct {
	Error string `json:"error" example:"Unauthorized"`
}

type ErrorResponse403 struct {
	Error string `json:"error" example:"Forbidden"`
}

type ErrorResponse404 struct {
	Error string `json:"error" example:"Resource not found"`
}

type ErrorResponse500 struct {
	Error string `json:"error" example:"Internal server error"`
}