package utils

import (
	"context"
	"fmt"
)

var (
	CheckOwnershipStmt = `SELECT COUNT(*) FROM "projects" WHERE oid = $1 AND owner_id = $2;`
)

func UpdateDataInDatabase(ctx context.Context, db Querier, query string, dest ...interface{}) error {
	_, err := db.Exec(ctx, query, dest...)
	return err
}

func CheckOwnershipQuery(ctx context.Context, projectId string, userId int, db Querier) (bool, error) {
	var count int
	err := db.QueryRow(ctx, CheckOwnershipStmt, projectId, userId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}
