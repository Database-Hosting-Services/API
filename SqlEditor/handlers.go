package sqleditor

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ResponseBodySwagger is a swagger documentation model for ResponseBody
// @Description SQL query execution response with results and metadata
type ResponseBodySwagger struct {
	// The JSON result of the query execution
	Result json.RawMessage `json:"result" swaggertype:"string" example:"[{\"id\":1,\"name\":\"test\"}]"`
	// Names of columns in the result set
	ColumnNames []string `json:"column_names" example:"[\"id\",\"name\"]"`
	// Query execution time in milliseconds
	ExecutionTime float64 `json:"execution_time" example:"10.45"`
}

// RunSqlQuery godoc
// @Summary Execute SQL query on project database
// @Description Execute a dynamic SQL query against a specific project's PostgreSQL database and return structured JSON results with metadata including column names and execution time
// @Tags sqlEditor
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID (OID)"
// @Param query body sqleditor.RequestBody true "SQL query to execute"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=sqleditor.ResponseBodySwagger} "Query executed successfully with results, column names, and execution time"
// @Failure 400 {object} response.ErrorResponse "Project ID missing, invalid request body, empty query, or dangerous SQL operations detected"
// @Failure 401 {object} response.ErrorResponse "Unauthorized access"
// @Failure 404 {object} response.ErrorResponse "Project not found or referenced table/column does not exist"
// @Failure 500 {object} response.ErrorResponse "Internal server error, query execution failed, or error parsing results"
// @Router /projects/{project_id}/sqlEditor/run-query [post]
func RunSqlQuery(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		// Get the request body
		var RequestBody RequestBody
		if err := json.NewDecoder(r.Body).Decode(&RequestBody); err != nil {
			response.BadRequest(w, r, "Invalid request body", nil)
			return
		}

		if RequestBody.Query == "" {
			response.BadRequest(w, r, "Query is required", nil)
			return
		}

		// Validate query for dangerous operations
		isValid, err := ValidateQuery(RequestBody.Query)
		if !isValid {
			response.BadRequest(w, r, err.Error(), nil)
			return
		}

		// Get the query response
		queryResponse, apiErr := GetQueryResponse(r.Context(), config.DB, projectOid, RequestBody.Query)
		if apiErr.Error() != nil {
			utils.ResponseHandler(w, r, apiErr)
			return
		}

		response.OK(w, r, "Query executed successfully", queryResponse)
		config.App.InfoLog.Println("Query executed successfully for project:", projectOid)
	}
}
