package ai

import (
	"DBHS/config"
	"DBHS/response"
	"net/http"
	"github.com/gorilla/mux"
)

func getAnalytics() Analytics { // this only a placeholder for now
	return Analytics{}
}


// @Summary Generate AI Report
// @Description Generate an AI-powered analytics report for a specific project
// @Tags AI
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} response.Response{data=object} "Report generated successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /projects/{project_id}/ai/report [get]
// @Security BearerAuth
func Report(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get project id from path
		vars := mux.Vars(r)
		projectID := vars["project_id"]

		// get user id from context
		userID := r.Context().Value("user-id").(int)

		Analytics := getAnalytics() // TODO: get real analytics
		AI := config.AI

		report, err := getReport(projectID, userID, Analytics, AI)
		if err != nil {
			response.InternalServerError(w, err.Error(), err)
			return
		}

		response.OK(w, "Report generated successfully", report)
	}
}