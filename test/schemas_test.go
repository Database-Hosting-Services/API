package test

import (
	"DBHS/config"
	"DBHS/schemas"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseSchema(t *testing.T) {
	email, username := setupUserTest(t)

	token, err := CreateTestUser(config.App, email, username, "Test@123456")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	defer Drop("users", map[string]interface{}{
		"email":    email,
		"username": username,
	})

	// Create request
	req := httptest.NewRequest("GET", "/api/projects/test-project-123/schema/tables", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Add project ID to request context using mux vars
	vars := map[string]string{
		"project_id": "test-project-123",
	}
	req = mux.SetURLVars(req, vars)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create handler and serve
	handler := schemas.GetDatabaseSchema(config.App)
	handler.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response []schemas.TableSchema
	json.NewDecoder(w.Body).Decode(&response)
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

func TestGetDatabaseTableSchema(t *testing.T) {
	email, username := setupUserTest(t)

	token, err := CreateTestUser(config.App, email, username, "Test@123456")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	defer Drop("users", map[string]interface{}{
		"email":    email,
		"username": username,
	})

	// Create request
	req := httptest.NewRequest("GET", "/api/projects/test-project-123/schema/tables/users", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	// Add project ID and table ID to request context using mux vars
	vars := map[string]string{
		"project_id": "test-project-123",
		"table-id":   "users",
	}
	req = mux.SetURLVars(req, vars)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create handler and serve
	handler := schemas.GetDatabaseTableSchema(config.App)
	handler.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var tableSchema schemas.TableSchema
	json.NewDecoder(w.Body).Decode(&tableSchema)
	assert.NoError(t, err)

	// Verify table schema structure
	assert.Equal(t, "users", tableSchema.TableName)
	assert.NotEmpty(t, tableSchema.Cols)

	// Verify column structure
	for _, col := range tableSchema.Cols {
		assert.NotEmpty(t, col.Name)
		assert.NotEmpty(t, col.Type)
		// Foreign key might be empty, so we don't assert on it
	}
}
