package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

var (
	App *Application
	Mux *http.ServeMux
)

var DB *pgxpool.Pool

func Init(infoLog, errorLog *log.Logger) {

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	App = &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	Mux = http.NewServeMux()

	// ---- database connection ---- //
	dbURL := os.Getenv("DATABASE_URL")

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL successfully! ✅")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed. 🔌")
	}
}
