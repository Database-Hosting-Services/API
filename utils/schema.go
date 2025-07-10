package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
)

// TableColumn represents a database column with its properties
type TableColumn struct {
	TableName              string  `db:"table_name" json:"TableName"`
	ColumnName             string  `db:"column_name" json:"ColumnName"`
	DataType               string  `db:"data_type" json:"DataType"`
	IsNullable             bool    `db:"is_nullable" json:"IsNullable"`
	ColumnDefault          *string `db:"column_default" json:"ColumnDefault"`
	CharacterMaximumLength *int    `db:"character_maximum_length" json:"CharacterMaximumLength"`
	NumericPrecision       *int    `db:"numeric_precision" json:"NumericPrecision"`
	NumericScale           *int    `db:"numeric_scale" json:"NumericScale"`
	OrdinalPosition        int     `db:"ordinal_position" json:"OrdinalPosition"`
}

// ConstraintInfo represents database constraints
type ConstraintInfo struct {
	TableName         string  `db:"table_name" json:"TableName"`
	ConstraintName    string  `db:"constraint_name" json:"ConstraintName"`
	ConstraintType    string  `db:"constraint_type" json:"ConstraintType"`
	ColumnName        *string `db:"column_name" json:"ColumnName"`
	ForeignTableName  *string `db:"foreign_table_name" json:"ForeignTableName"`
	ForeignColumnName *string `db:"foreign_column_name" json:"ForeignColumnName"`
	CheckClause       *string `db:"check_clause" json:"CheckClause"`
	OrdinalPosition   *int    `db:"ordinal_position" json:"OrdinalPosition"`
}

// IndexInfo represents database indexes
type IndexInfo struct {
	TableName  string `db:"table_name" json:"TableName"`
	IndexName  string `db:"index_name" json:"IndexName"`
	ColumnName string `db:"column_name" json:"ColumnName"`
	IsUnique   bool   `db:"is_unique" json:"IsUnique"`
	IndexType  string `db:"index_type" json:"IndexType"`
	IsPrimary  bool   `db:"is_primary" json:"IsPrimary"`
}

type Table struct {
	TableName   string           `db:"table_name" json:"TableName"`
	Columns     []TableColumn    `db:"columns" json:"Columns"`
	Constraints []ConstraintInfo `db:"constraints" json:"Constraints"`
	Indexes     []IndexInfo      `db:"indexes" json:"Indexes"`
}

type RenameRelation struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

