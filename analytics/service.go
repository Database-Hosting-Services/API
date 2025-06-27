package analytics

import (
	"DBHS/config"
	"DBHS/indexes"
	"DBHS/projects"

	api "DBHS/utils/apiError"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// --------------------- Background worker service functions -----------------------

func GetDatabaseStorage(ctx context.Context, db *pgxpool.Pool, projectOid string) (Storage, api.ApiError) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int64)
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
	conn, err := GetConnectionToAnalyticsPool(ctx, db, projectOid)
	defer conn.Close()

	if err.Error() != nil {
		return ExecutionTimeStats{}, err
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
	conn, err := GetConnectionToAnalyticsPool(ctx, db, projectOid)
	defer conn.Close()

	if err.Error() != nil {
		return DatabaseUsageCost{}, err
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

// --------------------- HTTP endpoint service functions -----------------------

func GetALLDatabaseStorage(ctx context.Context, db *pgxpool.Pool, projectOid string) ([]StorageWithDates, api.ApiError) {
	owner_id, ok := ctx.Value("user-id").(int64)
	if !ok || owner_id == 0 {
		return nil, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	id, err := projects.GetProjectID(ctx, db, owner_id, projectOid)
	if err != nil {
		if errors.Is(err, projects.ErrorProjectNotFound) {
			return nil, *api.NewApiError("Project not found", 404, err)
		}
		return nil, *api.NewApiError("Internal server error", 500, err)
	}

	// Get all storage records for the project
	rows, err := config.DB.Query(ctx, GET_ALL_CURRENT_STORAGE, id)
	if err != nil {
		return nil, *api.NewApiError("Internal server error", 500, errors.New("failed to retrieve storage records: "+err.Error()))
	}

	defer rows.Close()
	var storageRecords []StorageWithDates
	for rows.Next() {
		var storage StorageWithDates
		if err := rows.Scan(&storage.Timestamp, &storage.ManagementStorage, &storage.ActualData); err != nil {
			return nil, *api.NewApiError("Internal server error", 500, errors.New("failed to scan storage record: "+err.Error()))
		}
		storageRecords = append(storageRecords, storage)
	}

	if err := rows.Err(); err != nil {
		return nil, *api.NewApiError("Internal server error", 500, errors.New("error iterating over storage records: "+err.Error()))
	}

	return storageRecords, api.ApiError{} // Return empty ApiError to indicate success
}
