package tables

import (
	"DBHS/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func GetProjectNameID(ctx context.Context, projectId string, db utils.Querier) (interface{}, interface{}, error) {
	var name, id interface{}
	err := db.QueryRow(ctx,"SELECT id, name FROM projects WHERE oid = $1", projectId).Scan(&id, &name)
	if err != nil {
		return nil, nil, err
	}
	return name, id, nil
}

func InsertNewTable(ctx context.Context, table *Table, TableId *int, db utils.Querier) error {
	err := db.QueryRow(ctx, InsertNewTableRecordStmt, table.OID, table.Name, table.Description, table.ProjectID).Scan(TableId)
	if err != nil {
		return fmt.Errorf("failed to insert new table: %w", err)
	}
	return nil
}

func DeleteTableRecord(ctx context.Context, tableId int, db utils.Querier) error {
	_,err := db.Exec(ctx, fmt.Sprintf(DeleteTableStmt, "id"), tableId)
	if err != nil {
		return fmt.Errorf("failed to delete table record: %w", err)
	}
	return nil
}

func CheckOwnershipQuery(ctx context.Context, projectId string, userId int, db utils.Querier) (bool, error) {
	var count int
	err := db.QueryRow(ctx, CheckOwnershipStmt, projectId, userId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}

func ReadTable(ctx context.Context, db *pgxpool.Pool) (map[string]DbColumn, error) {
	var columns []DbColumn
	sqlxdb := sqlx.NewDb(stdlib.OpenDBFromPool(db), "pgx")
	err := sqlxdb.SelectContext(ctx, &columns, ReadTableStmt)
	if err != nil {
		return nil, fmt.Errorf("failed to read table: %w", err)
	}
	columnsMap := make(map[string]DbColumn)
	for _, column := range columns {
		columnsMap[column.Name] = column
	}
	return columnsMap, nil
}

func GetTableName(ctx context.Context, tableOID string, db utils.Querier) (string, error) {
	var tableName string
	err := db.QueryRow(ctx, GetTableNameStmt, tableOID).Scan(&tableName)
	if err != nil {
		return "", fmt.Errorf("failed to get table name: %w", err)
	}
	return tableName, nil
}

func DeleteTableFromHostingServer(ctx context.Context, tableName string, db utils.Querier) error {
	_, err := db.Exec(ctx, DropTableStmt, tableName)
	if err != nil {
		return fmt.Errorf("failed to delete table from hosting server: %w", err)
	}
	return nil
}

func DeleteTableFromServerDb(ctx context.Context, tableOID string, db utils.Querier) error {
	_, err := db.Exec(ctx, fmt.Sprintf(DeleteTableStmt, "oid"), tableOID)
	if err != nil {
		return fmt.Errorf("failed to delete table from server DB: %w", err)
	}
	return nil
}