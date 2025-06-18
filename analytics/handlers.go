package analytics

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"net/http"

	"github.com/gorilla/mux"
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

// ExecutionTime returns statistics about query execution times for a project
func ExecutionTime(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		stats, apiErr := GetExecutionTimeStats(r.Context(), config.DB, projectOid)
		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}

		response.OK(w, "Execution time statistics retrieved successfully", stats)
	}
}

// DatabaseUsage returns statistics about database usage and associated costs
func DatabaseUsage(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		stats, apiErr := GetDatabaseUsageStats(r.Context(), config.DB, projectOid)
		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}

		response.OK(w, "Database usage statistics retrieved successfully", stats)
	}
}