const (
	// Query to get all tables and their columns with detailed information
	getTablesAndColumnsQuery = `
		SELECT 
			t.table_name AS table_name,
			c.column_name AS column_name,
			c.data_type AS data_type,
			c.is_nullable = 'YES' AS is_nullable,
			c.column_default AS column_default,
			c.character_maximum_length AS character_maximum_length,
			c.numeric_precision AS numeric_precision,
			c.numeric_scale AS numeric_scale,
			c.ordinal_position AS ordinal_position
		FROM 
			information_schema.tables t
		JOIN 
			information_schema.columns c ON t.table_name = c.table_name 
			AND t.table_schema = c.table_schema
		WHERE 
			t.table_schema = 'public' 
			AND t.table_type = 'BASE TABLE'
		ORDER BY 
			t.table_name, c.ordinal_position;`

	// Query to get all constraints (PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK)
	getConstraintsQuery = `
		SELECT 
			tc.table_name AS table_name,
			tc.constraint_name AS constraint_name,
			tc.constraint_type AS constraint_type,
			kcu.column_name AS column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			cc.check_clause AS check_clause,
			kcu.ordinal_position AS ordinal_position
		FROM 
			information_schema.table_constraints tc
		LEFT JOIN 
			information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name 
			AND tc.table_schema = kcu.table_schema
		LEFT JOIN 
			information_schema.constraint_column_usage ccu 
			ON tc.constraint_name = ccu.constraint_name 
			AND tc.table_schema = ccu.table_schema
		LEFT JOIN 
			information_schema.check_constraints cc 
			ON tc.constraint_name = cc.constraint_name 
			AND tc.table_schema = cc.constraint_schema
		WHERE 
			tc.table_schema = 'public'
			AND tc.constraint_type IN ('PRIMARY KEY', 'FOREIGN KEY', 'UNIQUE', 'CHECK')
		ORDER BY 
			tc.table_name, tc.constraint_type, kcu.ordinal_position;`

	// Query to get all indexes (excluding those created by constraints)
	getIndexesQuery = `
		SELECT 
			t.relname AS table_name,
			i.relname AS index_name,
			a.attname AS column_name,
			ix.indisunique AS is_unique,
			am.amname AS index_type,
			ix.indisprimary AS is_primary
		FROM 
			pg_class t,
			pg_class i,
			pg_index ix,
			pg_attribute a,
			pg_am am
		WHERE 
			t.oid = ix.indrelid
			AND i.oid = ix.indexrelid
			AND a.attrelid = t.oid
			AND a.attnum = ANY(ix.indkey)
			AND t.relkind = 'r'
			AND am.oid = i.relam
			AND t.relnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public')
			AND NOT ix.indisprimary  -- Exclude primary key indexes (handled by constraints)
			AND NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints tc
				WHERE tc.table_name = t.relname 
				AND tc.constraint_type IN ('UNIQUE', 'FOREIGN KEY')
				AND tc.table_schema = 'public'
			)
		ORDER BY 
			t.relname, i.relname;`

	// Query to get the table schema by name
	getTableSchemaQuery = `
		SELECT 
			t.table_name AS table_name,
			c.column_name AS column_name,
			c.data_type AS data_type,
			c.is_nullable = 'YES' AS is_nullable,
			c.column_default AS column_default,
			c.character_maximum_length AS character_maximum_length,
			c.numeric_precision AS numeric_precision,
			c.numeric_scale AS numeric_scale,
			c.ordinal_position AS ordinal_position
		FROM 
			information_schema.tables t
		JOIN 
			information_schema.columns c ON t.table_name = c.table_name 
			AND t.table_schema = c.table_schema
		WHERE 
			t.table_schema = 'public' 
			AND t.table_type = 'BASE TABLE'
			AND t.table_name = $1
		ORDER BY 
			t.table_name, c.ordinal_position;`
	
	// Query to get constraints for a specific table
	getTableConstraintsQuery = `
		SELECT 
			tc.table_name AS table_name,
			tc.constraint_name AS constraint_name,
			tc.constraint_type AS constraint_type,
			kcu.column_name AS column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			cc.check_clause AS check_clause,
			kcu.ordinal_position AS ordinal_position
		FROM 
			information_schema.table_constraints tc
		LEFT JOIN 
			information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name 
			AND tc.table_schema = kcu.table_schema
		LEFT JOIN 
			information_schema.constraint_column_usage ccu 
			ON tc.constraint_name = ccu.constraint_name 
			AND tc.table_schema = ccu.table_schema
		LEFT JOIN 
			information_schema.check_constraints cc 
			ON tc.constraint_name = cc.constraint_name 
			AND tc.table_schema = cc.constraint_schema
		WHERE 
			tc.table_schema = 'public'
			AND tc.table_name = $1
			AND tc.constraint_type IN ('PRIMARY KEY', 'FOREIGN KEY', 'UNIQUE', 'CHECK')
		ORDER BY 
			tc.table_name, tc.constraint_type, kcu.ordinal_position;`

	// Query to get indexes for a specific table
	getTableIndexesQuery = `
		SELECT 
			t.relname AS table_name,
			i.relname AS index_name,
			a.attname AS column_name,
			ix.indisunique AS is_unique,
			am.amname AS index_type,
			ix.indisprimary AS is_primary
		FROM 
			pg_class t,
			pg_class i,
			pg_index ix,
			pg_attribute a,
			pg_am am
		WHERE 
			t.oid = ix.indrelid
			AND i.oid = ix.indexrelid
			AND a.attrelid = t.oid
			AND a.attnum = ANY(ix.indkey)
			AND t.relkind = 'r'
			AND am.oid = i.relam
			AND t.relnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'public')
			AND t.relname = $1
		ORDER BY 
			t.relname, i.relname;`
)

