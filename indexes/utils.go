package indexes

import "strings"

// generates a SQL query to create an index on a table
// CREATE INDEX index_name ON table_name USING index_type (column1, column2, column3);
func GenerateIndexQuery(Index IndexData) string {
	columns := strings.Join(Index.Columns, ", ")
	return "CREATE INDEX " + Index.IndexName + " ON " + Index.TableName + " USING " + Index.IndexType + " (" + columns + ")"
}
