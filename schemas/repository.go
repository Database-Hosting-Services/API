package schemas

import (
	"DBHS/projects"
	"DBHS/utils"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func getDatabaseByName(ctx context.Context, DB utils.Querier, name string) (projects.DatabaseConfig, error) {
	var database projects.DatabaseConfig
	err := pgxscan.Get(ctx, DB, &database, GetDatabaseByName, name)
	if err != nil {
		return projects.DatabaseConfig{}, err
	}
	return database, nil
}

func getSchema(ctx context.Context, DB utils.Querier, projectName string) (string, error) {
	var schema string
	err := DB.QueryRow(ctx, GetAllTablesSchema, projectName).Scan(&schema)
	if err != nil {
		return "", err
	}
	return schema, nil
}
