package response

import "net/http"

// you can use one of these frequently used response for more code readability

func OK(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	CreateResponse(w, r, http.StatusOK, message, nil, data, nil)
}

func Created(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	CreateResponse(w, r, http.StatusCreated, message, nil, data, nil)
}

func Redirect(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	CreateResponse(w, r, http.StatusFound, message, nil, data, nil) // 302
}
