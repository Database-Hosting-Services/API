package indexes

import (
	"DBHS/utils"
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
