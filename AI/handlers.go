package ai

import (
	"DBHS/config"
	"DBHS/response"
	"net/http"
	"time"

	"github.com/axiomhq/axiom-go/axiom"
	"github.com/axiomhq/axiom-go/axiom/ingest"
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
		userID := r.Context().Value("user-id").(int64)

		Analytics := getAnalytics() // TODO: get real analytics
		AI := config.AI

		report, err := getReport(projectID, userID, Analytics, AI)
		if err != nil {
			response.InternalServerError(w, err.Error(), err)
			config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
				{
					ingest.TimestampField: time.Now(),
					"project_id": projectID,
					"user_id":    userID,
					"error":      err.Error(),
					"message":    "Failed to generate AI report",
				},
			})
			return
		}

		config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
			{
				ingest.TimestampField: time.Now(),
				"project_id": projectID,
				"user_id":    userID,
				"report":     report,
				"status":     "success",
				"message":    "AI report generated successfully",
			},
		})

		response.OK(w, "Report generated successfully", report)
	}
}