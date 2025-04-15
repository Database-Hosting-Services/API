package test

import (
	"DBHS/accounts"
	"DBHS/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// Common response structure
type APIResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func setupUserTest(t *testing.T) (email, username string) {
	t.Helper()
	timestamp := time.Now().Unix()
	email = fmt.Sprintf("test%d@gmail.com", timestamp)
	username = fmt.Sprintf("testuser%d", timestamp)

	t.Cleanup(func() {
		err := Drop("users", map[string]interface{}{
			"email":    email,
			"username": username,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test data: %v", err)
		}
	})

	return email, username
}

// cleanupKeysFromCache removes all cache entries for the given keys
func cleanupKeysFromCache(keys ...string) {
	for _, key := range keys {
		config.VerifyCache.Delete(key)
	}
}

// Drop deletes a specific row from a table based on the provided conditions
/*
	tableName: the name of the table to delete from
	conditions: a map of column names and values to match against
*/
func Drop(tableName string, conditions map[string]interface{}) error {
	// Build the WHERE clause dynamically based on conditions
	whereClause := ""
	values := make([]interface{}, 0)
	i := 1

	for field, value := range conditions {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("%s = $%d", field, i)
		values = append(values, value)
		i++
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, whereClause)

	_, err := config.DB.Exec(context.Background(), query, values...)
	if err != nil {
		return fmt.Errorf("failed to delete from %s: %w", tableName, err)
	}

	return nil
}

// AuthenticateTestUser signs in with the existing user credentials from the .env file
func AuthenticateTestUser() (token string, err error) {
	// Get credentials from .env file
	email := os.Getenv("TEST_USER_EMAIL")
	password := os.Getenv("TEST_USER_PASSWORD")

	if email == "" || password == "" {
		return "", fmt.Errorf("TEST_USER_EMAIL or TEST_USER_PASSWORD not set in environment")
	}

	// Set up a mock SendMail function to avoid sending actual emails during tests
	originalSender := accounts.SetEmailSender(func(d interface{}, from, to, code, subject string) error {
		return nil
	})
	defer accounts.SetEmailSender(originalSender)

	// Create a sign-in request
	user := accounts.UserSignIn{
		Email:    email,
		Password: password,
	}

	// Create the signin request
	jsonData, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Register routes
	accounts.DefineURLs()

	// Create and execute the signin request
	req := httptest.NewRequest("POST", "/api/user/sign-in", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	config.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		return "", fmt.Errorf("sign-in handler returned wrong status code: got %v want %v with body %v",
			rr.Code, http.StatusOK, rr.Body.String())
	}

	// Parse the response to get the token
	var response struct {
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	token, ok := response.Data["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}

// CreateTestUser creates a test user and verifies it without sending emails
func CreateTestUser(app *config.Application, email, username, password string) (token string, err error) {
	// Create test user data
	user := accounts.UserUnVerified{
		User: accounts.User{
			Username: username,
			Email:    email,
			Password: password,
		},
	}

	// Create the signup request
	jsonData, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Create a new router and register the routes
	accounts.DefineURLs()

	// Create and execute the signup request
	req := httptest.NewRequest("POST", "/api/user/sign-up", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	config.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		cleanupKeysFromCache(user.Email, user.Username)
		return "", fmt.Errorf("Signup handler returned wrong status code: got %v want %v with body %v", rr.Code, http.StatusCreated, rr.Body.String())
	}

	// Get the verification code from cache
	var cachedData accounts.UserUnVerified
	_, err = config.VerifyCache.Get(user.Email, &cachedData)
	if err != nil {
		cleanupKeysFromCache(user.Email, user.Username)
		return "", fmt.Errorf("failed to get verification code from cache: %w", err)
	}

	// Now verify the user
	verifyData := map[string]string{
		"email": user.Email,
		"code":  cachedData.Code,
	}

	jsonData, err = json.Marshal(verifyData)
	if err != nil {
		cleanupKeysFromCache(user.Email, user.Username)
		return "", fmt.Errorf("failed to marshal verify data: %w", err)
	}

	req = httptest.NewRequest("POST", "/api/user/verify", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	config.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		cleanupKeysFromCache(user.Email, user.Username)
		return "", fmt.Errorf("verify handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}

	// Parse the verification response to get the token
	var response struct {
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	token, ok := response.Data["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}
