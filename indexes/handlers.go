package indexes

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// CreateIndex godoc
// @Summary Create a new index
// @Description Create a new index in the specified project
// @Tags indexes
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param index body IndexData true "Index information"
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse "Index created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input or index already exists"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Project not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /projects/{project_id}/indexes [post]
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
			config.App.ErrorLog.Println("Failed to create index:", err)
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
			config.App.ErrorLog.Println("Failed to get indexes:", err)
			response.InternalServerError(w, "Failed to get indexes", nil)
			return
		}

		response.OK(w, "Indexes retrieved successfully", indexes)
	}
}

func GetIndex(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		indexOid, projectOid := urlVariables["index_oid"], urlVariables["project_id"]
		if indexOid == "" || projectOid == "" {
			response.BadRequest(w, "Index Id and Project Id are required", nil)
			return
		}

		index, err := GetSpecificIndex(r.Context(), config.DB, projectOid, indexOid)
		if err != nil {
			config.App.ErrorLog.Println("Failed to get index:", err)
			if err.Error() == "index not found" {
				response.NotFound(w, "Index not found", nil)
			} else if err.Error() == "unauthorized" {
				response.UnAuthorized(w, "Unauthorized", nil)
			} else {
				response.InternalServerError(w, "Failed to get index", nil)
			}
			return
		}

		response.OK(w, "Index retrieved successfully", index)
	}
}

func DeleteIndex(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		indexOid, projectOid := urlVariables["index_oid"], urlVariables["project_id"]
		if indexOid == "" || projectOid == "" {
			response.BadRequest(w, "Index Id and Project Id are required", nil)
			return
		}

		err := DeleteSpecificIndex(r.Context(), config.DB, projectOid, indexOid)
		if err != nil {
			config.App.ErrorLog.Println("Failed to delete index:", err)
			if strings.Contains(err.Error(), "index not found") {
				response.NotFound(w, "Index not found", nil)
			} else if strings.Contains(err.Error(), "unauthorized") {
				response.UnAuthorized(w, "Unauthorized", nil)
			} else {
				response.InternalServerError(w, "Failed to delete index", nil)
			}
			return
		}

		response.OK(w, "Index deleted successfully", nil)
	}
}

func UpdateIndexName(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		indexOid, projectOid := urlVariables["index_oid"], urlVariables["project_id"]
		if indexOid == "" || projectOid == "" {
			response.BadRequest(w, "Index Id and Project Id are required", nil)
			return
		}

		var indexData IndexData
		if err := json.NewDecoder(r.Body).Decode(&indexData); err != nil {
			response.BadRequest(w, "Invalid request body", nil)
			return
		}

		if indexData.IndexName == "" {
			response.BadRequest(w, "Index name is required", nil)
			return
		}

		err := UpdateSpecificIndex(r.Context(), config.DB, projectOid, indexOid, indexData.IndexName)
		if err != nil {
			config.App.ErrorLog.Println("Failed to update index name:", err)
			if strings.Contains(err.Error(), "index not found") {
				response.NotFound(w, "Index not found", nil)
			} else if strings.Contains(err.Error(), "unauthorized") {
				response.UnAuthorized(w, "Unauthorized", nil)
			} else if strings.Contains(err.Error(), "index name is the same as the current name") {
				response.BadRequest(w, "Index name is the same as the current name", nil)
			} else if strings.Contains(err.Error(), "index already exists") {
				response.BadRequest(w, "Index already exists", nil)
			} else {
				response.InternalServerError(w, "Failed to update index name", nil)
			}
			return
		}

		response.OK(w, "Index name updated successfully", nil)
	}
}
