package config

import (
	"DBHS/caching"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

type Environment struct {
	AccessTokenExpiryHour  int
	AccessTokenSecret      []byte
	VerifyCodeExpiryMinute int
}

type DatabaseConfig struct {
	AdminConnStr string
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string // connection string controls how SSL/TLS encryption is used when connecting to the database
}

var (
	App         *Application
	Mux         *http.ServeMux
	Router      *mux.Router
	DB          *pgxpool.Pool // Regular database connection
	AdminDB     *pgxpool.Pool // Admin database connection
	VerifyCache *caching.RedisClient
	EmailSender *gomail.Dialer
	Env         *Environment
	DBConfig    *DatabaseConfig
)

func Init(infoLog, errorLog *log.Logger) {

	if err := godotenv.Load("../.env"); err != nil {
		errorLog.Fatal("Error loading .env file --> %s", err.Error())
	}

	App = &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	Mux = http.NewServeMux()
	Router = mux.NewRouter()

	// --------------------------------------- admin database connection ----------------------------------------- //

	// Get the admin connection string
	adminConnStr := os.Getenv("DATABASE_ADMIN_URL")
	if adminConnStr == "" {
		errorLog.Fatal("DATABASE_ADMIN_URL is not set")
	}

	// Parse the admin database URL to extract components
	DBConfig = ParseDatabaseURL(adminConnStr)
	DBConfig.AdminConnStr = adminConnStr

	adminConfig, err := pgxpool.ParseConfig(DBConfig.AdminConnStr)
	if err != nil {
		errorLog.Fatalf("Unable to parse admin database URL: %v", err)
	}

	AdminDB, err = pgxpool.NewWithConfig(context.Background(), adminConfig)
	if err != nil {
		errorLog.Fatalf("Unable to connect to admin database: %v", err)
	}

	if err := AdminDB.Ping(context.Background()); err != nil {
		errorLog.Fatalf("Unable to ping admin database: %v", err)
	}

	infoLog.Println("Connected to Admin PostgreSQL successfully! ✅")

	// --------------------------------------- database connection ----------------------------------------- //
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		errorLog.Fatal("DATABASE_URL is not set")
	}

	// Connect to the application database
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

	// --------------------------------------- redis connection ----------------------------------------- //
	VerifyCache, err = caching.NewRedisClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASS"), 0)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Connected to Redis successfully! ✅")

	AccessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	Env = &Environment{
		AccessTokenExpiryHour:  ACCESS_TOKEN_EXPIRY_HOUR,
		AccessTokenSecret:      []byte(AccessTokenSecret),
		VerifyCodeExpiryMinute: VERIFY_CODE_EXPIRY_MINUTE,
	}

	EmailSender = gomail.NewDialer("smtp.gmail.com", 587, "thunderdbhostingserver@gmail.com", os.Getenv("GMAIL_PASS"))
	_, err = EmailSender.Dial()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		App.InfoLog.Println("Application database connection closed. 🔌")
	}

	if AdminDB != nil {
		AdminDB.Close()
		App.InfoLog.Println("Admin database connection closed. 🔌")
	}
}
