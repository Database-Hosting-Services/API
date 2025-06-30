package ai

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"
	"io"
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
		userID := r.Context().Value("user-id").(int64)

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

func Agent(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.BadRequest(w, "Failed to read request body", err)
			return
		}

		// parse the request body
		var requestBody map[string]interface{}
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			response.BadRequest(w, "Failed to parse request body", err)
			return
		}
		// check if the request body is valid
		if requestBody["prompt"] == nil {
			response.BadRequest(w, "Prompt is required", nil)
			return
		}

		prompt := requestBody["prompt"].(string)
		// get project id from path
		vars := mux.Vars(r)
		projectUID := vars["project_id"]

		// get user id from context
		userID := r.Context().Value("user-id").(int64)

		AIresponse, err := AgentQuery(projectUID, userID, prompt, config.AI)
		if err != nil {
			response.InternalServerError(w, "error while querying agent", err)
			return
		}

		response.OK(w, "Agent query successful", AIresponse)
	}
}

func AgentAccept(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}