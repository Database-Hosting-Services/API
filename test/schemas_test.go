package test

import (
	"DBHS/config"
	"DBHS/schemas"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseSchema(t *testing.T) {
	// Setup first user (owner)
	ownerEmail, ownerUsername := setupUserTest(t)
	ownerToken, err := CreateTestUser(config.App, ownerEmail, ownerUsername, "Test@123456")
	assert.NoError(t, err)

	// Setup second user (unauthorized)
	unauthorizedEmail, unauthorizedUsername := setupUserTest(t)
	unauthorizedToken, err := CreateTestUser(config.App, unauthorizedEmail, unauthorizedUsername, "Test@123456")
	assert.NoError(t, err)

	// Create a project for the owner
	projectData := map[string]string{
		"name":        "test_project",
		"description": "Test project description",
	}
	jsonData, err := json.Marshal(projectData)
	assert.NoError(t, err)

	createReq := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+ownerToken)
	w := httptest.NewRecorder()
	config.Router.ServeHTTP(w, createReq)

	var projectResponse APIResponse
	err = json.NewDecoder(w.Body).Decode(&projectResponse)
	assert.NoError(t, err)
	projectOID := projectResponse.Data["oid"].(string)

	// Cleanup
	defer func() {
		// Clean up project
		err := Drop("projects", map[string]interface{}{
			"oid": projectOID,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test project: %v", err)
		}

		// Clean up users
		err = Drop("users", map[string]interface{}{
			"email": ownerEmail,
		})
		if err != nil {
			t.Errorf("Failed to cleanup owner user: %v", err)
		}

		err = Drop("users", map[string]interface{}{
			"email": unauthorizedEmail,
		})
		if err != nil {
			t.Errorf("Failed to cleanup unauthorized user: %v", err)
		}
	}()

	t.Run("Success - Owner gets project schema", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects/"+projectOID+"/schema/tables", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		w := httptest.NewRecorder()

		schemas.DefineURLs()
		config.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response APIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Status)
		assert.NotNil(t, response.Data["tables"])
	})

	t.Run("Fail - Unauthorized user tries to get project schema", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects/"+projectOID+"/schema/tables", nil)
		req.Header.Set("Authorization", "Bearer "+unauthorizedToken)
		w := httptest.NewRecorder()

		config.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response APIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, response.Status)
	})

	t.Run("Fail - Invalid project ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects/invalid-project-id/schema/tables", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		w := httptest.NewRecorder()

		config.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response APIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Status)
	})
}

func TestGetDatabaseTableSchema(t *testing.T) {
	// Setup first user (owner)
	ownerEmail, ownerUsername := setupUserTest(t)
	ownerToken, err := CreateTestUser(config.App, ownerEmail, ownerUsername, "Test@123456")
	assert.NoError(t, err)

	// Setup second user (unauthorized)
	unauthorizedEmail, unauthorizedUsername := setupUserTest(t)
	unauthorizedToken, err := CreateTestUser(config.App, unauthorizedEmail, unauthorizedUsername, "Test@123456")
	assert.NoError(t, err)

	// Create a project for the owner
	projectData := map[string]string{
		"name":        "test_project",
		"description": "Test project description",
	}
	jsonData, err := json.Marshal(projectData)
	assert.NoError(t, err)

	createReq := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+ownerToken)
	w := httptest.NewRecorder()
	config.Router.ServeHTTP(w, createReq)

	var projectResponse APIResponse
	err = json.NewDecoder(w.Body).Decode(&projectResponse)
	assert.NoError(t, err)
	projectOID := projectResponse.Data["oid"].(string)

	// Cleanup
	defer func() {
		// Clean up project
		err := Drop("projects", map[string]interface{}{
			"oid": projectOID,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test project: %v", err)
		}

		// Clean up users
		err = Drop("users", map[string]interface{}{
			"email": ownerEmail,
		})
		if err != nil {
			t.Errorf("Failed to cleanup owner user: %v", err)
		}

		err = Drop("users", map[string]interface{}{
			"email": unauthorizedEmail,
		})
		if err != nil {
			t.Errorf("Failed to cleanup unauthorized user: %v", err)
		}
	}()

	t.Run("Success - Owner gets table schema", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects/"+projectOID+"/schema/tables/users", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		w := httptest.NewRecorder()

		schemas.DefineURLs()
		config.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response APIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Status)
		assert.NotNil(t, response.Data["table_schema"])
	})

	t.Run("Fail - Unauthorized user tries to get table schema", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects/"+projectOID+"/schema/tables/users", nil)
		req.Header.Set("Authorization", "Bearer "+unauthorizedToken)
		w := httptest.NewRecorder()

		config.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response APIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, response.Status)
	})

	t.Run("Fail - Table doesn't exist", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects/"+projectOID+"/schema/tables/nonexistent_table", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		w := httptest.NewRecorder()

		config.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response APIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Status)
	})

	t.Run("Fail - Invalid project ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects/invalid-project-id/schema/tables/users", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		w := httptest.NewRecorder()

		config.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response APIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Status)
	})
}
