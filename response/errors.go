package response

import "net/http"

// you can use one of these frequently used response for more code readability
//

func BadRequest(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusBadRequest, message, err, nil)
}

func NotFound(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusNotFound, message, err, nil)
}

func InternalServerError(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusInternalServerError, message, err, nil)
}

func UnAuthorized(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusUnauthorized, message, err, nil)
}

func MethodNotAllowed(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusMethodNotAllowed,message, err, nil)
}