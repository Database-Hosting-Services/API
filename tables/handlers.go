package tables

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)


// GetAllTablesHanlder godoc
// @Summary Get all tables in a project
// @Description Get a list of all tables in the specified project
// @Tags tables
// @Produce json
// @Param project_id path string true "Project ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]Table} "List of tables"
// @Failure 404 {object} response.ErrorResponse404 "Project not found"
// @Failure 401 {object} response.ErrorResponse401 "Unauthorized"
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
// @Router /api/projects/{project_id}/tables [get]
func GetAllTablesHanlder(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectId := urlVariables["project_id"]
		if projectId == "" {
			response.NotFound(w, "Project ID is required", nil)
			return
		}
		
		data, err := GetAllTables(r.Context(), projectId, config.DB)
		if err != nil {
			if errors.Is(err, response.ErrUnauthorized) {
				response.UnAuthorized(w, "Unauthorized", nil)
				return
			}
			app.ErrorLog.Println("Tables reading failed:", err)
			response.InternalServerError(w, "Failed to read tables", err)
			return
		}

		response.OK(w, "", data)
	}
}

// CreateTableHandler godoc
// @Summary Create a new table
// @Description Create a new table in the specified project
// @Tags tables
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param table body Table true "Table information"
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/projects/{project_id}/tables [post]
func CreateTableHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for creating a table
		table := Table{}
		// Parse the request body to populate the table struct
		if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
			response.BadRequest(w, "Invalid request body", err)
			return
		}

		// Validate the table struct
		if !CheckForValidTable(&table) {
			response.BadRequest(w, "Invalid table definition", nil)
			return
		}

		// Get the project ID from the URL
		urlVariables := mux.Vars(r)
		projectId := urlVariables["project_id"]
		if projectId == "" {
			response.BadRequest(w, "Project ID is required", nil)
			return
		}
		// Call the service function to create the table
		tableOID, err := CreateTable(r.Context(), projectId, &table, config.DB)
		if err != nil {
			if errors.Is(err, response.ErrUnauthorized) {
				response.UnAuthorized(w, "Unauthorized", nil)
				return
			}
			app.ErrorLog.Println("Table creation failed:", err)
			response.InternalServerError(w, "Failed to create table", err)
			return
		}
		// Return a success response
		response.Created(w, "Table created successfully", map[string]string{
			"oid": tableOID,
		})
	}
}


// UpdateTableHandler godoc
// @Summary Update an existing table
// @Description Update table structure by adding, modifying, or deleting columns
// @Tags tables
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param updates body UpdateTableSchema true "new table schema updates"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 500 {object} response.ErrorResponse500
// @Router /api/projects/{project_id}/tables/{table_id} [put]
func UpdateTableHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		updates := UpdateTableSchema{}
		// Parse the request body to populate the UpdateTable struct
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			response.BadRequest(w, "Invalid request body", err)
			return
		}

		// Get the project ID and Table id from the URL
		urlVariables := mux.Vars(r)
		projectOID := urlVariables["project_id"]
		tableId := urlVariables["table_id"]
		if projectOID == "" || tableId == "" {
			response.BadRequest(w, "Project ID and Table ID are required", nil)
			return
		}

		// Call the service function to update the table
		if err := UpdateTable(r.Context(), projectOID, tableId, &updates, config.DB); err != nil {
			if errors.Is(err, response.ErrUnauthorized) {
				response.UnAuthorized(w, "Unauthorized", nil)
				return
			}
			app.ErrorLog.Println("Table update failed:", err)
			response.InternalServerError(w, "Failed to update table", err)
			return
		}
		// Return a success response
		response.OK(w, "Table updated successfully", nil)
	}
}

// DeleteTableHandler godoc
// @Summary Delete a table
// @Description Delete a table from the specified project
// @Tags tables
// @Produce json
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 500 {object} response.ErrorResponse500
// @Router /api/projects/{project_id}/tables/{table_id} [delete]
func DeleteTableHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOID := urlVariables["project_id"]
		tableOID := urlVariables["table_id"]
		if projectOID == "" || tableOID == "" {
			response.BadRequest(w, "Project ID and Table ID are required", nil)
			return
		}
		// Call the service function to delete the table
		if err := DeleteTable(r.Context(), projectOID, tableOID, config.DB); err != nil {
			if errors.Is(err, response.ErrUnauthorized) {
				response.UnAuthorized(w, "Unauthorized", nil)
				return
			}
			app.ErrorLog.Println("Table deletion failed:", err)
			response.InternalServerError(w, "Failed to delete table", err)
			return
		}
		// Return a success response
		response.OK(w, "Table deleted successfully", nil)
	}
}

// ReadTableHandler godoc
// @Summary Read table data
// @Description Get table structure and data with pagination, filtering and sorting
// @Tags tables
// @Produce json
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param page query int true "Page number"
// @Param limit query int true "Number of records per page"
// @Param order_by query string false "Column to order by"
// @Param order query string false "Sort order (asc or desc)"
// @Param filter query string false "Filter condition (e.g. name=value)"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=Data}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/projects/{project_id}/tables/{table_id} [get]
func ReadTableHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// url variables
		urlVariables := mux.Vars(r)
		projectId := urlVariables["project_id"]
		tableId := urlVariables["table_id"]
		if projectId == "" || tableId == "" {
			response.BadRequest(w, "Project ID and Table ID are required", nil)
			return
		}
		// query parameters
		parameters := r.URL.Query()
		if parameters == nil || parameters["page"] == nil || parameters["limit"] == nil {
			response.BadRequest(w, "Page and Limit are required", nil)
			return
		}

		// Call the service function to read the table
		data, err := ReadTable(r.Context(), projectId, tableId, parameters, config.DB)
		if err != nil {
			if errors.Is(err, response.ErrUnauthorized) {
				response.UnAuthorized(w, "Unauthorized", nil)
				return
			}
			app.ErrorLog.Println("Could not read table:", err)
			response.InternalServerError(w, "Could not read table", err)
			return
		}

		response.OK(w, "Table Read Succesfully", data)
	}
}
