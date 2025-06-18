package analytics

import (
	"DBHS/indexes"
	api "DBHS/utils/apiError"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDatabaseStorage(ctx context.Context, db *pgxpool.Pool, projectOid string) (Storage, api.ApiError) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return Storage{}, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := indexes.ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return Storage{}, *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return Storage{}, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}
	defer conn.Close()

	var storage Storage
	if err := conn.QueryRow(ctx, GET_CURRENT_STORAGE).Scan(&storage.ManagementStorage, &storage.ActualData); err != nil {
		return Storage{}, *api.NewApiError("Internal server error", 500, errors.New("failed to retrieve storage information: "+err.Error()))
	}

	return storage, api.ApiError{} // Return empty ApiError to indicate success
}

func GetExecutionTimeStats(ctx context.Context, db *pgxpool.Pool, projectOid string) (ExecutionTimeStats, api.ApiError) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return ExecutionTimeStats{}, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// Get the project pool connection
	conn, err := indexes.ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return ExecutionTimeStats{}, *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return ExecutionTimeStats{}, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}
	defer conn.Close()

	// Check if pg_stat_statements extension exists
	var extensionExists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_stat_statements')").Scan(&extensionExists)
	if err != nil {
		return ExecutionTimeStats{}, *api.NewApiError("Internal server error", 500, errors.New("failed to check for pg_stat_statements extension: "+err.Error()))
	}

	// If extension doesn't exist, try to create it
	if !extensionExists {
		_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS pg_stat_statements")
		if err != nil {
			return ExecutionTimeStats{}, *api.NewApiError("Internal server error", 500,
				errors.New("pg_stat_statements extension is not available. Please ensure it is installed in PostgreSQL and included in shared_preload_libraries"))
		}
	}

	// Get current database name
	var dbName string
	if err := conn.QueryRow(ctx, "SELECT current_database()").Scan(&dbName); err != nil {
		return ExecutionTimeStats{}, *api.NewApiError("Internal server error", 500, errors.New("failed to get database name: "+err.Error()))
	}

	var stats ExecutionTimeStats
	if err := conn.QueryRow(ctx, GET_MAX_AVG_TOTAL_EXECUTION_TIME, dbName).Scan(
		&stats.TotalTimeMs,
		&stats.MaxTimeMs,
		&stats.AvgTimeMs,
	); err != nil {
		return ExecutionTimeStats{}, *api.NewApiError("Internal server error", 500, errors.New("failed to retrieve execution time stats: "+err.Error()))
	}

	return stats, api.ApiError{} // Return empty ApiError to indicate success
}

func GetDatabaseUsageStats(ctx context.Context, db *pgxpool.Pool, projectOid string) (DatabaseUsageCost, api.ApiError) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return DatabaseUsageCost{}, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// Get the project pool connection
	conn, err := indexes.ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return DatabaseUsageCost{}, *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return DatabaseUsageCost{}, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}
	defer conn.Close()

	// Check if pg_stat_statements extension exists
	var extensionExists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_stat_statements')").Scan(&extensionExists)
	if err != nil {
		return DatabaseUsageCost{}, *api.NewApiError("Internal server error", 500, errors.New("failed to check for pg_stat_statements extension: "+err.Error()))
	}

	// If extension doesn't exist, try to create it
	if !extensionExists {
		_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS pg_stat_statements")
		if err != nil {
			return DatabaseUsageCost{}, *api.NewApiError("Internal server error", 500,
				errors.New("pg_stat_statements extension is not available. Please ensure it is installed in PostgreSQL and included in shared_preload_libraries"))
		}
	}

	var stats DatabaseUsageStats
	if err := conn.QueryRow(ctx, GET_READ_WRITE_CPU).Scan(
		&stats.ReadQueries,
		&stats.WriteQueries,
		&stats.TotalCPUTimeMs,
	); err != nil {
		return DatabaseUsageCost{}, *api.NewApiError("Internal server error", 500, errors.New("failed to retrieve database usage stats: "+err.Error()))
	}

	// Calculate costs based on the retrieved stats
	Cost := stats.CalculateCosts()

	return Cost, api.ApiError{} // Return empty ApiError to indicate success
}
