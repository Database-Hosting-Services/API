package tables

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

func GetProjectNameID(ctx context.Context, projectId string, db Querier) (interface{}, interface{}, error) {
	var name, id interface{}
	err := db.QueryRow(ctx,"SELECT id, name FROM projects WHERE oid = $1", projectId).Scan(&id, &name)
	if err != nil {
		return nil, nil, err
	}
	return name, id, nil
}

func InsertNewTable(ctx context.Context, table *Table, TableId *int, db Querier) error {
	err := db.QueryRow(ctx, InsertNewTableRecordStmt, table.OID, table.Name, table.Description, table.ProjectID).Scan(TableId)
	if err != nil {
		return fmt.Errorf("failed to insert new table: %w", err)
	}
	return nil
}

func DeleteTableRecord(ctx context.Context, tableId int, db Querier) error {
	_,err := db.Exec(ctx, DeleteTableRecordStmt, tableId)
	if err != nil {
		return fmt.Errorf("failed to delete table record: %w", err)
	}
	return nil
}

func CheckOwnership(ctx context.Context, projectId string, userId int, db Querier) (bool, error) {
	var count int
	err := db.QueryRow(ctx, CheckOwnershipStmt, projectId, userId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}