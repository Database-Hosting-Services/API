package ai

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"

	//"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
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

func ChatBotAsk(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		projectOID := vars["project_id"]
		projectID, err := GetProjectIDfromOID(r.Context(), config.DB, projectOID)
		if err != nil {
			response.InternalServerError(w, "Failed to get project ID", err)
			return
		}
		userID := r.Context().Value("user-id").(int)

		var userRequest ChatBotRequest
		if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
			response.BadRequest(w, "Invalid request body", err)
			return
		}

		transaction, err := config.DB.Begin(r.Context())
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Failed to start database transaction", err)
			return
		}

		chat_data, err := GetOrCreateChatData(r.Context(), transaction, userID, projectID)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Failed to get or create chat data", err)
			return
		}

		answer, err := config.AI.QueryChat(userRequest.Question)
		if err != nil {
			response.InternalServerError(w, err.Error(), err)
			return
		}

		err = SaveChatAction(r.Context(), transaction, chat_data.ID, userID, userRequest.Question, answer.ResponseText)
		if err != nil {
			response.InternalServerError(w, err.Error(), err)
			return
		}

		if err := transaction.Commit(r.Context()); err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Failed to commit database transaction", err)
			return
		}

		response.OK(w, "Answer generated successfully", answer)
	}
}
