package response

import "net/http"

// you can use one of these frequently used response for more code readability
//

func OK(w http.ResponseWriter, message string, data interface{}) {
	CreateResponse(w, http.StatusOK, message, nil, data)
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	CreateResponse(w, http.StatusCreated, message, nil, data)
}
