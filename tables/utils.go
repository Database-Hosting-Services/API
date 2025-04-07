package tables

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func ParseTableIntoSQLCreate(table *ClientTable) (string, error) {
	columns := make([]string, len(table.Columns))
	for i, column := range table.Columns {
		columns[i] = fmt.Sprintf("%s %s", column.Name, column.Type)
		if column.IsPrimaryKey {
			columns[i] += " PRIMARY KEY"
		}
		if column.IsUnique {
			columns[i] += " UNIQUE"
		}
		if !column.IsNullable {
			columns[i] += " NOT NULL"
		}
		if column.ForeignKey.TableName != "" {
			columns[i] += fmt.Sprintf(" REFERENCES %s(%s)", column.ForeignKey.TableName, column.ForeignKey.ColumnName)
		}
	}
	createTableSQL := fmt.Sprintf("CREATE TABLE %s (%s);", table.TableName, strings.Join(columns, ", "))
	return createTableSQL, nil
}

func CreateTableIntoHostingServer(ctx context.Context, table *ClientTable, tx pgx.Tx) (error) {
	DDLQuery, err := ParseTableIntoSQLCreate(table)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, DDLQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}