/*
the schema format for the tables is:

	table_name: {
		columns: [
			{
				column_name: "",
				data_type: "",
				is_nullable: "",
				column_default: "",
				character_maximum_length: "",
				numeric_precision: "",
				numeric_scale: "",
				ordinal_position: "",
			}
		],
		constraints: [
			{
				constraint_name: "",
				constraint_type: "",
			}
		],
		indexes: [
			{
				index_name: "",
				index_type: "",
			}
		]
	}
*/
func GetTables(ctx context.Context, db Querier) (map[string]*Table, error) {
	// Get all tables and columns
	columnsRows, err := db.Query(ctx, getTablesAndColumnsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query table columns: %w", err)
	}
	defer columnsRows.Close()

	var columns []TableColumn
	err = pgxscan.ScanAll(&columns, columnsRows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table columns: %w", err)
	}

	tableColumns := make(map[string][]TableColumn)
	for _, col := range columns {
		tableColumns[col.TableName] = append(tableColumns[col.TableName], col)
	}

	// Get all constraints
	constraints, err := GetConstraints(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get constraints: %w", err)
	}

	tableConstraints := make(map[string][]ConstraintInfo)
	for _, constraint := range constraints {
		tableConstraints[constraint.TableName] = append(tableConstraints[constraint.TableName], constraint)
	}

	// Get all indexes
	indexes, err := GetIndexes(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexes: %w", err)
	}

	tableIndexes := make(map[string][]IndexInfo)
	for _, index := range indexes {
		tableIndexes[index.TableName] = append(tableIndexes[index.TableName], index)
	}

	tablesMap := make(map[string]*Table)
	for tableName, columns := range tableColumns {
		tablesMap[tableName] = &Table{
			TableName:   tableName,
			Columns:     columns,
			Constraints: tableConstraints[tableName],
			Indexes:     tableIndexes[tableName],
		}
	}

	return tablesMap, nil
}

func GetTableSchema(ctx context.Context, tableName string, db Querier) (*Table, error) {
	// Get the table schema
	schemaRows, err := db.Query(ctx, getTableSchemaQuery, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query table schema: %w", err)
	}
	defer schemaRows.Close()

	var columns []TableColumn
	err = pgxscan.ScanAll(&columns, schemaRows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table schema: %w", err)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	return &Table{
		TableName: tableName,
		Columns:   columns,
	}, nil
}

func GetTableConstraints(ctx context.Context, tableName string, db Querier) ([]ConstraintInfo, error) {
	// Get the table constraints
	constraintsRows, err := db.Query(ctx, getTableConstraintsQuery, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query table constraints: %w", err)
	}
	defer constraintsRows.Close()

	var constraints []ConstraintInfo
	err = pgxscan.ScanAll(&constraints, constraintsRows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table constraints: %w", err)
	}

	for i := range len(constraints) {
		if constraints[i].ConstraintType == "CHECK" {
			constraints[i].ColumnName = &strings.Split(*constraints[i].CheckClause, " ")[0]
			if strings.Contains(*constraints[i].CheckClause, "IS NOT NULL") {
				constraints[i].ConstraintType = "NOT NULL"
			}
		}
	}

	return constraints, nil
}

func GetTableIndexes(ctx context.Context, tableName string, db Querier) ([]IndexInfo, error) {
	// Get the table indexes
	indexesRows, err := db.Query(ctx, getTableIndexesQuery, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query table indexes: %w", err)
	}
	defer indexesRows.Close()

	var indexes []IndexInfo
	err = pgxscan.ScanAll(&indexes, indexesRows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table indexes: %w", err)
	}

	return indexes, nil
}

func GetTable(ctx context.Context, tableName string, db Querier) (*Table, error) {
	// Get the table schema
	tableSchema, err := GetTableSchema(ctx, tableName, db)
	if err != nil {
		return nil, err
	}

	// Get the table constraints
	constraints, err := GetTableConstraints(ctx, tableName, db)
	if err != nil {
		return nil, err
	}

	// Get the table indexes
	indexes, err := GetTableIndexes(ctx, tableName, db)
	if err != nil {
		return nil, err
	}

	return &Table{
		TableName:   tableName,
		Columns:     tableSchema.Columns,
		Constraints: constraints,
		Indexes:     indexes,
	}, nil
}

