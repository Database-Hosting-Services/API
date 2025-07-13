package sqleditor

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// List of dangerous SQL keywords/operations
var dangerousOperations = []string{
	"CREATE",
	"DROP",
	"ALTER",
	"TRUNCATE",
	"GRANT",
	"REVOKE",
	"EXEC",
	"EXECUTE",
	"CALL",
	"MERGE",
	"REPLACE",
	"RENAME",
	"COMMENT",
}

// Additional check for dangerous patterns
var dangerousPatterns = []string{
	`\bINTO\s+OUTFILE\b`, // File operations
	`\bLOAD_FILE\b`,      // File operations
	`\bSYSTEM\b`,         // System commands
	`\bSHELL\b`,          // Shell commands
}

// System tables and schemas that should never be updated
var protectedTables = []string{
	"PG_DATABASE",
	"PG_CLASS",
	"PG_NAMESPACE",
	"PG_TABLES",
	"PG_ATTRIBUTE",
	"INFORMATION_SCHEMA",
}

// ValidateQuery checks if the SQL query contains dangerous operations
// Returns true if the query is safe (allows SELECT, INSERT, DELETE, and UPDATE), false if it contains dangerous operations
func ValidateQuery(sqlQuery string) (bool, error) {
	// Convert to uppercase for case-insensitive matching
	upperQuery := strings.ToUpper(strings.TrimSpace(sqlQuery))

	// Check for dangerous operations at the beginning of statements
	// Split by semicolon to handle multiple statements
	statements := strings.Split(upperQuery, ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		// Check if statement starts with any dangerous operation
		for _, operation := range dangerousOperations {
			// Use regex to match word boundaries to avoid false positives
			pattern := fmt.Sprintf(`^\s*%s\b`, regexp.QuoteMeta(operation))
			matched, _ := regexp.MatchString(pattern, statement)
			if matched {
				return false, fmt.Errorf("query contains forbidden operation: %s", operation)
			}
		}

		for _, pattern := range dangerousPatterns {
			matched, _ := regexp.MatchString(pattern, statement)
			if matched {
				return false, errors.New("query contains forbidden file or system operations")
			}
		}

		// Check if UPDATE statement targets protected system tables
		if strings.HasPrefix(statement, "UPDATE") {
			for _, protectedTable := range protectedTables {
				// Check if the statement contains references to protected tables
				pattern := fmt.Sprintf(`\b%s`, regexp.QuoteMeta(protectedTable))
				matched, _ := regexp.MatchString(pattern, statement)
				if matched {
					return false, fmt.Errorf("cannot update protected system table/schema: %s", protectedTable)
				}
			}
		}
	}

	return true, nil
}

// WrapQueryWithJSONAgg takes a SQL query and wraps it with json_agg and row_to_json
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
