package api

// ApiError represents an operational error with an HTTP status code and additional details.
type ApiError struct {
	Message    string
	StatusCode int
	err        error
}

// NewApiError creates a new ApiError instance.
// It sets the Status field based on the status code: "fail" for 4xx errors, otherwise "error".
func NewApiError(message string, statusCode int, err error) *ApiError {
	return &ApiError{
		Message:    message,
		StatusCode: statusCode,
		err:        err,
	}
}

// Error implements the error interface for ApiError.
func (e *ApiError) Error() error {
	return e.err
}
