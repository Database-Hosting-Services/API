package analytics

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// CurrentStorage godoc
// @Summary Get current storage information
// @Description Retrieve the current storage usage information for a specific project
// @Tags analytics
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Current storage retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is missing"
// @Failure 404 {object} response.ErrorResponse "No storage information found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /projects/{project_id}/analytics/storage [get]
func CurrentStorage(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		storage, apiErr := GetALLDatabaseStorage(r.Context(), config.DB, projectOid)
		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}
		if len(storage) == 0 {
			response.NotFound(w, "No storage information found for the project", nil)
			return
		}
		response.OK(w, "Current storage retrieved successfully", storage)
	}
}

// ExecutionTime godoc
// @Summary Get query execution time statistics
// @Description Retrieve statistics about query execution times for a specific project
// @Tags analytics
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Execution time statistics retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is missing"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /projects/{project_id}/analytics/execution-time [get]
func ExecutionTime(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		stats, apiErr := GetALLExecutionTimeStats(r.Context(), config.DB, projectOid)
		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}

		response.OK(w, "Execution time statistics retrieved successfully", stats)
	}
}

// DatabaseUsage godoc
// @Summary Get database usage statistics
// @Description Retrieve statistics about database usage and associated costs
// @Tags analytics
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Database usage statistics retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is missing"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /projects/{project_id}/analytics/usage [get]
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
