package tables

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllTables(ctx context.Context, projectOID string, servDb *pgxpool.Pool) ([]ShortTable, error) {
	userId, ok := ctx.Value("user-id").(int)
	if !ok || userId == 0 {
		return nil, errors.New("Unauthorized")
	}

	_, projectId, err := GetProjectNameID(ctx, projectOID, servDb)
	if err != nil {
		return nil, err
	}

	tables, err := GetAllTablesNameOid(ctx, projectId.(int64), servDb)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func CreateTable(ctx context.Context, projectOID string, table *ClientTable, servDb *pgxpool.Pool) (string, error) {
	userId, ok := ctx.Value("user-id").(int)
	if !ok || userId == 0 {
		return "", errors.New("Unauthorized")
	}

	projectId, userDb, err := ExtractDb(ctx, projectOID, userId, servDb)
	if err != nil {
		return "", err
	}
	// create the table in the user db
	tx, err := userDb.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	if err := CreateTableIntoHostingServer(ctx, table, tx); err != nil {
		return "", err
	}
	tableRecord := Table{
		Name:      table.TableName,
		ProjectID: projectId,
		OID:       utils.GenerateOID(),
	}
	var tableId int
	// insert table row into the tables table
	if err := InsertNewTable(ctx, &tableRecord, &tableId, servDb); err != nil {
		return "", err
	}
	if err := tx.Commit(ctx); err != nil {
		DeleteTableRecord(ctx, tableId, servDb)
		return "", err
	}
	config.App.InfoLog.Printf("Table %s created successfully in project %s by user %s", table.TableName, projectOID, ctx.Value("user-name").(string))
	return tableRecord.OID, nil
}

func UpdateTable(ctx context.Context, projectOID string, tableOID string, updates *TableUpdate, servDb *pgxpool.Pool) error {
	userId, ok := ctx.Value("user-id").(int)
	if !ok || userId == 0 {
		return errors.New("Unauthorized")
	}

	_, userDb, err := ExtractDb(ctx, projectOID, userId, servDb)
	if err != nil {
		return err
	}

	tableName, err := GetTableName(ctx, tableOID, servDb)
	if err != nil {
		return err
	}

	// Call the service function to read the table
	table, err := ReadTableColumns(ctx, tableName, userDb)
	if err != nil {
		return err
	}

	tx, err := userDb.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := ExecuteUpdate(tableName, table, updates, tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func DeletTable(ctx context.Context, projectOID, tableOID string, servDb *pgxpool.Pool) error {
	userId, ok := ctx.Value("user-id").(int)
	if !ok || userId == 0 {
		return errors.New("Unauthorized")
	}

	_, userDb, err := ExtractDb(ctx, projectOID, userId, servDb)
	if err != nil {
		return err
	}

	tableName, err := GetTableName(ctx, tableOID, servDb)
	if err != nil {
		return err
	}
	usertx, err := userDb.Begin(ctx)
	if err != nil {
		return err
	}
	defer usertx.Rollback(ctx)

	if err := DeleteTableFromHostingServer(ctx, tableName, usertx); err != nil {
		return err
	}

	servtx, err := servDb.Begin(ctx)
	if err != nil {
		return err
	}
	defer servtx.Rollback(ctx)
	if err := DeleteTableFromServerDb(ctx, tableOID, servtx); err != nil {
		return err
	}

	if err := servtx.Commit(ctx); err != nil {
		return err
	}

	if err := usertx.Commit(ctx); err != nil {
		config.App.ErrorLog.Println("Failed to commit transaction:", err, "server db and user db are not in sync")
		return err
	}

	config.App.InfoLog.Printf("Table %s deleted successfully in project %s by user %s", tableName, projectOID, ctx.Value("user-name").(string))

	return nil
}

/*
	the responce data will be in the form
	{
		"columns": [ // array of columns names and types
			{
				"name": "value",
				"type": "type"
			}
			...
		]
		"rows": [
			{
				"column1_name": "value"
				"column2_name": "value"
				...
			}
		]

	}

*/

func ReadTable(ctx context.Context, projectOID, tableOID string, parameters map[string][]string, servDb *pgxpool.Pool) (*Data, error) {
	userId, ok := ctx.Value("user-id").(int)
	if !ok || userId == 0 {
		return nil, response.ErrUnauthorized
	}

	_, userDb, err := ExtractDb(ctx, projectOID, userId, servDb)
	if err != nil {
		return nil, err
	}

	tableName, err := GetTableName(ctx, tableOID, servDb)
	if err != nil {
		return nil, err
	}

	data, err := ReadTableData(ctx, tableName, parameters, userDb)
	if err != nil {
		return nil, err
	}

	return data, nil
}
