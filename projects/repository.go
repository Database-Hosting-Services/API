package projects

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckDatabaseExists(ctx context.Context, db *pgxpool.Pool, query string, SearchField ...interface{}) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, query, SearchField...).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func InsertNewRecord(ctx context.Context, db *pgxpool.Pool, query string, values ...interface{}) error {
	_, err := db.Exec(ctx, query, values...)
	if err != nil {
		return err
	}
	return nil
}
