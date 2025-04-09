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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the table struct
		if !CheckForValidTable(&table) {
			http.Error(w, "Invalid table definition", http.StatusBadRequest)
			return
		}

		// Get the project ID from the URL
		urlVariables := mux.Vars(r)
		projectId := urlVariables["project_id"]
		if projectId == "" {
			http.Error(w, "Project ID is required", http.StatusBadRequest)
			return
		}
		// Call the service function to create the table
		if err := CreateTable(r.Context(), projectId, &table, config.DB); err != nil {
			if err.Error() == "Unauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			app.ErrorLog.Println("Table creation failed:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		"columns" : [
		]
	},
	"delete": {
		"columns" : [	
		]
	},
*/

func UpdateTableHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}