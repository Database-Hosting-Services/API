package response

import "net/http"

func ErrorResponse(w http.ResponseWriter, status int, message string, err error) {
	response := Response{
		Status:  status,
		Message: message,
		Error:   err.Error(),
	}
	sendResponse(w, status, response)
}

// you can use one of these frequently used response for more code readability
//

func BadRequest(w http.ResponseWriter, message string, err error) {
	ErrorResponse(w, http.StatusBadRequest, message, err)
}

func NotFound(w http.ResponseWriter, message string, err error) {
	ErrorResponse(w, http.StatusNotFound, message, err)
}

func InternalServerError(w http.ResponseWriter, message string, err error) {
	ErrorResponse(w, http.StatusInternalServerError, message, err)
}
