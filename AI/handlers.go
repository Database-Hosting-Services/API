package ai

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"
	"time"

	"github.com/axiomhq/axiom-go/axiom"
	"github.com/axiomhq/axiom-go/axiom/ingest"
	"github.com/gorilla/mux"
	"io"
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
		userID := r.Context().Value("user-id").(int64)

		Analytics := getAnalytics() // TODO: get real analytics
		AI := config.AI

		report, err := getReport(projectID, userID, Analytics, AI)
		if err != nil {
			response.InternalServerError(w, err.Error(), err)
			config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
				{
					ingest.TimestampField: time.Now(),
					"project_id":          projectID,
					"user_id":             userID,
					"error":               err.Error(),
					"message":             "Failed to generate AI report",
				},
			})
			return
		}

		config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
			{
				ingest.TimestampField: time.Now(),
				"project_id":          projectID,
				"user_id":             userID,
				"report":              report,
				"status":              "success",
				"message":             "AI report generated successfully",
			},
		})

		response.OK(w, "Report generated successfully", report)
	}
}

// ChatBotAsk godoc
// @Summary Chat Bot Ask
// @Description This endpoint allows users to ask questions to the chatbot, which will respond using AI. It also saves the chat history for future reference.
// @Tags AI
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param ChatBotRequest body ChatBotRequest true "Chat Bot Request"
// @Success 200 {object} response.Response{data=object} "Answer generated successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /projects/{project_id}/ai/chatbot/ask [post]
// @Security BearerAuth
func ChatBotAsk(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		projectOID := vars["project_id"]
		projectID, err := GetProjectIDfromOID(r.Context(), config.DB, projectOID)
		if err != nil {
			response.InternalServerError(w, "Failed to get project ID", err)
			return
		}
		userID64 := r.Context().Value("user-id").(int64)
		userID := int(userID64)

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

		// i should ignore this step if the client passed the chat history with the request
		chat_data, err := GetOrCreateChatData(r.Context(), transaction, userID, projectID)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Failed to get or create chat data", err)
			return
		}
		app.InfoLog.Printf("Chat data: %+v", chat_data)

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

func Agent(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.BadRequest(w, "Failed to read request body", err)
			return
		}

		// parse the request body
		request := &Request{}
		err = json.Unmarshal(body, request)
		if err != nil {
			response.BadRequest(w, "Failed to parse request body", err)
			return
		}
		// check if the request body is valid
		if request.Prompt == "" {
			response.BadRequest(w, "Prompt is required", nil)
			return
		}

		// get project id from path
		vars := mux.Vars(r)
		projectUID := vars["project_id"]

		// get user id from context
		userID := r.Context().Value("user-id").(int64)

		AIresponse, err := AgentQuery(projectUID, userID, request.Prompt, config.AI)
		if err != nil {
			response.InternalServerError(w, "error while querying agent", err)
			return
		}

		response.OK(w, "Agent query successful", AIresponse)
	}
}

func AgentAccept(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get project id from path
		vars := mux.Vars(r)
		projectUID := vars["project_id"]
		// get user id from context
		userID := r.Context().Value("user-id").(int64)

		// execute the agent query
		err := AgentExec(projectUID, userID, config.AI)
		if err != nil {
			if err.Error() == "changes expired or not found" {
				response.BadRequest(w, "No schema changes found or changes expired", nil)
			} else {
				response.InternalServerError(w, "error while executing agent", err)
			}
			// log the error to Axiom
			config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
				{
					ingest.TimestampField: time.Now(),
					"project_id":          projectUID,
					"user_id":             userID,
					"error":               err.Error(),
					"message":             "Failed to execute agent query",
				},
			})
			return
		}
		// log the success to Axiom
		config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
			{
				ingest.TimestampField: time.Now(),
				"project_id":          projectUID,
				"user_id":             userID,
				"message":             "Agent query executed successfully",
			},
		})
		response.OK(w, "query executed successfully", nil)
	}
}
