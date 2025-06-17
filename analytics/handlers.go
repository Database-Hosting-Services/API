package analytics

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func CurrentStorage(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		storage, apiErr := GetDatabaseStorage(r.Context(), config.DB, projectOid)
		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}
		if storage.ManagementStorage == "" && storage.ActualData == "" {
			response.NotFound(w, "No storage information found for the project", nil)
			return
		}
		response.OK(w, "Current storage retrieved successfully", storage)
	}
}
