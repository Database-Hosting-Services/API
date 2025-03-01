package response

import "net/http"

func SuccessResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	sendResponse(w, status, response)
}

// you can use one of these frequently used response for more code readability
//

func OK(w http.ResponseWriter, message string, data interface{}) {
	SuccessResponse(w, http.StatusOK, message, data)
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	SuccessResponse(w, http.StatusCreated, message, data)
}
