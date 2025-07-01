package sqleditor

import (
	api "DBHS/utils/apiError"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func FetchQueryData(ctx context.Context, conn *pgxpool.Pool, query string) (ResponseBody, api.ApiError) {
	startTime := time.Now()
	query = WrapQueryWithJSONAgg(query)

	var result string
	err := conn.QueryRow(ctx, query).Scan(&result)
	if err != nil {
		return ResponseBody{}, *api.NewApiError("Internal server error", 500, errors.New("failed to execute query: "+err.Error()))
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
