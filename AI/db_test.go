package ai_test

import (
	ai "DBHS/AI"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func loadEnv() {
	_, filename, _, _ := runtime.Caller(0) // Gets current file path
	rootDir := filepath.Dir(filepath.Dir(filename))
	envPath := filepath.Join(rootDir, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Error loading .env: ", err)
	}
}

func TestExtractDatabaseSchema(t *testing.T) {
	// load env
	loadEnv()

	// connect to db using pgxpool
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("Failed to parse database config: %v", err)
	}
	dbConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	// extract schema
	schema, err := ai.ExtractDatabaseSchema(context.Background(), db)
	if err != nil {
		t.Fatalf("Failed to extract database schema: %v", err)
	}

	// print schema
	fmt.Println(schema)
}
