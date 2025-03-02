package token_test

import (
	"DBHS/config"
	"DBHS/utils/token"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Mock User implementation for testing
type MockUser struct {
	OId      string
	Username string
}

func (m MockUser) GetOId() string {
	return m.OId
}

func (m MockUser) GetUserName() string {
	return m.Username
}

func TestCreateAccessToken_ValidToken(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	expiry := 24

	// Create token
	tokenString, err := token.CreateAccessToken(user, expiry)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	
	// Verify token can be parsed
	parsedToken, parseErr := token.ParseToken(tokenString)
	assert.NoError(t, parseErr)
	assert.True(t, parsedToken.Valid)
	
	// Verify claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, user.GetOId(), claims["id"])
	assert.Equal(t, user.GetUserName(), claims["username"])
}

func TestIsAuthorized_ValidToken(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	
	tokenString, _ := token.CreateAccessToken(user, 24)
	err := token.IsAuthorized(tokenString)
	
	assert.NoError(t, err)
}

func TestIsAuthorized_ExpiredToken(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	
	claims := jwt.MapClaims{
		"id":       user.GetOId(),
		"username": user.GetUserName(),
		"exp":      time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := jwtToken.SignedString(config.Env.AccessTokenSecret)
	
	err := token.IsAuthorized(tokenString)
	
	assert.Error(t, err)
}

func TestIsAuthorized_InvalidSigningMethod(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	
	claims := jwt.MapClaims{
		"id":       user.GetOId(),
		"username": user.GetUserName(),
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	
	err := token.IsAuthorized(tokenString)
	
	assert.Error(t, err)
}

func TestIsAuthorized_InvalidTokenFormat(t *testing.T) {
	tokenString := "invalid.token.format"
	err := token.IsAuthorized(tokenString)
	
	assert.Error(t, err)
}

func TestIsAuthorized_EmptyToken(t *testing.T) {
	tokenString := ""
	err := token.IsAuthorized(tokenString)
	
	assert.Error(t, err)
}

func TestParseToken_ValidToken(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	
	tokenString, _ := token.CreateAccessToken(user, 24)
	parsedToken, err := token.ParseToken(tokenString)
	
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)
	
	// Verify claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, user.GetOId(), claims["id"])
	assert.Equal(t, user.GetUserName(), claims["username"])
}

func TestParseToken_ExpiredToken(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	
	claims := jwt.MapClaims{
		"id":       user.GetOId(),
		"username": user.GetUserName(),
		"exp":      time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := jwtToken.SignedString(config.Env.AccessTokenSecret)
	
	parsedToken, err := token.ParseToken(tokenString)
	
	assert.Error(t, err)
	assert.Nil(t, parsedToken)
}

func TestParseToken_InvalidTokenFormat(t *testing.T) {
	tokenString := "invalid.token.format"
	parsedToken, err := token.ParseToken(tokenString)
	
	assert.Error(t, err)
	assert.Nil(t, parsedToken)
}

func TestGetIdFromToken_ValidToken(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	
	tokenString, _ := token.CreateAccessToken(user, 24)
	id, err := token.GetIdFromToken(tokenString)
	
	assert.NoError(t, err)
	assert.Equal(t, user.GetOId(), id)
}

func TestGetIdFromToken_ExpiredToken(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	
	claims := jwt.MapClaims{
		"id":       user.GetOId(),
		"username": user.GetUserName(),
		"exp":      time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := jwtToken.SignedString(config.Env.AccessTokenSecret)
	
	id, err := token.GetIdFromToken(tokenString)
	
	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestGetIdFromToken_NoIdClaim(t *testing.T) {
	// Setup
	config.Env = &config.Environment{
		AccessTokenSecret: []byte("test-secret"),
	}
	user := MockUser{
		OId:      "user123",
		Username: "testuser",
	}
	
	claims := jwt.MapClaims{
		"username": user.GetUserName(),
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := jwtToken.SignedString(config.Env.AccessTokenSecret)
	
	id, err := token.GetIdFromToken(tokenString)
	
	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestGetIdFromToken_InvalidTokenFormat(t *testing.T) {
	tokenString := "invalid.token.format"
	id, err := token.GetIdFromToken(tokenString)
	
	assert.Error(t, err)
	assert.Empty(t, id)
}