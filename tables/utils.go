package tables

import (
	"DBHS/utils"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
)

func CreateColumnDefinition(column *Column) string {
	res := fmt.Sprintf("%s %s", column.Name, column.Type)
	if column.IsPrimaryKey != nil && *column.IsPrimaryKey {
		res += " PRIMARY KEY"
	}
	if column.IsUnique != nil && *column.IsUnique {
		res += " UNIQUE"
	}
	if column.IsNullable != nil && !*column.IsNullable {
		res += " NOT NULL"
	}
	if column.ForeignKey.TableName != "" {
		res += fmt.Sprintf(" REFERENCES %s(%s)", column.ForeignKey.TableName, column.ForeignKey.ColumnName)
	}
	return res
}

func ParseTableIntoSQLCreate(table *ClientTable) (string, error) {
	columns := make([]string, len(table.Columns))
	for i, column := range table.Columns {
		columns[i] = CreateColumnDefinition(&column)
	}
	createTableSQL := fmt.Sprintf("CREATE TABLE %s (%s);", table.TableName, strings.Join(columns, ", "))
	return createTableSQL, nil
}

func CreateTableIntoHostingServer(ctx context.Context, table *ClientTable, tx pgx.Tx) error {
	DDLQuery, err := ParseTableIntoSQLCreate(table)
	if err != nil {
		return err
	}
	log.Println(DDLQuery)
	_, err = tx.Exec(ctx, DDLQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func CheckForValidTable(table *ClientTable) bool {
	if table.TableName == "" || len(table.Columns) == 0 {
		return false
	}
	for _, column := range table.Columns {
		if column.Name == "" || column.Type == "" {
			return false
		}
	}
	return true
}

func ExecuteUpdate(tableName string, table map[string]DbColumn, updates *TableUpdate, db utils.Querier) error {
	// inserts
	insertStmt := "ALTER TABLE %s ADD COLUMN %s"
	for _, insert := range updates.Inserts.Columns {
		column := CreateColumnDefinition(&insert)
		query := fmt.Sprintf(insertStmt, tableName, column)
		if _, err := db.Exec(context.Background(), query); err != nil {
			return fmt.Errorf("failed to insert column: %w", err)
		}
	}
	// updates
	for _, update := range updates.Updates {
		if update.Update.Name != "" {
			query := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", tableName, update.Name, update.Update.Name)
			if _, err := db.Exec(context.Background(), query); err != nil {
				return fmt.Errorf("failed to update column: %w", err)
			}
			update.Name = update.Update.Name
		}

		if update.Update.Type != "" {
			query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", tableName, update.Name, update.Update.Type)
			if _, err := db.Exec(context.Background(), query); err != nil {
				return fmt.Errorf("failed to update column: %w", err)
			}
		}

		if update.Update.IsNullable != nil && table[update.Name].IsNullable != *update.Update.IsNullable {
			// If the column is nullable, we need to drop the NOT NULL constraint
			// If the column is not nullable, we need to add the NOT NULL constraint
			option := "DROP"
			if !*update.Update.IsNullable {
				option = "SET"
			}
			query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s %s NOT NULL", tableName, update.Name, option)
			if _, err := db.Exec(context.Background(), query); err != nil {
				return fmt.Errorf("failed to update column: %w", err)
			}
		}

		if update.Update.IsUnique != nil {
			if *update.Update.IsUnique && table[update.Name].UniqueConstraintType == nil {
				// Add unique constraint
				query := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)", tableName, CreateUniqueConstraintName(tableName, update.Name), update.Name)
				if _, err := db.Exec(context.Background(), query); err != nil {
					return fmt.Errorf("failed to add unique constraint: %w", err)
				}
			} else if !*update.Update.IsUnique && table[update.Name].UniqueConstraintType != nil {
				// Drop unique constraint
				query := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", tableName, *table[update.Name].UniqueConstraintName)
				if _, err := db.Exec(context.Background(), query); err != nil {
					return fmt.Errorf("failed to drop unique constraint: %w", err)
				}
			}
		}
	}

	// deletes
	for _, delete := range updates.Deletes {
		query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, delete)
		if _, err := db.Exec(context.Background(), query); err != nil {
			return fmt.Errorf("failed to delete column: %w", err)
		}
	}

	return nil
}

func CreateUniqueConstraintName(tableName string, columnName string) string {
	return fmt.Sprintf("%s_%s_key", tableName, columnName)
}
