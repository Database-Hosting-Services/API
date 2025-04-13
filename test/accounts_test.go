package test

import (
	"DBHS/accounts"
	"DBHS/config"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	// Setup test data
	email, username := setupUserTest(t)
	userData := map[string]string{
		"username": username,
		"email":    email,
		"password": "Test@123456",
	}
	jsonData, err := json.Marshal(userData)
	assert.NoError(t, err)

	// Create request
	req := httptest.NewRequest("POST", "/api/user/sign-up", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	accounts.DefineURLs()
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response APIResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Status)
	assert.Equal(t, "User signed up successfully, check your email for verification", response.Message)

	// Cleanup
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()
}

func TestSignIn(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	_, err := CreateTestUser(config.App, email, username, "Test@123456")
	assert.NoError(t, err)

	// Cleanup test user
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()

	// Sign in data
	signInData := map[string]string{
		"email":    email,
		"password": "Test@123456",
	}
	jsonData, err := json.Marshal(signInData)
	assert.NoError(t, err)

	// Create request
	req := httptest.NewRequest("POST", "/api/user/sign-in", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "User signed in successfully", response.Message)
	assert.NotEmpty(t, response.Data["token"])
	assert.Equal(t, username, response.Data["username"])
	assert.Equal(t, email, response.Data["email"])
}

func TestUpdatePassword(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	token, err := CreateTestUser(config.App, email, username, "Test@123456")
	assert.NoError(t, err)

	// Cleanup test user
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()

	// Update password data
	updateData := map[string]string{
		"current_password": "Test@123456",
		"password":         "NewTest@123456",
		"confirm_password": "NewTest@123456",
	}
	jsonData, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// Create request
	req := httptest.NewRequest("POST", "/api/users/update-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Execute request
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "Password updated successfully", response.Message)
}

func TestUpdateUser(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	token, err := CreateTestUser(config.App, email, username, "Test@123456")
	assert.NoError(t, err)

	// Get user OID from sign-in response
	signInData := map[string]string{
		"email":    email,
		"password": "Test@123456",
	}
	jsonData, err := json.Marshal(signInData)
	assert.NoError(t, err)

	signInReq := httptest.NewRequest("POST", "/api/user/sign-in", bytes.NewBuffer(jsonData))
	signInReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	config.Router.ServeHTTP(w, signInReq)

	var signInResponse APIResponse
	err = json.NewDecoder(w.Body).Decode(&signInResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, signInResponse.Data["oid"])
	userOID := signInResponse.Data["oid"].(string)

	// Update user data
	updateData := map[string]string{
		"username": "updated_" + username,
		"image":    "new_image.jpg",
	}
	jsonData, err = json.Marshal(updateData)
	assert.NoError(t, err)

	// Create update request
	req := httptest.NewRequest("PATCH", "/api/users/"+userOID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()

	// Execute update request
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "User's data updated successfully", response.Message)
	assert.Equal(t, "updated_"+username, response.Data["username"])
	assert.Equal(t, "new_image.jpg", response.Data["image"])

	// Cleanup
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()
}

func TestForgetPassword(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	_, err := CreateTestUser(config.App, email, username, "Test@123456")
	assert.NoError(t, err)

	// Cleanup test user
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()

	// Forget password request data
	requestData := map[string]string{
		"email": email,
	}
	jsonData, err := json.Marshal(requestData)
	assert.NoError(t, err)

	// Create request
	req := httptest.NewRequest("POST", "/api/user/forget-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "Verification Code Sent", response.Message)
}

func TestForgetPasswordVerify(t *testing.T) {
	// Setup test user
	email, username := setupUserTest(t)
	_, err := CreateTestUser(config.App, email, username, "Test@123456")
	assert.NoError(t, err)

	// Cleanup test user
	defer func() {
		err := Drop("users", map[string]interface{}{
			"email": email,
		})
		if err != nil {
			t.Errorf("Failed to cleanup test user: %v", err)
		}
	}()

	// First, trigger forget password to get verification code
	requestData := map[string]string{
		"email": email,
	}
	jsonData, err := json.Marshal(requestData)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/user/forget-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	config.Router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get verification code from cache
	var userData accounts.UserUnVerified
	_, err = config.VerifyCache.Get("forget:"+email, &userData)
	assert.NoError(t, err)

	// Verify and reset password
	verifyData := map[string]string{
		"email":    email,
		"code":     userData.Code,
		"password": "NewTest@123456",
	}
	jsonData, err = json.Marshal(verifyData)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/api/user/forget-password/verify", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	// Execute request
	config.Router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "Password updated successfully", response.Message)
}
