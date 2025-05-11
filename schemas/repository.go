package schemas

import (
	"DBHS/projects"
	"DBHS/utils"
	"context"
	"database/sql"
	"fmt"

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

func getDatabaseTableName(ctx context.Context, DB utils.Querier, tableOID string) (string, error) {
	var tableName string
	err := pgxscan.Get(ctx, DB, &tableName, GetTableNameByOID, tableOID)
	if err != nil {
		return "", err
	}
	return tableName, nil
}

func getDatabaseSchema(ctx context.Context, DB utils.Querier) (*SchemaResponse, error) {
	rows, err := DB.Query(ctx, GetAllTablesSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	schemaMap := make(map[string]*TableSchema)

	for rows.Next() {
		var (
			tableName     string
			columnName    string
			dataType      string
			isNullable    bool
			isPrimaryKey  bool
			isUnique      bool
			foreignTable  sql.NullString
			foreignColumn sql.NullString
		)

		err := rows.Scan(
			&tableName,
			&columnName,
			&dataType,
			&isNullable,
			&isPrimaryKey,
			&isUnique,
			&foreignTable,
			&foreignColumn,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, exists := schemaMap[tableName]; !exists {
			schemaMap[tableName] = &TableSchema{
				TableName: tableName,
				Cols:      []Column{},
			}
		}

		// Create column with foreign key, if not exists it will be NULL
		column := Column{
			Name:         columnName,
			Type:         dataType,
			IsNullable:   isNullable,
			IsPrimaryKey: isPrimaryKey,
			IsUnique:     isUnique,
		}

		if foreignTable.Valid && foreignColumn.Valid {
			column.ForeignKeys = append(column.ForeignKeys, ForeignKey{
				TableName:  foreignTable.String,
				ColumnName: foreignColumn.String,
			})
		}

		schemaMap[tableName].Cols = append(schemaMap[tableName].Cols, column)
	}

	// Check for errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	// Convert map to slice for response
	response := &SchemaResponse{}
	for _, table := range schemaMap {
		response.Schema = append(response.Schema, *table)
	}

	return response, nil
}

func GetTableSchema(ctx context.Context, DB utils.Querier, tableName string) (*SchemaResponse, error) {
	rows, err := DB.Query(ctx, GetTableSchemaQuery, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query table schema: %w", err)
	}
	defer rows.Close()

	var (
		columns     []Column
		scannedName string
		hasRows     bool
	)

	for rows.Next() {
		hasRows = true
		var (
			col           Column
			foreignTable  sql.NullString
			foreignColumn sql.NullString
		)

		err := rows.Scan(
			&scannedName,
			&col.Name,
			&col.Type,
			&col.IsNullable,
			&col.IsPrimaryKey,
			&col.IsUnique,
			&foreignTable,
			&foreignColumn,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Initialize foreign keys slice
		col.ForeignKeys = make([]ForeignKey, 0)

		// Add foreign key
		if foreignTable.Valid && foreignColumn.Valid {
			col.ForeignKeys = append(col.ForeignKeys, ForeignKey{
				TableName:  foreignTable.String,
				ColumnName: foreignColumn.String,
			})
		}

		columns = append(columns, col)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if !hasRows {
		return nil, fmt.Errorf("table '%s' not found", tableName)
	}

	return &SchemaResponse{
		Schema: []TableSchema{{
			TableName: scannedName,
			Cols:      columns,
		}},
	}, nil
}
