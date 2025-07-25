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
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
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
			response.InternalServerError(w, r, err.Error(), err)
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

		response.OK(w, r, "Report generated successfully", report)
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
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
// @Router /projects/{project_id}/ai/chatbot/ask [post]
// @Security BearerAuth
func ChatBotAsk(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		projectOID := vars["project_id"]
		projectID, err := GetProjectIDfromOID(r.Context(), config.DB, projectOID)
		if err != nil {
			response.InternalServerError(w, r, "Failed to get project ID", err)
			return
		}
		userID64 := r.Context().Value("user-id").(int64)
		userID := int(userID64)

		var userRequest ChatBotRequest
		if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
			response.BadRequest(w, r, "Invalid request body", err)
			return
		}

		transaction, err := config.DB.Begin(r.Context())
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, r, "Failed to start database transaction", err)
			return
		}

		// i should ignore this step if the client passed the chat history with the request
		chat_data, err := GetOrCreateChatData(r.Context(), transaction, userID, projectID)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, r, "Failed to get or create chat data", err)
			return
		}
		app.InfoLog.Printf("Chat data: %+v", chat_data)

		answer, err := config.AI.QueryChat(userRequest.Question)
		if err != nil {
			response.InternalServerError(w, r, err.Error(), err)
			return
		}

		err = SaveChatAction(r.Context(), transaction, chat_data.ID, userID, userRequest.Question, answer.ResponseText)
		if err != nil {
			response.InternalServerError(w, r, err.Error(), err)
			return
		}

		if err := transaction.Commit(r.Context()); err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, r, "Failed to commit database transaction", err)
			return
		}

		response.OK(w, r, "Answer generated successfully", answer)
	}
}

// Agent godoc
// @Summary AI Agent Query
// @Description This endpoint allows users to query the AI agent with a prompt. The agent will respond with a schema change suggestion based on the prompt.
// @Tags AI
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param Request body Request true "Request"
// @Success 200 {object} response.Response{data=AgentResponse} "Agent query successful"
// @Failure 400 {object} response.ErrorResponse400 "Bad request"
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
// @Router /projects/{project_id}/ai/agent [post]
// @Security JWTAuth
func Agent(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.BadRequest(w, r, "Failed to read request body", err)
			return
		}

		// parse the request body
		request := &Request{}
		err = json.Unmarshal(body, request)
		if err != nil {
			response.BadRequest(w, r, "Failed to parse request body", err)
			return
		}
		// check if the request body is valid
		if request.Prompt == "" {
			response.BadRequest(w, r, "Prompt is required", nil)
			return
		}

		// get project id from path
		vars := mux.Vars(r)
		projectUID := vars["project_id"]

		// get user id from context
		userID := r.Context().Value("user-id").(int64)

		AIresponse, err := AgentQuery(projectUID, userID, request.Prompt, config.AI)
		if err != nil {
			response.InternalServerError(w, r, "error while querying agent", err)
			// log the error to Axiom
			config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
				{
					ingest.TimestampField: time.Now(),
					"project_id":          projectUID,
					"user_id":             userID,
					"error":               err.Error(),
					"message":             "Failed to query agent",
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
				"message":             "Agent query successful",
				"ddl":                 AIresponse.SchemaDDL,
				"prompt":              request.Prompt,
				"SchemaChanges":       AIresponse.SchemaChanges,
				"response":            AIresponse.Response,
			},
		})

		response.OK(w, r, "Agent query successful", AIresponse)
	}
}

// AgentAccept godoc
// @Summary Accept AI Agent Query
// @Description This endpoint allows users to accept the AI agent's query and execute the schema changes
// @Tags AI
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} response.Response "Query executed successfully"
// @Failure 400 {object} response.ErrorResponse400 "Bad request"
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
// @Router /projects/{project_id}/ai/agent/accept [post]
// @Security JWTAuth
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
				response.BadRequest(w, r, "No schema changes found or changes expired", nil)
			} else {
				response.InternalServerError(w, r, "error while executing agent", err)
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
		response.OK(w, r, "query executed successfully", nil)
	}
}

// AgentCancel godoc
// @Summary Cancel AI Agent Query
// @Description This endpoint allows users to cancel an AI agent query
// @Tags AI
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} response.Response "Agent query cancelled successfully"
// @Failure 400 {object} response.ErrorResponse400 "Bad request"
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
// @Router /projects/{project_id}/ai/agent/cancel [post]
// @Security JWTAuth
func AgentCancel(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get project id from path
		vars := mux.Vars(r)
		projectUID := vars["project_id"]
		// get user id from context
		userID := r.Context().Value("user-id").(int64)
		// cancel the agent query
		err := ClearCacheForProject(projectUID)
		if err != nil {
			response.InternalServerError(w, r, "error while cancelling agent query", err)
			// log the error to Axiom
			config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
				{
					ingest.TimestampField: time.Now(),
					"project_id":          projectUID,
					"user_id":             userID,
					"error":               err.Error(),
					"message":             "Failed to cancel agent query",
				},
			})
			return
		}
		// log the cancellation to Axiom
		config.AxiomLogger.IngestEvents(r.Context(), "ai-logs", []axiom.Event{
			{
				ingest.TimestampField: time.Now(),
				"project_id":          projectUID,
				"user_id":             userID,
				"message":             "Agent query cancelled",
			},
		})

		response.OK(w, r, "Agent query cancelled successfully", nil)
	}
}
