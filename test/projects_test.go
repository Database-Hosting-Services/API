package test

import (
	"DBHS/config"
	"DBHS/projects"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	token, err := CreateTestUser(config.App, email, username, "Test@123456")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	// Cleanup test user
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()

	// Create project request data
	projectData := map[string]string{
		"name":        "test_project",
		"description": "Test project description",
	}
	jsonData, err := json.Marshal(projectData)
	if err != nil {
		t.Fatalf("Failed to marshal project data: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Register routes and execute request
	projects.DefineURLs()
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response to get project OID
	var response struct {
		Message string
		Data    projects.SafeProjectData
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Data.Oid)

	// Cleanup test project
	defer func() {
		err := Drop("projects", map[string]interface{}{
			"oid": response.Data.Oid,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test project: %v", err)
		}
	}()
}

func TestGetUserProjects(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	token, err := CreateTestUser(config.App, email, username, "Test@123456")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	// Cleanup test user
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()

	// Create a test project first
	projectData := map[string]string{
		"name":        "test_project",
		"description": "Test project description",
	}
	jsonData, err := json.Marshal(projectData)
	if err != nil {
		t.Fatalf("Failed to marshal project data: %v", err)
	}

	// Create project
	createReq := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	projects.DefineURLs()
	config.Router.ServeHTTP(w, createReq)

	var createResponse struct {
		Message string
		Data    projects.SafeProjectData
	}
	err = json.NewDecoder(w.Body).Decode(&createResponse)
	assert.NoError(t, err)

	// Cleanup test project
	defer func() {
		err := Drop("projects", map[string]interface{}{
			"oid": createResponse.Data.Oid,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test project: %v", err)
		}
	}()

	// Test getting projects
	req := httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Message string
		Data    []projects.SafeProjectData
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, createResponse.Data.Oid, response.Data[0].Oid)
}

func TestUpdateProject(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	token, err := CreateTestUser(config.App, email, username, "Test@123456")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	// Cleanup test user
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()

	// Create a test project first
	projectData := map[string]string{
		"name":        "test_project",
		"description": "Test project description",
	}
	jsonData, err := json.Marshal(projectData)
	if err != nil {
		t.Fatalf("Failed to marshal project data: %v", err)
	}

	// Create project
	createReq := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	projects.DefineURLs()
	config.Router.ServeHTTP(w, createReq)

	var createResponse struct {
		Message string
		Data    projects.SafeProjectData
	}
	err = json.NewDecoder(w.Body).Decode(&createResponse)
	assert.NoError(t, err)

	// Cleanup
	defer func() {
		err := Drop("projects", map[string]interface{}{
			"oid": createResponse.Data.Oid,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test project: %v", err)
		}
	}()

	// Update project data
	updateData := map[string]string{
		"name":        "updated_project",
		"description": "Updated project description",
	}
	jsonData, err = json.Marshal(updateData)
	assert.NoError(t, err)

	// Create update request
	req := httptest.NewRequest("PATCH", "/api/projects/"+createResponse.Data.Oid, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse struct {
		Message string
		Data    map[string]interface{}
	}
	err = json.NewDecoder(w.Body).Decode(&updateResponse)
	assert.NoError(t, err)
	assert.Equal(t, "updated_project", updateResponse.Data["name"])
}

func TestDeleteProject(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	token, err := CreateTestUser(config.App, email, username, "Test@123456")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	// Cleanup test user
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()

	// Create a test project first
	projectData := map[string]string{
		"name":        "test_project",
		"description": "Test project description",
	}
	jsonData, err := json.Marshal(projectData)
	if err != nil {
		t.Fatalf("Failed to marshal project data: %v", err)
	}

	// Create project
	createReq := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	projects.DefineURLs()
	config.Router.ServeHTTP(w, createReq)

	var createResponse struct {
		Message string
		Data    projects.SafeProjectData
	}
	err = json.NewDecoder(w.Body).Decode(&createResponse)
	assert.NoError(t, err)

	// Create delete request
	req := httptest.NewRequest("DELETE", "/api/projects/"+createResponse.Data.Oid, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify project is deleted
	verifyReq := httptest.NewRequest("GET", "/api/projects/"+createResponse.Data.Oid, nil)
	verifyReq.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	config.Router.ServeHTTP(w, verifyReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
