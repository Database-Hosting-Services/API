package accounts

import (
	"DBHS/caching"
	"DBHS/config"
	"DBHS/utils"
	"DBHS/utils/token"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignupUser(ctx context.Context, db *pgxpool.Pool, user *User) error {
	/*
		store user's data in cache and
		generate a verification code and send it to the user
	*/
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.OID = utils.GenerateOID()
	user.Password = string(hashedPassword)
	user.Verified = false

	// store user's data in cache
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	config.VerifyCache.Set(user.Email, userJson, time.Minute*30)
	config.VerifyCache.Set(user.Username, true, time.Minute*30)

	// send the verification code
	userVerficationCode := utils.GenerateVerficationCode()
	if err = SendMail(config.EmailSender, os.Getenv("GMAIL"), user.Email, userVerficationCode, "Verification Code"); err != nil {
		return err
	}

	return nil
}

func SignInUser(ctx context.Context, db *pgxpool.Pool, cache *caching.RedisClient, user *UserSignIn) (map[string]interface{}, error) {
	exits, err := cache.Exists(user.Email)
	if err != nil {
		return nil, errors.New("InCorrect email or password")
	}

	if exits {
		return SendUserVerificationCode(cache, user.Email, user.Password)
	}

	var authenticatedUser User
	err = GetUser(ctx, db, user.Email, SELECT_USER_BY_Email, []interface{}{
		&authenticatedUser.ID,
		&authenticatedUser.OID,
		&authenticatedUser.Username,
		&authenticatedUser.Email,
		&authenticatedUser.Password,
		&authenticatedUser.Image,
		&authenticatedUser.CreatedAt,
		&authenticatedUser.LastLogin,
	}...)

	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(user.Password, authenticatedUser.Password) {
		return nil, errors.New("InCorrect Email or Password")
	}

	UserTokenData := &User{
		OID:      authenticatedUser.OID,
		Username: authenticatedUser.Username,
	}

	token, err := token.CreateAccessToken(UserTokenData, config.Env.AccessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"oid":      authenticatedUser.OID,
		"username": authenticatedUser.Username,
		"email":    authenticatedUser.Email,
		"image":    authenticatedUser.Image,
		"token":    token,
	}

	return resp, nil
}

func SendUserVerificationCode(cache *caching.RedisClient, email, Password string) (map[string]interface{}, error) {
	userData, err := cache.Get(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	var user UserVerify
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, err
	}

	if !CheckPasswordHash(Password, user.Password) {
		return nil, errors.New("InCorrect Email or Password")
	}

	SendMail(config.EmailSender, os.Getenv("GMAIL"), user.Email, user.Code, "Your Verification Code")
	return map[string]interface{}{
		"Verification": "The verification code has been sent to your email",
	}, nil
}

func UpdateVerificationCode(cache *caching.RedisClient, user UserSignIn) error {
	data, err := cache.Get(user.Email)
	if err != nil {
		return errors.New("invalid email")
	}

	var UserData UserVerify
	if err := json.Unmarshal([]byte(data), &UserData); err != nil {
		return err
	}

	NewCode := utils.GenerateVerficationCode()
	UserData.Code = NewCode

	expiryMinutes, err := strconv.Atoi(os.Getenv("VERIFY_CODE_EXPIRY_MINUTE"))
	if err != nil {
		return err
	}

	cache.Set(user.Email, UserData, time.Duration(expiryMinutes)*time.Minute)
	cache.Set(UserData.Username, 1, time.Duration(expiryMinutes)*time.Minute)

	SendMail(config.EmailSender, os.Getenv("GMAIL"), user.Email, NewCode, "Your Verification Code")
	return nil
}

func VerifyUser(ctx context.Context, db *pgxpool.Pool, cache *caching.RedisClient, user *UserVerify) (map[string]interface{}, error) {
	
	userCode := user.Code
	if err := cache.GetJson(user.Email, user); err != nil {
		return nil, err
	}
	
	if userCode != user.Code {
		return nil, fmt.Errorf("Wrong verification code")
	}

	// add the user into postgres
	transaction, err := db.Begin(ctx) // we should replace this with a middleware
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback(ctx)

	// set user as verified
	user.Verified = true

	if err := CreateUser(ctx, transaction, &user.User); err != nil {
		return nil, err
	}

	err = GetUser(ctx, transaction, user.Email, SELECT_ID_FROM_USER_BY_EMAIL, []interface{}{&user.ID}...)
	if err != nil {
		return nil, err
	}

	token, err := token.CreateAccessToken(&user.User, config.Env.AccessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	// remove user from the cache
	delResult, err := cache.Eval(ctx, luaDeleteScript, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	if delResult.(string) == "ERROR" {
		return nil, fmt.Errorf("error while removing user from cache")
	}

	if err := transaction.Commit(ctx); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"id":       user.OID, // sent to the client
		"email":    user.Email,
		"username": user.Username,
		"verified": user.Verified,
		"token":    token,
	}

	return data, nil
}
