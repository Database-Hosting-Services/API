package indexes

import (
	"DBHS/config"
	"DBHS/projects"
	"DBHS/utils"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
)

// generates a SQL query to create an index on a table
// CREATE INDEX index_name ON table_name USING index_type (column1, column2, column3);
func GenerateIndexQuery(Index IndexData) string {
	columns := strings.Join(Index.Columns, ", ")
	Index.IndexName = utils.ReplaceWhiteSpacesWithUnderscore(Index.IndexName)
	return "CREATE INDEX " + Index.IndexName + " ON " + Index.TableName + " USING " + Index.IndexType + " (" + columns + ")"
}

// generates a SQL query to delete an index
// DROP INDEX IF EXISTS index_name;
func GenerateDeleteIndexQuery(indexName string) string {
	return "DROP INDEX IF EXISTS " + indexName
}

// generates a SQL query to rename an index
// ALTER INDEX old_index_name RENAME TO new_index_name;
func GenerateRenameIndexQuery(oldName string, newName string) string {
	return "ALTER INDEX IF EXISTS " + oldName + " RENAME TO " + newName
}

func ProjectPoolConnection(ctx context.Context, db *pgxpool.Pool, UserID int, projectOid string) (*pgxpool.Pool, error) {
	projectDB, err := projects.GetUserSpecificProject(ctx, db, UserID, projectOid)
	if err != nil {
		return nil, err
	}

	if projectDB == nil {
		return nil, errors.New("project not found")
	}

	// ------------------------ Get The project connection Pool ------------------------

	DBname := strings.ToLower(projectDB.Name) + "_" + strconv.Itoa(UserID)
	conn, err := config.ConfigManager.GetDbConnection(ctx, DBname)
	if err != nil {
		return nil, err
	}

	// ------------------------ Check if the connection pool is already created ------------------------
	if conn == nil {
		return nil, errors.New("connection pool not found")
	}

	return conn, nil
}
