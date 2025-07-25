package schemas

import (
	"DBHS/config"
	"DBHS/projects"
	"DBHS/response"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
	"strings"
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
		userId := r.Context().Value("user-id").(int64)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project-id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		project, err := projects.GetUserSpecificProject(r.Context(), config.DB, userId, projectOid)
		if err != nil {
			if errors.Is(err, projects.ErrorProjectNotFound) {
				response.BadRequest(w, r, "Project is not found", err)
				return
			}
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		projectName := strings.ToLower(project.Name)
		projectName += "_" + strconv.FormatInt(userId, 10)

		databaseConn, err := config.ConfigManager.GetDbConnection(r.Context(), projectName)
		if err != nil {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		schema, err := getDatabaseSchema(r.Context(), databaseConn)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		response.OK(w, r, "Schema Fetched successfully", schema)
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
		userId := r.Context().Value("user-id").(int64)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project-id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		project, err := projects.GetUserSpecificProject(r.Context(), config.DB, userId, projectOid)
		if err != nil {
			if errors.Is(err, projects.ErrorProjectNotFound) {
				response.BadRequest(w, r, "Project is not found", err)
				return
			}
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		projectName := strings.ToLower(project.Name)
		projectName += "_" + strconv.FormatInt(userId, 10)

		databaseConn, err := config.ConfigManager.GetDbConnection(r.Context(), projectName)
		if err != nil {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		tableOID := urlVariables["table-id"]
		tableName, err := getDatabaseTableName(r.Context(), config.DB, tableOID)
		if err != nil {
			print(err.Error())
			response.BadRequest(w, r, "Invalid Table Id", err)
			return
		}

		schema, err := GetTableSchema(r.Context(), databaseConn, tableName)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		response.OK(w, r, "Schema Fetched successfully", schema)
	}
}
