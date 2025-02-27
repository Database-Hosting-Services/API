package accounts

import (
	"DBHS/utils"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignupUser(ctx context.Context, db *pgxpool.Pool, user *User) error {
	transaction, err := db.Begin(ctx) // we should replace this with a middleware
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.OID = utils.GenerateOID()
	user.Password = string(hashedPassword)
	user.Verified = false

	ret := CreateUser(ctx, transaction, user)
	if err := transaction.Commit(ctx); err != nil {
		return err
	}
	return ret
}
