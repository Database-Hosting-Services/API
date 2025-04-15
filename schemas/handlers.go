package schemas

import (
	"DBHS/config"
	"DBHS/projects"
	"DBHS/response"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// GetDatabaseSchema godoc
// @Summary Get database schema
// @Description Get the database schema for a specific project
// @Tags schemas
// @Produce json
// @Param project-id path string true "Project ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Database schema retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is required"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/projects/{project-id}/schema/tables [get]
func GetDatabaseSchema(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user-id").(int)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project-id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		project, err := projects.GetUserSpecificProject(r.Context(), config.DB, userId, projectOid)
		if err != nil {
			if errors.Is(err, projects.ErrorProjectNotFound) {
				response.BadRequest(w, "Project is not found", err)
				return
			}
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		databaseConn, err := config.ConfigManager.GetDbConnection(r.Context(), project.Name)
		if err != nil {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		fmt.Println(databaseConn)

		schema, err := getSchema(r.Context(), databaseConn, project.Name)
		if err != nil {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		response.OK(w, "Schema Fetched successfully", schema)
	}
}

// GetDatabaseTableSchema godoc
// @Summary Get table schema
// @Description Get the schema for a specific table in a project
// @Tags schemas
// @Produce json
// @Param project-id path string true "Project ID"
// @Param table-id path string true "Table ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Table schema retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is required"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/projects/{project-id}/schema/tables/{table-id} [get]
func GetDatabaseTableSchema(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// userId := r.Context().Value("user-id").(int)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project-id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}
	}
}
