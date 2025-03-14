package accounts

import (
	"DBHS/caching"
	"DBHS/config"
	"DBHS/utils"
	"DBHS/utils/token"
	"github.com/jackc/pgx/v5"

	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignupUser(ctx context.Context, db *pgxpool.Pool, user *UserUnVerified) error {
	/*
		store user's data in cache and
		generate a verification code and send it to the user
	*/
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Code = utils.GenerateVerficationCode()
	user.OID = utils.GenerateOID()
	user.Password = string(hashedPassword)

	// store user's data in cache
	config.VerifyCache.Set(user.Email, user, time.Minute*30)
	config.VerifyCache.Set(user.Username, true, time.Minute*30)

	// send the verification code
	if err = SendMail(config.EmailSender, os.Getenv("GMAIL"), user.Email, user.Code, "Verification Code"); err != nil {
		config.VerifyCache.Delete(user.Email)
		config.VerifyCache.Delete(user.Username)
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
		if err.Error() == "user with "+user.Email+" not found" {
			return nil, errors.New("InCorrect Email or Password")
		}
		return nil, err
	}

	if !CheckPasswordHash(user.Password, authenticatedUser.Password) {
		return nil, errors.New("InCorrect Email or Password")
	}

	UserTokenData := User{
		OID:      authenticatedUser.OID,
		Username: authenticatedUser.Username,
	}

	token, err := token.CreateAccessToken(&UserTokenData, config.Env.AccessTokenExpiryHour)
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
	var user UserUnVerified

	_, err := cache.Get(email, &user)
	if err != nil {
		return nil, errors.New("invalid email or password")
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
	var UserData UserUnVerified

	_, err := cache.Get(user.Email, &user)
	if err != nil {
		return errors.New("invalid email")
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

func UpdateUserPassword(ctx context.Context, db *pgxpool.Pool, UserPassword *UpdatePasswordModel) error {
	if UserPassword.Password != UserPassword.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	UserId, ok := ctx.Value("user-id").(string)
	if !ok || UserId == "" {
		return errors.New("Unauthorized")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UserPassword.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	err = UpdateUserPasswordInDatabase(ctx, db, UserId, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func VerifyUser(ctx context.Context, db *pgxpool.Pool, cache *caching.RedisClient, user *UserUnVerified) (map[string]interface{}, error) {

	userCode := user.Code
	if _, err := cache.Get(user.Email, user); err != nil {
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

	if err := CreateUser(ctx, transaction, &user.User); err != nil {
		return nil, err
	}

	if err := GetUser(ctx, transaction, user.Email, SELECT_ID_FROM_USER_BY_EMAIL, []interface{}{&user.ID}...); err != nil {
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
		"token":    token,
	}

	return data, nil
}

func ForgetPasswordService(ctx context.Context, db *pgxpool.Pool, cache *caching.RedisClient, email string) error {
	// check if a user exist with this email
	var user UserUnVerified
	err := GetUser(ctx, db, email, SELECT_USER_BY_Email, []interface{}{
		&user.ID,
		&user.OID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.CreatedAt,
		&user.LastLogin,
	}...)

	if err != nil {
		return fmt.Errorf("User does not exist")
	}

	code := utils.GenerateVerficationCode()
	user.Code = code

	if err := cache.Set("forget:"+user.Email, &user, time.Minute*time.Duration(config.Env.VerifyCodeExpiryMinute)); err != nil {
		return err
	}

	if err := SendMail(config.EmailSender, os.Getenv("GMAIL"), email, code, "Verifacation Code"); err != nil {
		cache.Delete("forget:" + user.Email)
		return err
	}
	return nil
}

func ForgetPasswordVerifyService(ctx context.Context, db *pgxpool.Pool, cache *caching.RedisClient, resetForm *ResetPasswordForm) error {
	var user UserUnVerified
	if _, err := cache.Get("forget:"+resetForm.Email, &user); err != nil {
		return err
	}
	if user.Code != resetForm.Code {
		return fmt.Errorf("Wrong verification code")
	}

	if err := GetUser(ctx, db, resetForm.Email, SELECT_USER_BY_Email, []interface{}{
		&user.ID,
		&user.OID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.CreatedAt,
		&user.LastLogin,
	}...); err != nil {
		return err
	}

	transaction, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)

	if err := UpdateUserPasswordInDatabase(ctx, transaction, user.OID, utils.HashedPassword(resetForm.Password)); err != nil {
		return err
	}

	if err := cache.Delete("forget:" + resetForm.Email); err != nil {
		return err
	}

	if err := transaction.Commit(ctx); err != nil {
		cache.Set("forget:"+user.Email, &user, time.Minute*time.Duration(config.Env.VerifyCodeExpiryMinute))
		return err
	}
	return nil
}

func UpdateUserData(ctx context.Context, db pgx.Tx, query string, args []interface{}) error {
	if err := UpdateUserDataInDatabase(ctx, db, query, args); err != nil {
		return err
	}
	return nil
}
