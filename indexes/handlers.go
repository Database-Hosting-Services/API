package indexes

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
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
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		var indexData IndexData
		if err := json.NewDecoder(r.Body).Decode(&indexData); err != nil {
			response.BadRequest(w, r, "Invalid request body", nil)
			return
		}

		if indexData.IndexName == "" || indexData.IndexType == "" || len(indexData.Columns) == 0 || indexData.TableName == "" {
			response.BadRequest(w, r, "Index name, type, columns and table name are required", nil)
			return
		}

		// Create the index in the database
		err := CreateIndexInDatabase(r.Context(), config.DB, projectOid, indexData)

		if err.Error() != nil {
			utils.ResponseHandler(w, r, err)
			return
		}

		response.Created(w, r, "Index created successfully", nil)
	}
}

// ProjectIndexes godoc
// @Summary Get all indexes in a project
// @Description Retrieve a list of all indexes for the specified project
// @Tags indexes
// @Produce json
// @Param project_id path string true "Project ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Indexes retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is required"
// @Failure 500 {object} response.ErrorResponse "Failed to get indexes"
// @Router /projects/{project_id}/indexes [get]
func ProjectIndexes(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		indexes, err := GetIndexes(r.Context(), config.DB, projectOid)

		if err.Error() != nil {
			config.App.ErrorLog.Println("Failed to get indexes:", err)
			utils.ResponseHandler(w, r, err)
			return
		}

		response.OK(w, r, "Indexes retrieved successfully", indexes)
	}
}

// GetIndex godoc
// @Summary Get details of a specific index
// @Description Retrieve details of a specific index by its ID and project ID
// @Tags indexes
// @Produce json
// @Param project_id path string true "Project ID"
// @Param index_oid path string true "Index ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Index retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Index ID and Project ID are required"
// @Failure 404 {object} response.ErrorResponse "Index not found"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to get index"
// @Router /projects/{project_id}/indexes/{index_oid} [get]
func GetIndex(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		indexOid, projectOid := urlVariables["index_oid"], urlVariables["project_id"]
		if indexOid == "" || projectOid == "" {
			response.BadRequest(w, r, "Index Id and Project Id are required", nil)
			return
		}

		index, err := GetSpecificIndex(r.Context(), config.DB, projectOid, indexOid)
		if err.Error() != nil {
			config.App.ErrorLog.Println("Failed to get index:", err)
			utils.ResponseHandler(w, r, err)
			return
		}

		response.OK(w, r, "Index retrieved successfully", index)
	}
}

// DeleteIndex godoc
// @Summary Delete a specific index
// @Description Remove a specific index by its ID and project ID
// @Tags indexes
// @Produce json
// @Param project_id path string true "Project ID"
// @Param index_oid path string true "Index ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Index deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Index ID and Project ID are required"
// @Failure 404 {object} response.ErrorResponse "Index not found"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to delete index"
// @Router /projects/{project_id}/indexes/{index_oid} [delete]
func DeleteIndex(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		indexOid, projectOid := urlVariables["index_oid"], urlVariables["project_id"]
		if indexOid == "" || projectOid == "" {
			response.BadRequest(w, r, "Index Id and Project Id are required", nil)
			return
		}

		err := DeleteSpecificIndex(r.Context(), config.DB, projectOid, indexOid)
		if err.Error() != nil {
			utils.ResponseHandler(w, r, err)
			return
		}

		response.OK(w, r, "Index deleted successfully", nil)
	}
}

// UpdateIndexName godoc
// @Summary Update the name of a specific index
// @Description Modify the name of a specific index by its ID and project ID
// @Tags indexes
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param index_oid path string true "Index ID"
// @Param name body object true "New index name" Schema example: {"name": "new_name"}
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Index name updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid input or index name is the same as the current name"
// @Failure 404 {object} response.ErrorResponse "Index not found"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to update index name"
// @Router /projects/{project_id}/indexes/{index_oid} [put]
func UpdateIndexName(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		indexOid, projectOid := urlVariables["index_oid"], urlVariables["project_id"]
		if indexOid == "" || projectOid == "" {
			response.BadRequest(w, r, "Index Id and Project Id are required", nil)
			return
		}

		var indexData UpdateName
		if err := json.NewDecoder(r.Body).Decode(&indexData); err != nil {
			response.BadRequest(w, r, "Invalid request body", nil)
			return
		}

		if indexData.Name == "" {
			response.BadRequest(w, r, "Index name is required", nil)
			return
		}

		err := UpdateSpecificIndex(r.Context(), config.DB, projectOid, indexOid, indexData.Name)
		if err.Error() != nil {
			config.App.ErrorLog.Println("Failed to update index name:", err)
			utils.ResponseHandler(w, r, err)
			return
		}

		response.OK(w, r, "Index name updated successfully", nil)
	}
}
