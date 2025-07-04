package tables

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllTables(ctx context.Context, projectOID string, servDb *pgxpool.Pool) ([]Table, error) {
	userId, ok := ctx.Value("user-id").(int64)
	if !ok || userId == 0 {
		return nil, response.ErrUnauthorized
	}

	projectId, userDb, err := utils.ExtractDb(ctx, projectOID, userId, servDb)
	if err != nil {
		return nil, err
	}

	tables, err := GetAllTablesRepository(ctx, projectId, userDb, servDb)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func GetTableSchema(ctx context.Context, projectOID string, tableOID string, servDb *pgxpool.Pool) (*Table, error) {
	userId, ok := ctx.Value("user-id").(int64)
	if !ok || userId == 0 {
		return nil, response.ErrUnauthorized
	}

	_, userDb, err := utils.ExtractDb(ctx, projectOID, userId, servDb)
	if err != nil {
		return nil, err
	}

	tableName, err := GetTableName(ctx, tableOID, servDb)
	if err != nil {
		return nil, err
	}

	schema, err := utils.GetTable(ctx, tableName, userDb)
	if err != nil {
		return nil, err
	}

	return &Table{
		Schema: schema,
		OID:    tableOID,
		Name:   schema.TableName,
	}, nil
}

func CreateTable(ctx context.Context, projectOID string, table *Table, servDb *pgxpool.Pool) (string, error) {
	userId, ok := ctx.Value("user-id").(int64)
	if !ok || userId == 0 {
		return "", response.ErrUnauthorized
	}

	projectId, userDb, err := utils.ExtractDb(ctx, projectOID, userId, servDb)
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

	table.OID = utils.GenerateOID()
	table.ProjectID = projectId
	var tableId int64
	// insert table row into the tables table
	if err := InsertNewTable(ctx, table, &tableId, servDb); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		DeleteTableRecord(ctx, tableId, servDb)
		return "", err
	}
	config.App.InfoLog.Printf("Table %s created successfully in project %s by user %s", table.Name, projectOID, ctx.Value("user-name").(string))
	return table.OID, nil
}

func UpdateTable(ctx context.Context, projectOID string, tableOID string, newSchema *UpdateTableSchema, servDb *pgxpool.Pool) error {
	userId, ok := ctx.Value("user-id").(int64)
	if !ok || userId == 0 {
		return response.ErrUnauthorized
	}

	_, userDb, err := utils.ExtractDb(ctx, projectOID, userId, servDb)
	if err != nil {
		return err
	}

	oldSchema, err := utils.GetTableSchema(ctx, tableOID, servDb)
	if err != nil {
		return err
	}

	DDLUpdate, err := utils.CompareTableSchemas(oldSchema, newSchema.Schema, newSchema.Renames)
	if err != nil {
		return err
	}

	tx, err := userDb.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, DDLUpdate); err != nil {
		return err
	}

	// Update the table name in the service database
	if oldSchema.TableName != newSchema.Schema.TableName {
		if _, err := servDb.Exec(ctx, "UPDATE \"Ptable\" SET name = $1 WHERE oid = $2",
			newSchema.Schema.TableName, tableOID); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func DeleteTable(ctx context.Context, projectOID, tableOID string, servDb *pgxpool.Pool) error {
	userId, ok := ctx.Value("user-id").(int64)
	if !ok || userId == 0 {
		return response.ErrUnauthorized
	}

	_, userDb, err := utils.ExtractDb(ctx, projectOID, userId, servDb)
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
	userId, ok := ctx.Value("user-id").(int64)
	if !ok || userId == 0 {
		return nil, response.ErrUnauthorized
	}

	_, userDb, err := utils.ExtractDb(ctx, projectOID, userId, servDb)
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
