package accounts

import (
	"DBHS/utils"
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

	userToken := utils.NewToken()
	userToken.AddClaims(map[string]interface{}{"id": user.ID})
	token, err := userToken.String()
	if err != nil {
		return nil, err
	}

	data := &map[string]interface{}{
		"id":       user.ID,
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
