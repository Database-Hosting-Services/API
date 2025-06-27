package tables

import (
	"DBHS/config"
	"DBHS/utils"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func GetAllTablesNameOid(ctx context.Context, projectId int64, db pgxscan.Querier) ([]ShortTable, error) {
	var tables []ShortTable
	err := pgxscan.Select(ctx, db, &tables, `SELECT oid, name FROM "Ptable" WHERE project_id = $1`, projectId)
	if err != nil {
		return nil, err
	}

	return tables, err
}

func InsertNewTable(ctx context.Context, table *Table, TableId *int, db utils.Querier) error {
	err := db.QueryRow(ctx, InsertNewTableRecordStmt, table.OID, table.Name, table.Description, table.ProjectID).Scan(TableId)
	if err != nil {
		return fmt.Errorf("failed to insert new table: %w", err)
	}
	return nil
}

func DeleteTableRecord(ctx context.Context, tableId int, db utils.Querier) error {
	_, err := db.Exec(ctx, fmt.Sprintf(DeleteTableStmt, "id"), tableId)
	if err != nil {
		return fmt.Errorf("failed to delete table record: %w", err)
	}
	return nil
}

func CheckOwnershipQuery(ctx context.Context, projectId string, userId int, db utils.Querier) (bool, error) {
	var count int
	err := db.QueryRow(ctx, CheckOwnershipStmt, projectId, userId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}

func ReadTableColumns(ctx context.Context, tableName string, db pgxscan.Querier) (map[string]DbColumn, error) {
	var columns []DbColumn
	err := pgxscan.Select(ctx, db, &columns, ReadTableStmt, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to read table: %w", err)
	}
	columnsMap := make(map[string]DbColumn)
	for _, column := range columns {
		columnsMap[column.Name] = column
	}
	return columnsMap, nil
}

func GetTableName(ctx context.Context, tableOID string, db utils.Querier) (string, error) {
	var tableName string
	err := db.QueryRow(ctx, GetTableNameStmt, tableOID).Scan(&tableName)
	if err != nil {
		return "", fmt.Errorf("failed to get table name: %w", err)
	}
	return tableName, nil
}

func DeleteTableFromHostingServer(ctx context.Context, tableName string, db utils.Querier) error {
	_, err := db.Exec(ctx, fmt.Sprintf(DropTableStmt, tableName))
	if err != nil {
		return fmt.Errorf("failed to delete table from hosting server: %w", err)
	}
	return nil
}

func DeleteTableFromServerDb(ctx context.Context, tableOID string, db utils.Querier) error {
	_, err := db.Exec(ctx, fmt.Sprintf(DeleteTableStmt, "oid"), tableOID)
	if err != nil {
		return fmt.Errorf("failed to delete table from server DB: %w", err)
	}
	return nil
}

/*
	the query has x parts
	SELECT * FROM [TABLE_NAME]
	WHERE [FILTERS]
	ORDER BY [ORDERED BY]
	LIMIT [LIMIT]
	OFFSET [PAGE * LIMIT]
*/

func ReadTableData(ctx context.Context, tableName string, parameters map[string][]string, db utils.Querier) (*Data, error) {
	query, err := PrepareQuery(tableName, parameters)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()
	if columns == nil {
		return nil, err
	}
	data := Data{
		Columns: make([]ShowColumn, len(columns)),
	}
	// get column name, type
	for i, col := range columns {
		data.Columns[i] = ShowColumn{
			Name: col.Name,
			Type: config.PgTypes[col.DataTypeOID],
		}
	}

	//reserve memory where a row will be read
	values := make([]interface{}, len(columns))
	ptr := make([]interface{}, len(columns))
	for i := range values {
		ptr[i] = &values[i]
	}

	row := make(map[string]interface{})
	for rows.Next() {
		if err := rows.Scan(ptr...); err != nil {
			return nil, err
		}
		for i := range columns {
			row[columns[i].Name] = values[i]
		}
		data.Rows = append(data.Rows, row)
	}

	return &data, nil
}

func PrepareQuery(tableName string, parameters map[string][]string) (string, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	query, err := AddFilters(query, parameters["filter"])
	if err != nil {
		return "", err
	}

	query, err = AddOrder(query, parameters["order"])
	if err != nil {
		return "", err
	}
	// Add Limit and Offset
	limit, err := strconv.Atoi(parameters["limit"][0])
	if err != nil {
		return "", err
	}

	page, err := strconv.Atoi(parameters["page"][0])
	if err != nil {
		return "", err
	}

	query = query + fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, page*limit)

	return query, nil
}

// filter will be a string in the format "column:op:value"
func AddFilters(query string, filters []string) (string, error) {
	if filters == nil || len(filters) == 0 {
		return query, nil
	}
	query = query + " WHERE "
	var opMap = map[string]string{
		"eq":   "=",
		"neq":  "!=",
		"lt":   "<",
		"lte":  "<=",
		"gt":   ">",
		"gte":  ">=",
		"like": "LIKE", // if needed
	}

	predicates := make([]string, 0, len(filters))
	for _, filter := range filters {
		parts := strings.Split(filter, ":")
		column, op, value := parts[0], parts[1], parts[2]
		if op == "like" {
			predicates = append(predicates, fmt.Sprintf("%s %s %s", column, opMap[op], value))
		} else {
			intV, err := strconv.Atoi(value)
			if err != nil {
				return "", err
			}
			predicates = append(predicates, fmt.Sprintf("%s %s %d", column, opMap[op], intV))
		}
	}

	return query + strings.Join(predicates, " AND "), nil

}

// order will be a string in the format "column:op"
func AddOrder(query string, orders []string) (string, error) {
	if orders == nil || len(orders) == 0 {
		return query, nil
	}

	query = query + " ORDER BY "
	var opMap = map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	predicates := make([]string, 0, len(orders))
	for _, order := range orders {
		parts := strings.Split(order, ":")
		column, op := parts[0], parts[1]
		predicates = append(predicates, fmt.Sprintf("%s %s", column, opMap[op]))
	}

	return query + strings.Join(predicates, ", "), nil
}
