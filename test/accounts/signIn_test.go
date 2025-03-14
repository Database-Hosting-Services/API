package accounts_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"DBHS/accounts"
	global "DBHS/config"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
)

func TestSignInUser_Success(t *testing.T) {
	StartUp()
	// defer Cleanup()

	app := &global.Application{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// Prepare user details.
	userOID := "test-oid-128"
	username := "testuser_4"
	email := "test_4@example.com"
	image := ""
	userPassword, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(context.Background(),
		`INSERT INTO users (oid, username, email, password, image, created_at, last_login)
		 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		userOID, username, email, userPassword, image,
	)

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"/api/user/sign-in",
		CreateBody(map[string]interface{}{
			"email":    email,
			"password": "pass",
			"username": username,
			"image":    image,
			"oid":      userOID,
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := accounts.SignIn(app)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "User signed in successfully")
}

func TestSignInUser_Fail(t *testing.T) {
	StartUp()
	// defer Cleanup()

	app := &global.Application{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	userEmail := "test_70@example.com"
	userPassword := "pass"

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		"/api/user/sign-in",
		CreateBody(map[string]interface{}{
			"email":    userEmail,
			"password": userPassword,
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := accounts.SignIn(app)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "InCorrect Email or Password")
}
