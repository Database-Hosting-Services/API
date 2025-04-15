package test

import (
	"DBHS/config"
	"DBHS/indexes"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
	// Get authentication token using the existing user credentials
	token, err := AuthenticateTestUser()
	if err != nil {
		t.Fatalf("Failed to authenticate test user: %v", err)
	}

	// Get project ID from environment
	projectID := os.Getenv("TEST_PROJECT_ID")
	if projectID == "" {
		t.Fatalf("TEST_PROJECT_ID not set in environment")
	}

	// Create index request data
	indexData := indexes.IndexData{
		IndexName: "employees_last_name_test",
		IndexType: "btree",
		Columns:   []string{"first_name", "last_name"},
		TableName: "employees",
	}
	jsonData, err := json.Marshal(indexData)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}

	// Create index request
	req := httptest.NewRequest("POST", "/api/projects/"+projectID+"/indexes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	indexes.DefineURLs()
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code, "Response body: %s", w.Body.String())

	// Cleanup: Get all indexes to find the one we just created and delete it
	t.Cleanup(func() {
		// Delete the index directly using API
		deleteReq := httptest.NewRequest("GET", "/api/projects/"+projectID+"/indexes", nil)
		deleteReq.Header.Set("Authorization", "Bearer "+token)
		getResp := httptest.NewRecorder()
		config.Router.ServeHTTP(getResp, deleteReq)

		var getResponse struct {
			Message string
			Data    []indexes.RetrievedIndex
		}
		if err := json.NewDecoder(getResp.Body).Decode(&getResponse); err != nil {
			t.Logf("Warning: Failed to decode response for cleanup: %v", err)
			return
		}

		// Find and delete the test index
		for _, idx := range getResponse.Data {
			if idx.IndexName == "employees_last_name_test" {
				deleteReq := httptest.NewRequest("DELETE", "/api/projects/"+projectID+"/indexes/"+idx.IndexOid, nil)
				deleteReq.Header.Set("Authorization", "Bearer "+token)
				deleteResp := httptest.NewRecorder()
				config.Router.ServeHTTP(deleteResp, deleteReq)
				break
			}
		}
	})
}

func TestGetProjectIndexes(t *testing.T) {
	// Get authentication token using the existing user credentials
	token, err := AuthenticateTestUser()
	if err != nil {
		t.Fatalf("Failed to authenticate test user: %v", err)
	}

	// Get project ID from environment
	projectID := os.Getenv("TEST_PROJECT_ID")
	if projectID == "" {
		t.Fatalf("TEST_PROJECT_ID not set in environment")
	}

	// Get project indexes
	req := httptest.NewRequest("GET", "/api/projects/"+projectID+"/indexes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	indexes.DefineURLs()
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code, "Response body: %s", w.Body.String())

	var getResponse struct {
		Message string
		Data    []indexes.RetrievedIndex
	}
	err = json.NewDecoder(w.Body).Decode(&getResponse)
	assert.NoError(t, err)

	// There should be at least one index
	assert.NotEmpty(t, getResponse.Data, "No indexes found")
}

func TestUpdateIndexName(t *testing.T) {
	// Get authentication token using the existing user credentials
	token, err := AuthenticateTestUser()
	if err != nil {
		t.Fatalf("Failed to authenticate test user: %v", err)
	}

	// Get project ID from environment
	projectID := os.Getenv("TEST_PROJECT_ID")
	if projectID == "" {
		t.Fatalf("TEST_PROJECT_ID not set in environment")
	}

	// Create an index specifically for this test with unique timestamp
	indexName := "updatetest_" + time.Now().Format("150405")
	indexData := indexes.IndexData{
		IndexName: indexName,
		IndexType: "btree",
		Columns:   []string{"id"},
		TableName: "employees",
	}
	jsonData, err := json.Marshal(indexData)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}

	createReq := httptest.NewRequest("POST", "/api/projects/"+projectID+"/indexes", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	indexes.DefineURLs()
	config.Router.ServeHTTP(w, createReq)

	// Check if creation was successful
	assert.Equal(t, http.StatusCreated, w.Code, "Failed to create index for update test: %s", w.Body.String())

	// Get all indexes to find the one we just created
	getReq := httptest.NewRequest("GET", "/api/projects/"+projectID+"/indexes", nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	config.Router.ServeHTTP(w, getReq)

	var getResponse struct {
		Message string
		Data    []indexes.RetrievedIndex
	}
	err = json.NewDecoder(w.Body).Decode(&getResponse)
	assert.NoError(t, err)

	// Find the index we created
	var indexID string
	for _, idx := range getResponse.Data {
		if idx.IndexName == indexName {
			indexID = idx.IndexOid
			break
		}
	}

	assert.NotEmpty(t, indexID, "Could not find created index")

	// Define new name for the index
	newName := "updated_" + indexName

	// Update index name
	updateData := map[string]string{
		"name": newName,
	}
	jsonData, err = json.Marshal(updateData)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/projects/"+projectID+"/indexes/"+indexID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code, "Response body: %s", w.Body.String())

	// Do a direct cleanup without further API calls
	deleteReq := httptest.NewRequest("DELETE", "/api/projects/"+projectID+"/indexes/"+indexID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteResp := httptest.NewRecorder()
	config.Router.ServeHTTP(deleteResp, deleteReq)

	// Check for appropriate status code
	if deleteResp.Code != http.StatusOK && deleteResp.Code != http.StatusNotFound {
		t.Logf("Warning: Cleanup delete returned status %d: %s", deleteResp.Code, deleteResp.Body.String())
	}
}

func TestDeleteIndex(t *testing.T) {
	// Get authentication token using the existing user credentials
	token, err := AuthenticateTestUser()
	if err != nil {
		t.Fatalf("Failed to authenticate test user: %v", err)
	}

	// Get project ID from environment
	projectID := os.Getenv("TEST_PROJECT_ID")
	if projectID == "" {
		t.Fatalf("TEST_PROJECT_ID not set in environment")
	}

	// Create a test index specifically to delete with unique timestamp
	indexName := "index_to_delete_" + time.Now().Format("150405")
	indexData := indexes.IndexData{
		IndexName: indexName,
		IndexType: "btree",
		Columns:   []string{"id"},
		TableName: "employees",
	}
	jsonData, err := json.Marshal(indexData)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}

	createReq := httptest.NewRequest("POST", "/api/projects/"+projectID+"/indexes", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	indexes.DefineURLs()
	config.Router.ServeHTTP(w, createReq)

	// Check if creation was successful
	assert.Equal(t, http.StatusCreated, w.Code, "Failed to create index for deletion test: %s", w.Body.String())

	// Get all indexes to find the one we just created
	getReq := httptest.NewRequest("GET", "/api/projects/"+projectID+"/indexes", nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	config.Router.ServeHTTP(w, getReq)

	var getResponse struct {
		Message string
		Data    []indexes.RetrievedIndex
	}
	err = json.NewDecoder(w.Body).Decode(&getResponse)
	assert.NoError(t, err)

	// Find the index we just created
	var indexID string
	for _, idx := range getResponse.Data {
		if idx.IndexName == indexName {
			indexID = idx.IndexOid
			break
		}
	}

	assert.NotEmpty(t, indexID, "Could not find the created index")

	// Delete the index
	req := httptest.NewRequest("DELETE", "/api/projects/"+projectID+"/indexes/"+indexID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code, "Response body: %s", w.Body.String())

	var deleteResponse struct {
		Message string
	}
	err = json.NewDecoder(w.Body).Decode(&deleteResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Index deleted successfully", deleteResponse.Message)
}
