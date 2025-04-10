package schemas

import (
	"DBHS/config"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseSchema(t *testing.T) {
	// Setup request with dummy data
	app := &config.Application{
		ErrorLog: nil,
		InfoLog:  nil,
	}
	req := httptest.NewRequest("GET", "/api/projects/123/schema/tables", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user-id", 1))
	req = mux.SetURLVars(req, map[string]string{"project_id": "123"})
	w := httptest.NewRecorder()

	// Execute
	handler := GetDatabaseSchema(app)
	handler(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []TableSchema
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify response structure
	for _, table := range response {
		assert.NotEmpty(t, table.TableName)
		assert.NotNil(t, table.Cols)

		for _, col := range table.Cols {
			assert.NotEmpty(t, col.Name)
			assert.NotEmpty(t, col.Type)
		}
	}
}
