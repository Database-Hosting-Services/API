package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`	
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// sendResponse sends the response with custom status
// we user separated status parameters to separate the HTTP metadata and response body (response struct)
// and also to make sure that the status exists (http.ResponseWriter  requires to set the HTTP status code)

// Note : you can make the response.Data as map[string]interface{}
// you can view accounts.service & accounts.handlers for more details

func SendResponse(w http.ResponseWriter, status int, headers map[string]string, response *Response) {
	w.Header().Set("Content-Type", "application/json")
	for k, v := range headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func CreateResponse(w http.ResponseWriter, status int, message string, err error, data interface{}, headers map[string]string) {
	var response *Response
	if err != nil {
		response = &Response{
			Status:  status,
			Error:   err.Error(),
		}
	} else {
		response = &Response{
			Status:  status,
			Data:    data,
		}
	}
	if message != "" && err != nil {
		response.Message = message
	}
	SendResponse(w, status, headers, response)
}
