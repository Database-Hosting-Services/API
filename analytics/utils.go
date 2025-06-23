package analytics

import (
	"DBHS/indexes"
	api "DBHS/utils/apiError"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CalculateCosts calculates the costs associated with database usage based on read/write queries and CPU time.
func (d *DatabaseUsageStats) CalculateCosts() DatabaseUsageCost {
	Cost := DatabaseUsageCost{
		ReadWriteCost: float64(d.ReadQueries)/1_000_000*1.00 + float64(d.WriteQueries)/1_000_000*1.50,
		CPUCost:       (d.TotalCPUTimeMs / 1000 / 3600) * 0.000463,
	}
	Cost.TotalCost = Cost.ReadWriteCost + Cost.CPUCost
	return Cost
}

func GetConnectionToAnalyticsPool(ctx context.Context, db *pgxpool.Pool, projectOid string) (*pgxpool.Pool, api.ApiError) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return nil, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// Get the project pool connection
	conn, err := indexes.ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return nil, *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return nil, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}

	// Check if pg_stat_statements extension exists
	var extensionExists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_stat_statements')").Scan(&extensionExists)
	if err != nil {
		return nil, *api.NewApiError("Internal server error", 500, errors.New("failed to check for pg_stat_statements extension: "+err.Error()))
	}

	// If extension doesn't exist, try to create it
	if !extensionExists {
		_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS pg_stat_statements")
		if err != nil {
			return nil, *api.NewApiError("Internal server error", 500,
				errors.New("pg_stat_statements extension is not available. Please ensure it is installed in PostgreSQL and included in shared_preload_libraries"))
		}
	}

	return conn, api.ApiError{} // Return empty ApiError to indicate success
}