func GetConstraints(ctx context.Context, db Querier) ([]ConstraintInfo, error) {
	// Get all constraints
	constraintsRows, err := db.Query(ctx, getConstraintsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query constraints: %w", err)
	}
	defer constraintsRows.Close()

	var constraints []ConstraintInfo
	err = pgxscan.ScanAll(&constraints, constraintsRows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan constraints: %w", err)
	}

	return constraints, nil
}

func GetIndexes(ctx context.Context, db Querier) ([]IndexInfo, error) {
	// Get all indexes
	indexesRows, err := db.Query(ctx, getIndexesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer indexesRows.Close()

	var indexes []IndexInfo
	err = pgxscan.ScanAll(&indexes, indexesRows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan indexes: %w", err)
	}

	return indexes, nil
}

// ExtractDatabaseSchema extracts the complete database schema as DDL statements
func ExtractDatabaseSchema(ctx context.Context, db Querier) (string, error) {
	var ddlStatements strings.Builder
	ddlStatements.WriteString("-- Database Schema DDL Export\n")
	ddlStatements.WriteString("-- Generated automatically\n\n")

	// Get all tables and columns
	columnsRows, err := db.Query(ctx, getTablesAndColumnsQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query table columns: %w", err)
	}
	defer columnsRows.Close()

	var columns []TableColumn
	err = pgxscan.ScanAll(&columns, columnsRows)
	if err != nil {
		return "", fmt.Errorf("failed to scan table columns: %w", err)
	}

	// Get all constraints
	constraintsRows, err := db.Query(ctx, getConstraintsQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query constraints: %w", err)
	}
	defer constraintsRows.Close()

	var constraints []ConstraintInfo
	err = pgxscan.ScanAll(&constraints, constraintsRows)
	if err != nil {
		return "", fmt.Errorf("failed to scan constraints: %w", err)
	}

	// Get all indexes
	indexesRows, err := db.Query(ctx, getIndexesQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query indexes: %w", err)
	}
	defer indexesRows.Close()

	var indexes []IndexInfo
	err = pgxscan.ScanAll(&indexes, indexesRows)
	if err != nil {
		return "", fmt.Errorf("failed to scan indexes: %w", err)
	}

	// Group data by table
	tableColumns := make(map[string][]TableColumn)
	tableConstraints := make(map[string][]ConstraintInfo)
	tableIndexes := make(map[string][]IndexInfo)

	for _, col := range columns {
		tableColumns[col.TableName] = append(tableColumns[col.TableName], col)
	}

	for _, constraint := range constraints {
		tableConstraints[constraint.TableName] = append(tableConstraints[constraint.TableName], constraint)
	}

	for _, index := range indexes {
		tableIndexes[index.TableName] = append(tableIndexes[index.TableName], index)
	}

	// Generate CREATE TABLE statements
	for tableName, cols := range tableColumns {
		ddlStatements.WriteString(generateCreateTableStatement(tableName, cols, tableConstraints[tableName]))
		ddlStatements.WriteString("\n")

		// Add indexes for this table
		if idxs, exists := tableIndexes[tableName]; exists {
			for _, index := range idxs {
				ddlStatements.WriteString(generateCreateIndexStatement(&index))
				ddlStatements.WriteString("\n")
			}
		}
		ddlStatements.WriteString("\n")
	}

	return ddlStatements.String(), nil
}

func GenerateCreateTableDDL(table *Table) (string, error) {
	var ddlStatements strings.Builder
	ddlStatements.WriteString("-- Database Schema DDL Export\n")
	ddlStatements.WriteString("-- Generated automatically\n\n")
	ddlStatements.WriteString(generateCreateTableStatement(table.TableName, table.Columns, table.Constraints))
	ddlStatements.WriteString("\n")
	// Add indexes for this table
	ddlStatements.WriteString("\n-- Indexes\n")
	for _, index := range table.Indexes {
		ddlStatements.WriteString(generateCreateIndexStatement(&index))
		ddlStatements.WriteString("\n")
	}
	return ddlStatements.String(), nil
}

// generateCreateTableStatement creates a CREATE TABLE DDL statement
func generateCreateTableStatement(tableName string, columns []TableColumn, constraints []ConstraintInfo) string {
	var stmt strings.Builder
	stmt.WriteString(fmt.Sprintf("CREATE TABLE \"%s\" (\n", tableName))

	// Add columns
	columnDefs := make([]string, 0, len(columns))
	for _, col := range columns {
		columnDef := fmt.Sprintf("    \"%s\" %s", col.ColumnName, formatDataType(col))

		if !col.IsNullable {
			columnDef += " NOT NULL"
		}

		if col.ColumnDefault != nil {
			columnDef += fmt.Sprintf(" DEFAULT %s", *col.ColumnDefault)
		}

		columnDefs = append(columnDefs, columnDef)
	}

	// Group constraints by type
	primaryKeys := make([]string, 0)
	uniqueConstraints := make(map[string][]string)
	foreignKeys := make([]ConstraintInfo, 0)
	checkConstraints := make([]ConstraintInfo, 0)

	for _, constraint := range constraints {
		switch constraint.ConstraintType {
		case "PRIMARY KEY":
			if constraint.ColumnName != nil {
				primaryKeys = append(primaryKeys, *constraint.ColumnName)
			}
		case "UNIQUE":
			if constraint.ColumnName != nil {
				uniqueConstraints[constraint.ConstraintName] = append(uniqueConstraints[constraint.ConstraintName], *constraint.ColumnName)
			}
		case "FOREIGN KEY":
			foreignKeys = append(foreignKeys, constraint)
		case "CHECK":
			checkConstraints = append(checkConstraints, constraint)
		}
	}

	// Add PRIMARY KEY constraint
	if len(primaryKeys) > 0 {
		columnDefs = append(columnDefs, fmt.Sprintf("    PRIMARY KEY (\"%s\")", strings.Join(primaryKeys, "\", \"")))
	}

	// Add UNIQUE constraints
	for constraintName, cols := range uniqueConstraints {
		columnDefs = append(columnDefs, fmt.Sprintf("    CONSTRAINT \"%s\" UNIQUE (\"%s\")", constraintName, strings.Join(cols, "\", \"")))
	}

	// Add FOREIGN KEY constraints
	for _, fk := range foreignKeys {
		if fk.ColumnName != nil && fk.ForeignTableName != nil && fk.ForeignColumnName != nil {
			columnDefs = append(columnDefs, fmt.Sprintf("    CONSTRAINT \"%s\" FOREIGN KEY (\"%s\") REFERENCES \"%s\" (\"%s\")",
				fk.ConstraintName, *fk.ColumnName, *fk.ForeignTableName, *fk.ForeignColumnName))
		}
	}

	// Add CHECK constraints
	for _, check := range checkConstraints {
		if check.CheckClause != nil {
			columnDefs = append(columnDefs, fmt.Sprintf("    CONSTRAINT \"%s\" CHECK (%s)", check.ConstraintName, *check.CheckClause))
		}
	}

	stmt.WriteString(strings.Join(columnDefs, ",\n"))
	stmt.WriteString("\n);")

	return stmt.String()
}

// generateCreateIndexStatement creates a CREATE INDEX DDL statement
func generateCreateIndexStatement(index *IndexInfo) string {
	indexType := ""
	if index.IsUnique {
		indexType = "UNIQUE "
	}

	return fmt.Sprintf("CREATE %sINDEX \"%s\" ON \"%s\" USING %s (\"%s\");",
		indexType, index.IndexName, index.TableName, index.IndexType, index.ColumnName)
}

// formatDataType formats the PostgreSQL data type with precision/scale if applicable
func formatDataType(col TableColumn) string {
	dataType := strings.ToUpper(col.DataType)

	switch dataType {
	case "CHARACTER VARYING", "VARCHAR":
		if col.CharacterMaximumLength != nil {
			return fmt.Sprintf("VARCHAR(%d)", *col.CharacterMaximumLength)
		}
		return "VARCHAR"
	case "CHARACTER", "CHAR":
		if col.CharacterMaximumLength != nil {
			return fmt.Sprintf("CHAR(%d)", *col.CharacterMaximumLength)
		}
		return "CHAR"
	case "NUMERIC", "DECIMAL":
		if col.NumericPrecision != nil && col.NumericScale != nil {
			return fmt.Sprintf("NUMERIC(%d,%d)", *col.NumericPrecision, *col.NumericScale)
		} else if col.NumericPrecision != nil {
			return fmt.Sprintf("NUMERIC(%d)", *col.NumericPrecision)
		}
		return "NUMERIC"
	default:
		return dataType
	}
}

func formatConstraint(constraint *ConstraintInfo) string {
	if constraint == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CONSTRAINT \"%s\"", constraint.ConstraintName))
	switch constraint.ConstraintType {
	case "PRIMARY KEY":
		sb.WriteString(fmt.Sprintf(" PRIMARY KEY (\"%s\")", *constraint.ColumnName))
	case "UNIQUE":
		sb.WriteString(fmt.Sprintf(" UNIQUE (\"%s\")", *constraint.ColumnName))
	case "FOREIGN KEY":
		sb.WriteString(fmt.Sprintf(" FOREIGN KEY (\"%s\") REFERENCES \"%s\" (\"%s\")",
			*constraint.ColumnName, *constraint.ForeignTableName, *constraint.ForeignColumnName))
	case "CHECK":
		sb.WriteString(fmt.Sprintf(" CHECK (%s)", *constraint.CheckClause))
	}
	return sb.String()
}

