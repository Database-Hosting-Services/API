package indexes

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func CreateIndex(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		var indexData IndexData
		if err := json.NewDecoder(r.Body).Decode(&indexData); err != nil {
			response.BadRequest(w, "Invalid request body", nil)
			return
		}

		if indexData.IndexName == "" || indexData.IndexType == "" || len(indexData.Columns) == 0 || indexData.TableName == "" {
			response.BadRequest(w, "Index name, type, columns and table name are required", nil)
			return
		}

		// Create the index in the database
		err := CreateIndexInDatabase(r.Context(), config.DB, projectOid, indexData)
		if err != nil {
			app.ErrorLog.Println("Failed to create index:", err)
			if strings.Contains(err.Error(), "already exists") {
				response.BadRequest(w, "Index already exists", nil)
			} else if strings.Contains(err.Error(), "project not found") {
				response.NotFound(w, "Project not found", nil)
			} else if strings.Contains(err.Error(), "Unauthorized") {
				response.UnAuthorized(w, "Unauthorized", nil)
			} else {
				response.InternalServerError(w, "Failed to create index", nil)
			}
			return
		}

		response.Created(w, "Index created successfully", nil)
	}
}

func ProjectIndexes(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		indexes, err := GetIndexes(r.Context(), config.DB, projectOid)
		if err != nil {
			app.ErrorLog.Println("Failed to get indexes:", err)
			response.InternalServerError(w, "Failed to get indexes", nil)
			return
		}

		response.OK(w, "Indexes retrieved successfully", indexes)
	}
}
