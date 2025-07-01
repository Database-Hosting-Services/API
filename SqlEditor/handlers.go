package sqleditor

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RunSqlQuery(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		// Get the request body
		var RequestBody RequestBody
		if err := json.NewDecoder(r.Body).Decode(&RequestBody); err != nil {
			response.BadRequest(w, "Invalid request body", nil)
			return
		}

		if RequestBody.Query == "" {
			response.BadRequest(w, "Query is required", nil)
			return
		}

		// Get the query response
		queryResponse, apiErr := GetQueryResponse(r.Context(), config.DB, projectOid, RequestBody.Query)
		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}

		response.OK(w, "Query executed successfully", queryResponse)
		config.App.InfoLog.Println("Query executed successfully for project:", projectOid)
	}
}
