package response

import (
	"errors"
	"net/http"
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
	ErrBadRequest = errors.New("BadRequest")
)

// you can use one of these frequently used response for more code readability
//

func BadRequest(w http.ResponseWriter, r *http.Request, message string, err error) {
	CreateResponse(w, r, http.StatusBadRequest, message, err, nil, nil)
}

func NotFound(w http.ResponseWriter, r *http.Request, message string, err error) {
	CreateResponse(w, r, http.StatusNotFound, message, err, nil, nil)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, message string, err error) {
	CreateResponse(w, r, http.StatusInternalServerError, message, err, nil, nil)
}

func UnAuthorized(w http.ResponseWriter, r *http.Request, message string, err error) {
	CreateResponse(w, r, http.StatusUnauthorized, message, err, nil, nil)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request, allowed string, message string, err error) {
	CreateResponse(w, r, http.StatusMethodNotAllowed, message, err, nil, map[string]string{
		"Allow": allowed,
	})
}

func TooManyRequests(w http.ResponseWriter, r *http.Request, message string, err error) {
	CreateResponse(w, r, http.StatusTooManyRequests, message, err, nil, nil)
}