// function that compares two tables schema and returns the DDL statements to update the schema to turn old -> new
func CompareTableSchemas(oldTable, newTable *Table, renames []RenameRelation) (string, error) {
	// each of the three aspects of the schema (columns, constraints, indexes) will be compared separately
	// the result will be three actions: add, drop and alter
	// and a spiecial case for the renaming of the table itself
	var ddlStatements strings.Builder
	ddlStatements.WriteString(fmt.Sprintf("-- Comparing table schema for %s\n", oldTable.TableName))
	ddlStatements.WriteString("-- Generated automatically\n\n")
	// Compare names
	if oldTable.TableName != newTable.TableName {
		ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\";\n",
			oldTable.TableName, newTable.TableName))
	}
	// Compare columns
	oldColumns := make(map[string]*TableColumn)
	for _, col := range oldTable.Columns {
		oldColumns[col.ColumnName] = &col
	}
	newColumns := make(map[string]*TableColumn)
	for _, col := range newTable.Columns {
		newColumns[col.ColumnName] = &col
	}
	// renames
	renamedColumns := make(map[string]string)
	newRenamedColumns := make(map[string]bool)
	for _, rename := range renames {
		renamedColumns[rename.OldName] = rename.NewName
		newRenamedColumns[rename.NewName] = true
	}

	// drop columns that are in old but not in new onl
	for colName, oldCol := range oldColumns {
		// check for renames first
		if newColName, exists := renamedColumns[colName]; exists {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\";\n",
				newTable.TableName, colName, newColName))
			continue
		}

		// if the column is not in the new table, drop it
		if _, exists := newColumns[colName]; !exists {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\";\n",
				newTable.TableName, oldCol.ColumnName))
			continue
		}
	}
	// add columns that are in new but not in old
	for colName, newCol := range newColumns {
		if _, exists := newRenamedColumns[colName]; exists {
			continue // skip renamed columns
		}
		// check if the column is in the old table and add it if not
		if _, exists := oldColumns[colName]; !exists {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN \"%s\" %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
			continue
		}

		// check for changes in column properties
		oldColumn := oldColumns[colName]
		// type changes
		if oldColumn.DataType != newCol.DataType {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		}
		// nullability changes
		if oldColumn.IsNullable != newCol.IsNullable {
			if newCol.IsNullable {
				ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" DROP NOT NULL;\n",
					newTable.TableName, colName))
			} else {
				ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" SET NOT NULL;\n",
					newTable.TableName, colName))
			}
		}
		// default value changes
		if oldColumn.ColumnDefault != nil && newCol.ColumnDefault == nil {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" DROP DEFAULT;\n",
				newTable.TableName, colName))
		} else if oldColumn.ColumnDefault == nil && newCol.ColumnDefault != nil {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" SET DEFAULT %s;\n",
				newTable.TableName, colName, *newCol.ColumnDefault))
		} else if oldColumn.ColumnDefault != nil && newCol.ColumnDefault != nil && *oldColumn.ColumnDefault != *newCol.ColumnDefault {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" SET DEFAULT %s;\n",
				newTable.TableName, colName, *newCol.ColumnDefault))
		}
		// character maximum length changes
		if oldColumn.CharacterMaximumLength != nil && newCol.CharacterMaximumLength == nil {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		} else if oldColumn.CharacterMaximumLength == nil && newCol.CharacterMaximumLength != nil {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		} else if oldColumn.CharacterMaximumLength != nil && newCol.CharacterMaximumLength != nil &&
			*oldColumn.CharacterMaximumLength != *newCol.CharacterMaximumLength {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		}

		// numeric precision changes
		if oldColumn.NumericPrecision != nil && newCol.NumericPrecision == nil {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		} else if oldColumn.NumericPrecision == nil && newCol.NumericPrecision != nil {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		} else if oldColumn.NumericPrecision != nil && newCol.NumericPrecision != nil &&
			*oldColumn.NumericPrecision != *newCol.NumericPrecision {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		}

		// numeric scale changes
		if oldColumn.NumericScale != nil && newCol.NumericScale == nil {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		} else if oldColumn.NumericScale == nil && newCol.NumericScale != nil {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		} else if oldColumn.NumericScale != nil && newCol.NumericScale != nil &&
			*oldColumn.NumericScale != *newCol.NumericScale {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE %s;\n",
				newTable.TableName, colName, formatDataType(*newCol)))
		}
	}

	// Compare constraints
	oldConstraints := make(map[string]*ConstraintInfo)
	for _, constraint := range oldTable.Constraints {
		oldConstraints[constraint.ConstraintName] = &constraint
	}
	newConstraints := make(map[string]*ConstraintInfo)
	for _, constraint := range newTable.Constraints {
		newConstraints[constraint.ConstraintName] = &constraint
	}

	// drop constraints that are in old but not in new add executed with IF EXISTS to avoid errors
	for constraintName, oldConstraint := range oldConstraints {
		if _, exists := newConstraints[constraintName]; !exists {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" DROP CONSTRAINT IF EXISTS \"%s\";\n",
				newTable.TableName, oldConstraint.ConstraintName))
			continue
		}
	}

	// add constraints that are in new but not in old 
	for constraintName, newConstraint := range newConstraints {
		if _, exists := oldConstraints[constraintName]; !exists {
			ddlStatements.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ADD %s;\n",
				newTable.TableName, formatConstraint(newConstraint)))
		}
	}

	// compare indexes
	oldIndexes := make(map[string]*IndexInfo)
	for _, index := range oldTable.Indexes {
		oldIndexes[index.IndexName] = &index
	}
	newIndexes := make(map[string]*IndexInfo)
	for _, index := range newTable.Indexes {
		newIndexes[index.IndexName] = &index
	}
	// drop indexes that are in old but not in new
	for indexName, oldIndex := range oldIndexes {
		if _, exists := newIndexes[indexName]; !exists {
			ddlStatements.WriteString(fmt.Sprintf("DROP INDEX IF EXISTS \"%s\";\n", oldIndex.IndexName))
			continue
		}
	}

	// add indexes that are in new but not in old
	for indexName, newIndex := range newIndexes {
		if _, exists := oldIndexes[indexName]; !exists {
			ddlStatements.WriteString(generateCreateIndexStatement(newIndex))
		}
	}

	// return the DDL statements
	if ddlStatements.Len() == 0 {
		return "", fmt.Errorf("no changes detected between old and new table schemas")
	}
	return ddlStatements.String(), nil
}