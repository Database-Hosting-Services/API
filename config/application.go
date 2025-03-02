package config

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

type Environment struct {
	AccessTokenExpiryHour 	int
	AccessTokenSecret 		[]byte
}

var (
	App		*Application
	Mux 	*http.ServeMux
	DB  	*pgxpool.Pool
	Env 	*Environment	
)

func Init(infoLog, errorLog *log.Logger) {

	if err := godotenv.Load("../.env"); err != nil {
		errorLog.Fatal("Error loading .env file")
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
		errorLog.Fatalf("Unable to parse database URL: %v", err)
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		errorLog.Fatalf("Unable to connect to database: %v", err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		errorLog.Fatalf("Unable to ping database: %v", err)
	}

	infoLog.Println("Connected to PostgreSQL successfully! ✅")

	AccessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	AccessTokenExpiryHour, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_HOUR"))
	if err != nil {
		errorLog.Fatal(err)
	}

	Env = &Environment{
		AccessTokenExpiryHour: 	AccessTokenExpiryHour,
		AccessTokenSecret: 		[]byte(AccessTokenSecret),
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		App.InfoLog.Println("Database connection closed. 🔌")
	}
}
