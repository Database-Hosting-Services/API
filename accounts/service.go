package accounts

import (
	"DBHS/utils"
	"DBHS/utils/token"
	"DBHS/config"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignupUser(ctx context.Context, db *pgxpool.Pool, user *User) (*map[string]interface{}, error) {
	transaction, err := db.Begin(ctx) // we should replace this with a middleware
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback(ctx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user.OID = utils.GenerateOID()
	user.Password = string(hashedPassword)
	user.Verified = false

	if err := CreateUser(ctx, transaction, user); err != nil {
		return nil, err
	}

	if err := GetUserID(ctx, transaction, user); err != nil {
		return nil, err
	}

	token, err := token.CreateAccessToken(user, config.Env.AccessTokenExpiryHour)

	if err != nil {
		return nil, err
	}

	data := &map[string]interface{}{
		"id":       user.OID, 		//sent to the clinte
		"email":    user.Email,
		"username": user.Username,
		"verified": user.Verified,
		"token":    token,
	}

	if err := transaction.Commit(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func SignInUser(ctx context.Context, db *pgxpool.Pool, user *UserSignIn) (*map[string]interface{}, error) {
	authenticatedUser, err := GetUser(ctx, config.DB, user.Email)
	if err != nil {
		return nil, err;
	}
 
	if !CheckPasswordHash(user.Password, authenticatedUser.Password) {
		return nil, errors.New("InCorrect Email or Password")
	}

	var UserTokenData = &User{
		OID: authenticatedUser.OID,
		Username: authenticatedUser.Username,
	}
 
	token, err := token.CreateAccessToken(UserTokenData, config.Env.AccessTokenExpiryHour)
 
	if err != nil {
		return nil, err
	}
 
	resp := map[string]interface{}{
		"Data": map[string]interface{}{
			"oid":      authenticatedUser.OID,
			"username": authenticatedUser.Username,
			"email":    authenticatedUser.Email,
			"image":    authenticatedUser.Image,
			"token":    token,
		},
	}
 
	return &resp, nil
}
 