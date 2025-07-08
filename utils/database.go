package utils

import (
	"DBHS/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	CheckOwnershipStmt = `SELECT COUNT(*) FROM "projects" WHERE oid = $1 AND owner_id = $2;`
	CheckProjectExistStmt = `SELECT COUNT(*) FROM "projects" WHERE oid = $1;`
	CheckOwnershipTableStmt = `SELECT COUNT(*) FROM "Ptable" WHERE oid = $1 AND project_id = $2;`
	CheckTableExistStmt = `SELECT COUNT(*) FROM "Ptable" WHERE oid = $1;`
)

func UpdateDataInDatabase(ctx context.Context, db Querier, query string, dest ...interface{}) error {
	_, err := db.Exec(ctx, query, dest...)
	return err
}

func CheckOwnershipQuery(ctx context.Context, projectOID string, userId int64, db Querier) (bool, error) {

	var count int
	err := db.QueryRow(ctx, CheckOwnershipStmt, projectOID, userId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}

func CheckProjectExist(ctx context.Context, projectOID string, db Querier) (bool, error) {
	var count int
	err := db.QueryRow(ctx, CheckProjectExistStmt, projectOID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check project existence: %w", err)
	}
	return count > 0, nil
}

func CheckOwnershipQueryTable(ctx context.Context, tableOID string, projectID int64, db Querier) (bool, error) {
	var count int
	err := db.QueryRow(ctx, CheckOwnershipTableStmt, tableOID, projectID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}

func CheckTableExist(ctx context.Context, tableOID string, db Querier) (bool, error) {
	var count int
	err := db.QueryRow(ctx, CheckTableExistStmt, tableOID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}
	return count > 0, nil
}

func GetProjectNameID(ctx context.Context, projectOID string, db Querier) (interface{}, interface{}, error) {
	var name, id interface{}
	err := db.QueryRow(ctx, "SELECT id, name FROM projects WHERE oid = $1", projectOID).Scan(&id, &name)
	if err != nil {
		return nil, nil, err
	}
	return name, id, nil
}

func ExtractDb(ctx context.Context, projectOID string, UserID int64, servDb *pgxpool.Pool) (int64, *pgxpool.Pool, error) {
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