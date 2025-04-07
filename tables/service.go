package tables

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)



func CreateTable(ctx context.Context, table *Table, servDb *pgxpool.Pool) error {
	// get the dbname to connect to
	dbName, err := GetProjcetFeild(ctx, table.ProjectID, "name", servDb) 
	if err != nil {
		return err
	}
	// get the db connection
	

}