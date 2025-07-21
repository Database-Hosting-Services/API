package analytics

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// CurrentStorage godoc
// @Summary Get historical storage information
// @Description Retrieve all historical storage usage records for a specific project with timestamps
// @Tags analytics
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID (OID)"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]analytics.StorageWithDates} "Storage history retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is missing"
// @Failure 404 {object} response.ErrorResponse "No storage information found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /projects/{project_id}/analytics/storage [get]
func CurrentStorage(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		storage, apiErr := GetALLDatabaseStorage(r.Context(), config.DB, projectOid)
		if len(storage) == 0 || (apiErr.Error() != nil && strings.Contains(apiErr.Error().Error(), "cannot scan NULL into")) {
			response.NotFound(w, r, "No storage information found for the project", nil)
			return
		}

		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}
		
		response.OK(w, r, "Storage history retrieved successfully", storage)
	}
}

// ExecutionTime godoc
// @Summary Get query execution time statistics
// @Description Retrieve all historical statistics about query execution times for a specific project with timestamps
// @Tags analytics
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID (OID)"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]analytics.DatabaseActivityWithDates} "Execution time statistics retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is missing"
// @Failure 404 {object} response.ErrorResponse "No execution time records found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /projects/{project_id}/analytics/execution-time [get]
func ExecutionTime(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		stats, apiErr := GetALLExecutionTimeStats(r.Context(), config.DB, projectOid)
        if len(stats) == 0 || (apiErr.Error() != nil && strings.Contains(apiErr.Error().Error(), "cannot scan NULL into")) {
			response.NotFound(w, r, "No execution time records found for the project", nil)
			return
		}

		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}

		response.OK(w, r, "Execution time statistics retrieved successfully", stats)
	}
}

// DatabaseUsage godoc
// @Summary Get database usage statistics and costs
// @Description Retrieve all historical statistics about database usage and associated costs for a specific project with timestamps
// @Tags analytics
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID (OID)"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]analytics.DatabaseUsageCostWithDates} "Database usage statistics retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is missing"
// @Failure 404 {object} response.ErrorResponse "No database usage records found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /projects/{project_id}/analytics/usage [get]
func DatabaseUsage(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		stats, apiErr := GetALLDatabaseUsageStats(r.Context(), config.DB, projectOid)
		if len(stats) == 0 || (apiErr.Error() != nil && strings.Contains(apiErr.Error().Error(), "cannot scan NULL into")) {
			response.NotFound(w, r, "No database usage records found for the project", nil)
			return
		}

		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}

		response.OK(w, r, "Database usage statistics retrieved successfully", stats)
	}
}
