package accounts

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func CreateUser(ctx context.Context, db pgx.Tx, user *User) error {
	query := `INSERT INTO "User" (oid, username, email, password, image, verified, created_at, last_login)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.Exec(ctx, query,
		user.OID,
		user.Username,
		user.Email,
		user.Password,
		user.Image,
		user.Verified,
		user.CreatedAt,
		user.LastLogin,
	)

	return err
}
