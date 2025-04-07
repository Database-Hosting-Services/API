package tables

import (
	"DBHS/config"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)



func CreateTable(ctx context.Context,projectID string, table *ClientTable, servDb *pgxpool.Pool) error {
	// get the dbname to connect to
	dbName, err := GetProjcetFeild(ctx, projectID, "name", servDb) 
	if err != nil {
		return err
	}
	// get the db connection
	userDb, err := config.ConfigManager.GetDbConnection(ctx, dbName.(string))
	if err != nil {
		return err
	}
	// create the table in the user db
	tx, err := userDb.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := CreateTableSQL(ctx, table, tx) ; err != nil {
		return err
	}

	// insert table row into the tables table
	

}