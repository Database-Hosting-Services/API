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

func GetDatabaseTableSchema(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// userId := r.Context().Value("user-id").(int)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}
	}
}
