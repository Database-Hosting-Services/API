package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// sendResponse sends the response with custom status
// we user separated status parameters to separate the HTTP metadata and response body (response struct)
// and also to make sure that the status exists (http.ResponseWriter  requires to set the HTTP status code)

// Note : you can make the response.Data as map[string]interface{}
// you can view accounts.service & accounts.handlers for more details
func sendResponse(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
