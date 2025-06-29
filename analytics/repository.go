package analytics

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// --------------------- Functions to get last records -----------------------

func GetLastExecutionTimeRecord(ctx context.Context, db *pgxpool.Pool, projectId int64) (*DatabaseActivityWithDates, error) {
	var record DatabaseActivityWithDates
	err := db.QueryRow(ctx, GET_LAST_EXECUTION_TIME_STATS, projectId).Scan(
		&record.Timestamp, &record.TotalTimeMs, &record.TotalQueries,
	)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func GetLastDatabaseUsageRecord(ctx context.Context, db *pgxpool.Pool, projectId int64) (*DatabaseUsageCostWithDates, error) {
	var record DatabaseUsageCostWithDates
	err := db.QueryRow(ctx, GET_LAST_DATABASE_USAGE_STATS, projectId).Scan(
		&record.Timestamp, &record.ReadWriteCost, &record.CPUCost, &record.TotalCost,
	)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
