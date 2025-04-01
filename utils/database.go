package utils

import (
	"context"
)

func UpdateDataInDatabase(ctx context.Context, db Querier, query string, dest ...interface{}) error {
	_, err := db.Exec(ctx, query, dest...)
	return err
}
