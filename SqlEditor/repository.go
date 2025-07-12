package sqleditor

import (
	api "DBHS/utils/apiError"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pingcap/tidb/parser"
	// "github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

// validateSQLWithParser validates SQL syntax using TiDB SQL parser
func validateSQLWithParser(query string) error {
	p := parser.New()

	// Parse the SQL statement
	stmts, _, err := p.Parse(query, "", "")
	if err != nil {
		return fmt.Errorf("SQL syntax error: %v", err)
	}

	// Check if we have any statements
	if len(stmts) == 0 {
		return errors.New("empty SQL statement")
	}

	return nil
}

func FetchQueryData(ctx context.Context, conn *pgxpool.Pool, query string) (ResponseBody, api.ApiError) {
	// Validate SQL syntax using parser
	if err := validateSQLWithParser(query); err != nil {
		return ResponseBody{}, *api.NewApiError("Invalid SQL syntax", 400, err)
	}

	wrappedQuery := WrapQueryWithJSONAgg(query)

	var result string
	startTime := time.Now()
	err := conn.QueryRow(ctx, wrappedQuery).Scan(&result)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return ResponseBody{}, *api.NewApiError("Table/Column not found", 404, errors.New("Table or column does not exist: "+err.Error()))
		}
		if strings.Contains(err.Error(), "cannot scan NULL into *string") {
			return ResponseBody{}, *api.NewApiError("There is no data", 200, errors.New("No data found for the query"))
		}
		return ResponseBody{}, *api.NewApiError("Query execution failed", 500, errors.New("failed to execute query: "+err.Error()))
	}

	// Extract column names from the JSON result
	columnNames, err := ExtractColumnNames(result)

	if err != nil {
		return ResponseBody{}, *api.NewApiError("Internal server error", 500, errors.New("failed to extract column names: "+err.Error()))
	}

	executionTime := time.Since(startTime)
	return ResponseBody{
		Result:        json.RawMessage(result),
		ColumnNames:   columnNames,
		ExecutionTime: float64(executionTime.Nanoseconds()) / 1e6, // Convert to milliseconds as float64
	}, api.ApiError{} // Return empty ApiError to indicate success
}
