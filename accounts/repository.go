package accounts

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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

// Querier defines an interface for executing a single-row query.
// Both *pgxpool.Pool and pgx.Tx implement this interface through the QueryRow method.
type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
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

func GetUser(ctx context.Context, db Querier, SearchField string, query string, dest ...interface{}) error {
	err := db.QueryRow(ctx, query, SearchField).Scan(dest...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user with %s not found", SearchField)
		}
		return err
	}
	return nil
}
