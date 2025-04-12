package response

import (
	"errors"
	"net/http"
)


var (
	 ErrUnauthorized = errors.New("Unauthorized")
)

// you can use one of these frequently used response for more code readability
//

func BadRequest(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusBadRequest, message, err, nil, nil)
}

func NotFound(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusNotFound, message, err, nil, nil)
}

func InternalServerError(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusInternalServerError, message, err, nil, nil)
}

func UnAuthorized(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusUnauthorized, message, err, nil, nil)
}

func MethodNotAllowed(w http.ResponseWriter, allowed string, message string, err error) {
	CreateResponse(w, http.StatusMethodNotAllowed, message, err, nil, map[string]string{
		"Allow": allowed,
	})
}

func TooManyRequests(w http.ResponseWriter, message string, err error) {
	CreateResponse(w, http.StatusTooManyRequests, message, err, nil, nil)
}
