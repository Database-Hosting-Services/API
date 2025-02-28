package token_test

import (
	"DBHS/config"
	"DBHS/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenLifecycle(t *testing.T) {
	// Save original key to restore after test
	originalKey := config.Secret_Key
	defer func() { config.Secret_Key = originalKey }()
	
	// Set a test secret key
	config.Secret_Key = "test-secret-key-for-jwt-validation"

	// Create a new token
	token := utils.NewToken()

	// Add standard claims
	token.AddClaim("sub", "1234567890")
	token.AddClaim("name", "Test User")
	token.AddClaim("iat", time.Now().Unix())
	token.AddClaim("exp", time.Now().Add(time.Hour).Unix())

	// Generate token string
	tokenString, err := token.String()
	assert.NoError(t, err, "Token string generation should not error")
	assert.NotEmpty(t, tokenString, "Token string should not be empty")

	// Parse the token string back to a token
	parsedToken, err := utils.NewTokenString(tokenString)
	assert.NoError(t, err, "Parsing valid token should not error")
	assert.True(t, parsedToken.IsValid(), "Parsed token should be valid")

	// Verify claims were preserved
	sub, ok := parsedToken.GetClaim("sub")
	assert.True(t, ok, "Subject claim should exist")
	assert.Equal(t, "1234567890", sub)

	name, ok := parsedToken.GetClaim("name")
	assert.True(t, ok, "Name claim should exist")
	assert.Equal(t, "Test User", name)
}

func TestInvalidToken(t *testing.T) {
	// Save original key to restore after test
	originalKey := config.Secret_Key
	defer func() { config.Secret_Key = originalKey }()
	
	// Set a test secret key
	config.Secret_Key = "test-secret-key-for-jwt-validation"

	// Test with invalid token string
	_, err := utils.NewTokenString("invalid.token.string")
	assert.Error(t, err, "Invalid token should return error")

	// Test with valid format but invalid signature
	_, err = utils.NewTokenString("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.invalid_signature")
	assert.Error(t, err, "Token with invalid signature should return error")

	// Test with empty secret key
	config.Secret_Key = ""
	token := utils.NewToken()
	token.AddClaim("sub", "1234567890")
	_, err = token.String()
	assert.Error(t, err, "Token generation with empty secret key should error")
}

func TestExpiredToken(t *testing.T) {
	// Save original key to restore after test
	originalKey := config.Secret_Key
	defer func() { config.Secret_Key = originalKey }()
	
	// Set a test secret key
	config.Secret_Key = "test-secret-key-for-jwt-validation"

	// Create a token that's already expired
	token := utils.NewToken()
	token.AddClaim("exp", time.Now().Add(-time.Hour).Unix()) // Expired an hour ago

	// Generate token string
	tokenString, err := token.String()
	assert.NoError(t, err, "Token string generation should not error")

	// Parsing an expired token should fail
	_, err = utils.NewTokenString(tokenString)
	assert.Error(t, err, "Parsing expired token should return error")
	assert.Contains(t, err.Error(), "token is expired", "Error should mention expiration")
}

func TestHeaderManipulation(t *testing.T) {
	token := utils.NewToken()

	// Test adding a single header
	token.AddHeader("custom", "value")
	value, ok := token.GetHeader("custom")
	assert.True(t, ok, "Custom header should exist")
	assert.Equal(t, "value", value)

	// Test overriding existing header
	token.AddHeader("custom", "new-value")
	value, ok = token.GetHeader("custom")
	assert.True(t, ok, "Custom header should still exist")
	assert.Equal(t, "new-value", value)

	// Test adding multiple headers
	headers := map[string]interface{}{
		"header1": "value1",
		"header2": "value2",
	}
	token.AddHeaders(headers)
	
	value1, ok := token.GetHeader("header1")
	assert.True(t, ok, "Header1 should exist")
	assert.Equal(t, "value1", value1)
	
	value2, ok := token.GetHeader("header2")
	assert.True(t, ok, "Header2 should exist")
	assert.Equal(t, "value2", value2)

	// Test getting non-existent header
	_, ok = token.GetHeader("nonexistent")
	assert.False(t, ok, "Nonexistent header should not exist")
}

func TestClaimManipulation(t *testing.T) {
	token := utils.NewToken()

	// Test adding a single claim
	token.AddClaim("custom", "value")
	value, ok := token.GetClaim("custom")
	assert.True(t, ok, "Custom claim should exist")
	assert.Equal(t, "value", value)

	// Test overriding existing claim
	token.AddClaim("custom", "new-value")
	value, ok = token.GetClaim("custom")
	assert.True(t, ok, "Custom claim should still exist")
	assert.Equal(t, "new-value", value)

	// Test adding multiple claims
	claims := map[string]interface{}{
		"claim1": "value1",
		"claim2": 42,
		"claim3": true,
	}
	token.AddClaims(claims)
	
	value1, ok := token.GetClaim("claim1")
	assert.True(t, ok, "Claim1 should exist")
	assert.Equal(t, "value1", value1)
	
	value2, ok := token.GetClaim("claim2")
	assert.True(t, ok, "Claim2 should exist")
	assert.Equal(t, 42, value2) // JSON numbers are unmarshaled as float64
	
	value3, ok := token.GetClaim("claim3")
	assert.True(t, ok, "Claim3 should exist")
	assert.Equal(t, true, value3)

	// Test getting non-existent claim
	_, ok = token.GetClaim("nonexistent")
	assert.False(t, ok, "Nonexistent claim should not exist")
}

func TestStringTest(t *testing.T) {
	// Create a token
	token := utils.NewToken()
	token.AddClaim("sub", "1234567890")
	
	// Test with custom secret key
	config.Secret_Key = "custom-secret-key-for-testing-string"
	tokenString, err := token.String()
	assert.NoError(t, err, "should not have error")
	// Verify the token was signed with the custom key
	// To do this properly, we need to manually parse and verify
	// Since the utils.NewTokenString uses config.Secret_Key, we need a different approach
	
	// Save original key to restore after test
	originalKey := config.Secret_Key
	defer func() { config.Secret_Key = originalKey }()
	
	// Set config key to match our custom key
	config.Secret_Key = "custom-secret-key-for-testing-string"
	
	// Now parsing should work
	parsedToken, err := utils.NewTokenString(tokenString)
	assert.NoError(t, err, "Parsing token signed with custom key should work when config key matches")
	
	// Verify claim
	sub, ok := parsedToken.GetClaim("sub")
	assert.True(t, ok, "Subject claim should exist")
	assert.Equal(t, "1234567890", sub)
	config.Secret_Key = ""
	// Test with empty key
	_, err = token.String()
	assert.Error(t, err, "StringTest should error with empty key")
}