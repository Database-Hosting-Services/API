package tables

import (
	"encoding/json"
	"net/http"
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
			ForiegnKey: {
				tableName: "",
				columnName: "",
			},
		}
	]

*/
func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for creating a table
	table := ClientTable{}
	// Parse the request body to populate the table struct
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	

}