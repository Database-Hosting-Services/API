package utils

import (
	"DBHS/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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

func GetProjectNameID(ctx context.Context, projectId string, db Querier) (interface{}, interface{}, error) {
	var name, id interface{}
	err := db.QueryRow(ctx, "SELECT id, name FROM projects WHERE oid = $1", projectId).Scan(&id, &name)
	if err != nil {
		return nil, nil, err
	}
	return name, id, nil
}

func ExtractDb(ctx context.Context, projectOID string, UserID int, servDb *pgxpool.Pool) (int64, *pgxpool.Pool, error) {
	// get the dbname to connect to
	dbName, projectId, err := GetProjectNameID(ctx, projectOID, servDb)
	if err != nil {
		return 0, nil, err
	}
	// get the db connection
	userDb, err := config.ConfigManager.GetDbConnection(ctx, UserServerDbFormat(dbName.(string), UserID))
	if err != nil {
		return 0, nil, err
	}

	return projectId.(int64), userDb, nil
}