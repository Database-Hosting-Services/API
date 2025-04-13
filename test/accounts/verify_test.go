package accounts_test

import (
	"DBHS/accounts"
	"DBHS/caching"
	global "DBHS/config"
	"DBHS/utils"

	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	db    *pgxpool.Pool
	cache *caching.RedisClient
)

func CacheUser(user *accounts.UserUnVerified) error {
	err := cache.Set(user.Email, user, time.Minute*30)
	if err != nil {
		return err
	}
	err = cache.Set(user.Username, true, time.Minute*30)
	if err != nil {
		return err
	}
	return nil
}

func GenerateUnVerifiedUser() *accounts.UserUnVerified {
	user := &accounts.UserUnVerified{
		User: accounts.User{
			Username: rand.Text()[0:10],
			Password: utils.HashedPassword(rand.Text()[0:10]),
			Email:    rand.Text()[0:10] + "@gmail.com",
		},
		Code: utils.GenerateVerficationCode(),
	}
	return user
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
	global.DB = db
	global.VerifyCache = cache

	global.Env = &global.Environment{
		AccessTokenExpiryHour:  global.ACCESS_TOKEN_EXPIRY_HOUR,
		VerifyCodeExpiryMinute: global.VERIFY_CODE_EXPIRY_MINUTE,
	}
}

func CreateBody(payload map[string]interface{}) *bytes.Buffer {
	JsonData, _ := json.Marshal(payload)
	return bytes.NewBuffer(JsonData)
}

func TestBasicVerifySuccess(t *testing.T) {
	StartUp()
	defer Cleanup()
	app := &global.Application{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	user := GenerateUnVerifiedUser()
	err := CacheUser(user)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8000/api/user/verify", CreateBody(map[string]interface{}{"email": user.Email, "code": user.Code}))
	if err != nil {
		assert.Fail(t, err.Error())
	}

	res := httptest.NewRecorder()
	handler := accounts.Verify(app)
	handler(res, req)
	assert.Equal(t, res.Code, http.StatusCreated)
}

func TestBasicVerifyFail(t *testing.T) {
	StartUp()
	defer Cleanup()
	app := &global.Application{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	user := GenerateUnVerifiedUser()
	err := CacheUser(user)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8000/api/user/verify", CreateBody(map[string]interface{}{"email": user.Email, "code": utils.GenerateVerficationCode()}))
	if err != nil {
		assert.Fail(t, err.Error())
	}

	res := httptest.NewRecorder()
	handler := accounts.Verify(app)
	handler(res, req)
	assert.NotEqual(t, res.Code, http.StatusCreated)
}

func TestCheckCommit(t *testing.T) {
	StartUp()
	defer Cleanup()
	app := &global.Application{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	user := GenerateUnVerifiedUser()
	err := CacheUser(user)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8000/api/user/verify", CreateBody(map[string]interface{}{"email": user.Email, "code": user.Code}))
	if err != nil {
		assert.Fail(t, err.Error())
	}

	res := httptest.NewRecorder()
	handler := accounts.Verify(app)
	handler(res, req)
	assert.Equal(t, res.Code, http.StatusCreated)
	var userv2 accounts.User
	err = accounts.GetUser(ctx, db, user.Email, accounts.SELECT_USER_BY_Email, []interface{}{
		&userv2.ID,
		&userv2.OID,
		&userv2.Username,
		&userv2.Email,
		&userv2.Password,
		&userv2.Image,
		&userv2.CreatedAt,
		&userv2.LastLogin,
	}...)
	assert.Nil(t, err)
	assert.Equal(t, user.OID, userv2.OID)
	assert.Equal(t, user.Email, userv2.Email)
	assert.Equal(t, user.Username, userv2.Username)
}
