package tables

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

/*
	body
	tableName: "",
	cols: [
		{
			name: "",
			type: "",
			isUnique: "",
			isNullable: "",
			isPrimaryKey: "",
			foriegnKey: {
				tableName: "",
				columnName: "",
			},
		}
	]

*/
func CreateTableHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for creating a table
		table := ClientTable{}
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
		if err := CreateTable(r.Context(), projectId, &table, config.DB); err != nil {
			if err.Error() == "Unauthorized" {
				response.UnAuthorized(w, "Unauthorized", nil)
				return
			}
			app.ErrorLog.Println("Table creation failed:", err)
			response.InternalServerError(w, "Failed to create table", err)
			return
		}
		// Return a success response
		response.Created(w, "Table created successfully", nil)
	}
}

/*
	"insert": {
		"columns" : [

		]
	},
	"update": {
		"oldName": "oldName",
		"columns": [
			// Only include the changed parts
		]
	},
	"delete": [
		"columnName1",
		"columnName2"
	]
*/

func UpdateTableHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		updates := TableUpdate{}
		// Parse the request body to populate the UpdateTable struct
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			response.BadRequest(w, "Invalid request body", err)
			return
		}

		// Get the project ID and Table id from the URL
		urlVariables := mux.Vars(r)
		projectId := urlVariables["project_id"]
		tableId := urlVariables["table_id"]
		if projectId == "" || tableId == "" {
			response.BadRequest(w, "Project ID and Table ID are required", nil)
			return
		}

		// Call the service function to update the table
		if err := UpdateTable(r.Context(), projectId, tableId, &updates, config.DB); err != nil {
			if err.Error() == "Unauthorized" {
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