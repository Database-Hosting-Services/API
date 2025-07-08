package response

import (
	"DBHS/config"
	"encoding/json"
	"net/http"
	"time"

	"github.com/axiomhq/axiom-go/axiom"
	"github.com/axiomhq/axiom-go/axiom/ingest"
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

func CreateResponse(w http.ResponseWriter, r *http.Request, status int, message string, err error, data interface{}, headers map[string]string) {
	var response *Response
	event := axiom.Event{
		ingest.TimestampField: time.Now(),
		"user-id":   r.Context().Value("user-id"),
		"user-oid":  r.Context().Value("user-oid"),
		"user-name": r.Context().Value("user-name"),
		"status-code": status,
		"URI": r.RequestURI,
	}
	if err != nil {
		response = &Response{
			Status: status,
			Error:  err.Error(),
		}
		// log the error to axiom
		event["error"] = err.Error()
	} else {
		response = &Response{
			Status: status,
			Data:   data,
		}
		event["response"] = response
	}
	if message != "" {
		response.Message = message
	}
	event["massage"] = response.Message
	config.AxiomLogger.IngestEvents(r.Context(), "api", []axiom.Event{event})
	SendResponse(w, status, headers, response)
}
