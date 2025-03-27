package accounts

import (
	"context"
	"fmt"

	"DBHS/utils"
	"github.com/jackc/pgx/v5"
)

func CreateUser(ctx context.Context, db pgx.Tx, user *User) error {
	query := `INSERT INTO "users" (oid, username, email, password, image, created_at, last_login)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Exec(ctx, query,
		user.OID,
		user.Username,
		user.Email,
		user.Password,
		user.Image,
		user.CreatedAt,
		user.LastLogin,
	)

	return err
}

/*

	GetUser executes a SQL query to retrieve a single user record based on a search field,
	and scans the result into the provided destination interface(s).

	Parameters:
	- ctx:        Context for request cancellation/timeout
	- db:         Database querier interface (pgxpool.Pool or pgx.Tx)
	- SearchField: Value used for the WHERE clause search (e.g., user ID, email)
	- query:      SQL query string containing a single placeholder ($1)
	- dest:       Variadic slice of pointers to scan results into (must match query columns)

  Example usage:
	err = GetUser(ctx, transaction, user.Email, SELECT_ID_FROM_USER_BY_EMAIL, []interface{}{&user.ID}...);

*/

func GetUser(ctx context.Context, db utils.Querier, SearchField string, query string, dest ...interface{}) error {
	err := db.QueryRow(ctx, query, SearchField).Scan(dest...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user with %s not found", SearchField)
		}
		return err
	}
	return nil
}

func UpdateDataInDatabase(ctx context.Context, db utils.Querier, query string, dest ...interface{}) error {
	_, err := db.Exec(ctx, query, dest...)
	return err
}
