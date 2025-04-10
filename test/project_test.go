package test

import (
	"DBHS/config"
	"log"
	"os"
	"testing"
)

var app *config.Application

func TestMain(m *testing.M) {
	// Setup
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize the application
	config.Init(infoLog, errorLog)
	app = config.App

	// Run tests
	code := m.Run()

	// Cleanup
	config.CloseDB()

	os.Exit(code)
}

func TestCreateTestUser(t *testing.T) {
	email, username := setupUserTest(t)

	err := CreateTestUser(app, email, username, "Test@123456")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
}
