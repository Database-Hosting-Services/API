package sqleditor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WrapQueryWithJSONAgg takes a SQL query and wraps it with json_agg and row_to_json
// to return the result as a JSON array
func WrapQueryWithJSONAgg(sqlQuery string) string {
	// Clean the input query by removing leading/trailing whitespace and semicolon
	cleanQuery := strings.TrimSuffix(strings.TrimSpace(sqlQuery), ";")
	// Build the wrapped query
	wrappedQuery := fmt.Sprintf(`SELECT json_agg(row_to_json(t))::text AS result FROM (%s) t;`, cleanQuery)
	return wrappedQuery
}

// extractColumnNames extracts column names from the JSON result
func ExtractColumnNames(jsonResult string) ([]string, error) {
	var data []map[string]interface{}

	err := json.Unmarshal([]byte(jsonResult), &data)
	if err != nil {
		return nil, err
	}

	// If no data, return empty slice
	if len(data) == 0 {
		return []string{}, nil
	}

	// Extract column names from the first row
	var columnNames []string
	for key := range data[0] {
		columnNames = append(columnNames, key)
	}

	return columnNames, nil
}
