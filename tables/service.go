package tables

import (
	"DBHS/config"
	"DBHS/utils"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTable(ctx context.Context, projectOID string, table *ClientTable, servDb *pgxpool.Pool) error {
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return errors.New("Unauthorized")
	}
	config.App.InfoLog.Println(UserID)
	// get the dbname to connect to
	dbName, projectId, err := GetProjcetNameID(ctx, projectOID, servDb)
	if err != nil {
		return err
	}
	// get the db connection
	userDb, err := config.ConfigManager.GetDbConnection(ctx, utils.UserServerDbFromat(dbName.(string), UserID))
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
		ProjectID: projectId.(int64),
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
