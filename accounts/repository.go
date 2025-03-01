package accounts

import (
	"context"
	"database/sql"
	"fmt"
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

func GetUserByEmail(ctx context.Context, db pgx.Tx, user *User) error {
	query := `SELECT
				id, oid, username, email, image, verified, created_at, last_login
			  FROM "User" WHERE email = $1`

	err := db.QueryRow(ctx, query, user.Email).Scan(
		&user.ID,
		&user.OID,
		&user.Username,
		&user.Email,
		&user.Image,
		&user.Verified,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with ID %d not found", user.Email)
		}
		return err
	}

	return nil
}
