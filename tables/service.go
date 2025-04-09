package tables

import (
	"DBHS/config"
	"DBHS/utils"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTable(ctx context.Context, projectOID string, table *ClientTable, servDb *pgxpool.Pool) error {
	userId, ok := ctx.Value("user-id").(int)
	if !ok || userId == 0 {
		return errors.New("Unauthorized")
	}

	projectId, userDb, err := ExtractDb(ctx, projectOID, userId, servDb)
	if err != nil {
		return err
	}
	// create the table in the user db
	tx, err := userDb.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := CreateTableIntoHostingServer(ctx, table, tx); err != nil {
		return err
	}
	tableRecord := Table{
		Name:      table.TableName,
		ProjectID: projectId,
		OID:       utils.GenerateOID(),
	}
	var tableId int
	// insert table row into the tables table
	if err := InsertNewTable(ctx, &tableRecord, &tableId, servDb); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		DeleteTableRecord(ctx, tableId, servDb)
		return err
	}
	config.App.InfoLog.Printf("Table %s created successfully in project %s by user %s", table.TableName, projectOID, ctx.Value("user-name").(string))
	return nil
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
	table, err := ReadTable(ctx, userDb)
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
