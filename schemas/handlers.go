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

		projectName := strings.ToLower(project.Name)
		projectName += "_" + strconv.Itoa(userId)

		databaseConn, err := config.ConfigManager.GetDbConnection(r.Context(), projectName)
		if err != nil {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		schema, err := getDatabaseSchema(r.Context(), databaseConn)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		response.OK(w, "Schema Fetched successfully", schema)
	}
}

func GetDatabaseTableSchema(app *config.Application) http.HandlerFunc {
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

		projectName := strings.ToLower(project.Name)
		projectName += "_" + strconv.Itoa(userId)

		databaseConn, err := config.ConfigManager.GetDbConnection(r.Context(), projectName)
		if err != nil {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		tableOID := urlVariables["table-id"]
		tableName, err := getDatabaseTableName(r.Context(), config.DB, tableOID)
		if err != nil {
			print(err.Error())
			response.BadRequest(w, "Invalid Table Id", err)
			return
		}

		schema, err := GetTableSchema(r.Context(), databaseConn, tableName)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			config.App.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		response.OK(w, "Schema Fetched successfully", schema)
	}
}
