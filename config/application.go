package config

import (
	"DBHS/caching"
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/axiomhq/axiom-go/axiom"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	Agent "github.com/Database-Hosting-Services/AI-Agent/RAG"
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

// UserDbConfig manages database connection pools
type UserDbConfig struct {
	baseConfig *pgxpool.Config
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

	ConfigManager *UserDbConfig

	AI Agent.RAGmodel

	AxiomLogger *axiom.Client
)

func loadEnv() {
	_, filename, _, _ := runtime.Caller(0) // Gets current file path
	rootDir := filepath.Dir(filepath.Dir(filename))
	envPath := filepath.Join(rootDir, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Error loading .env: ", err)
	}
}

const (
	deploy = true // Set to true if running in production, false for development
)

func Init(infoLog, errorLog *log.Logger) {

	if !deploy {
		loadEnv()
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

	// --------------------------------------- DataBase Pool Manager ----------------------------------------- //

	ConfigManager, err = NewUserDbConfig(adminConnStr)
	if err != nil {
		errorLog.Fatalf("Unable to create UserDbConfig: %v", err)
	}

	infoLog.Println("UserDbConfig created successfully! ✅")

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

	// this clear the cached prepared statement
	//محدش يعدل فيها عشان انا اتبضنت من كتفم دي لغه
	// there is an error occurs when you restart the server :
	// ERROR: prepared statement "stmtcache_d40c25297f5a9db6d92b9594942d1217a18da17e46487cf5" already exists (SQLSTATE 42P05)
	// it means that the prepared statement already exists and you cannot recache it
	// so this function should remove all cached prepared statements when the server starts
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Clear any existing statements
		_, err := conn.Exec(ctx, "DISCARD ALL")
		return err
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

	PgTypes = make(map[uint32]string)
	if err := LoadTypeMap(context.Background(), DB); err != nil {
		errorLog.Fatal(err)
	}

	AI = Agent.GetRAG(&Agent.RAGConfig{
		GeminiAPIKey:         os.Getenv("GEMINI_API_KEY"),
		GeminiModel:          os.Getenv("GEMINI_MODEL"),
		GeminiEmbeddingModel: os.Getenv("GEMINI_EMBEDDING_MODEL"),
		PineconeAPIKey:       os.Getenv("PINECONE_API_KEY"),
		PineconeIndexName:    os.Getenv("PINECONE_INDEX_NAME"),
		PineconeIndexHost:    os.Getenv("PINECONE_INDEX_HOST"),
	})
	infoLog.Println("Connected to AI successfully! ✅")

	AxiomLogger, err = axiom.NewClient(
		axiom.SetPersonalTokenConfig(os.Getenv("AXIOM_TOKEN"), os.Getenv("AXIOM_ORG_ID")), // Set your Axiom personal token and organization ID
	)
	if err != nil {
		errorLog.Fatalf("Failed to create Axiom client: %v", err)
	}
	infoLog.Println("Connected to Axiom successfully! ✅")
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
	// config.CloseAllPools();
}
