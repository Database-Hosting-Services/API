package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(ctx context.Context, db pgx.Tx, user *User) error {
	query := `INSERT INTO "users" (oid, username, email, password, image, verified, created_at, last_login)
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
			  FROM "users" WHERE email = $1`

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
			return fmt.Errorf("user with Email %s not found", user.Email)
		}
		return err
	}

	return nil
}

func GetUserID(ctx context.Context, db pgx.Tx, user *User) error {
	query := `SELECT id
			  FROM "users" WHERE username = $1`
	err := db.QueryRow(ctx, query, user.Username).Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with username %s not found", user.Username)
		}
		return err
	}
	return nil
}

/*

GetUser queries the database to fetch a user based on a dynamic search field.

Arguments:
 * SearchField: The field value used to identify the user (email, username, etc.).
 * query: The parameterized SQL query with a single placeholder ($1).

*/

func GetUser(ctx context.Context, db *pgxpool.Pool, SearchField string, query string) (User, error) {
	var user User
	err := db.QueryRow(ctx, query, SearchField).Scan(
		&user.ID,
		&user.OID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.CreatedAt,
		&user.LastLogin,
	)
	return user, err
}
