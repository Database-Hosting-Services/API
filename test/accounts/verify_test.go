package accounts_test

import (
	"DBHS/accounts"
	"DBHS/caching"
	"DBHS/utils"

	"os"
	"fmt"
	"bcrypt"
	"context"
	"testing"

	// "time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	// "github.com/stretchr/testify/assert"

)

var (
	db 		*pgxpool.Pool
	cache 	*caching.RedisClient
)

func CachUser(username , email, password, OID, code string) {
	user := accounts.UserVerify{
		User: accounts.User{
			OID: OID,
			Username: username,
			Email: email,
			Password: bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost),
		},
		Code: code,
	}
	cac
}

func Cleanup() {
	_, err := db.Exec(context.Background(), "select truncate_all_tables();")
	if err != nil {
		panic(err)
	}
}

func StartUp() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	dbURL := os.Getenv("TEST_DATABASE_URL")
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		panic(err)
	}

	db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}

	cache, err = caching.NewRedisClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASS"), 0)
	if err != nil {
		panic(err)
	}
	
}

func TestBasicVerify(t *testing.T) {
	StartUp()







	Cleanup()
}