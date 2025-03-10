package config

import (
	"DBHS/caching"
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
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

var (
	App         *Application
	Mux         *http.ServeMux
	Router      *mux.Router
	DB          *pgxpool.Pool
	VerifyCache *caching.RedisClient
	EmailSender *gomail.Dialer
	Env         *Environment
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

	// ---- redis connection ---- //
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
		App.InfoLog.Println("Database connection closed. 🔌")
	}
}
