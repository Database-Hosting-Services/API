package workers

import (
	"DBHS/analytics"
	"DBHS/config"
	"context"
)

const (
	STORAGE_ANALYTICS        = "Storage"
	EXECUTION_TIME_ANALYTICS = "ExecutionTimeStats"
	DATABASE_USAGE_ANALYTICS = "DatabaseUsageStats"
)

const (
	GET_ALL_PROJECTS = `
		SELECT id, oid, owner_id FROM projects
	`
	INSERT_ANALYTICS = `
		INSERT INTO analytics ("projectId", "type", "data") VALUES ($1, $2, $3)
	`
)

func gatherAndInsertStorageAnalytics(app *config.Application, ctx context.Context, projectId int64, projectOid string) error {
	// Get current storage data
	currentStorage, apiErr := analytics.GetDatabaseStorage(ctx, config.DB, projectOid)
	if apiErr.Error() != nil {
		return apiErr.Error()
	}

	_, err := config.DB.Exec(ctx, INSERT_ANALYTICS, projectId, STORAGE_ANALYTICS, currentStorage)
	if err != nil {
		return err
	}

	return nil
}

func gatherAndInsertExecutionTimeAnalytics(app *config.Application, ctx context.Context, projectId int64, projectOid string) error {
	// Get current execution time stats
	currentStats, apiErr := analytics.GetExecutionTimeStats(ctx, config.DB, projectOid)
	if apiErr.Error() != nil {
		return apiErr.Error()
	}

	// Get last execution time record
	lastRecord, err := analytics.GetLastExecutionTimeRecord(ctx, config.DB, projectId)
	if err != nil {
		// If no previous record exists, insert current data directly
		app.InfoLog.Println("No previous execution time record found, inserting current data")
	}

	// Calculate difference or use current data
	finalStats := analytics.CalculateExecutionTimeDifference(currentStats, lastRecord)

	_, err = config.DB.Exec(ctx, INSERT_ANALYTICS, projectId, EXECUTION_TIME_ANALYTICS, finalStats)
	if err != nil {
		return err
	}

	return nil
}

func gatherAndInsertDatabaseUsageAnalytics(app *config.Application, ctx context.Context, projectId int64, projectOid string) error {
	// Get current database usage stats
	currentUsage, apiErr := analytics.GetDatabaseUsageStats(ctx, config.DB, projectOid)
	if apiErr.Error() != nil {
		return apiErr.Error()
	}

	// Get last database usage record
	lastRecord, err := analytics.GetLastDatabaseUsageRecord(ctx, config.DB, projectId)
	if err != nil {
		// If no previous record exists, insert current data directly
		app.InfoLog.Println("No previous database usage record found, inserting current data")
	}

	// Calculate difference or use current data
	finalUsage := analytics.CalculateDatabaseUsageDifference(currentUsage, lastRecord)

	_, err = config.DB.Exec(ctx, INSERT_ANALYTICS, projectId, DATABASE_USAGE_ANALYTICS, finalUsage)
	if err != nil {
		return err
	}

	return nil
}

func GatherAnalytics(app *config.Application) {
	app.InfoLog.Println("Starting analytics gathering for all projects...")

	// get all projects from the database
	projects, err := config.DB.Query(context.Background(), GET_ALL_PROJECTS)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}
	defer projects.Close()
	app.InfoLog.Println("Gathering analytics for all projects 🔍")

	for projects.Next() {
		var id int64
		var projectOid string
		var ownerId int64
		err := projects.Scan(&id, &projectOid, &ownerId)
		if err != nil {
			app.ErrorLog.Println(err)
			continue
		}

		app.InfoLog.Println("Gathering analytics for project", projectOid, "🔍")
		ctx := context.WithValue(context.Background(), "user-id", ownerId)

		// Gather storage analytics
		if err := gatherAndInsertStorageAnalytics(app, ctx, id, projectOid); err != nil {
			app.ErrorLog.Println("Storage analytics error:", err)
			continue
		}

		// Gather execution time analytics
		if err := gatherAndInsertExecutionTimeAnalytics(app, ctx, id, projectOid); err != nil {
			app.ErrorLog.Println("Execution time analytics error:", err)
			continue
		}

		// Gather database usage analytics
		if err := gatherAndInsertDatabaseUsageAnalytics(app, ctx, id, projectOid); err != nil {
			app.ErrorLog.Println("Database usage analytics error:", err)
			continue
		}

		// log the analytics
		app.InfoLog.Println("Analytics gathered for project", projectOid, "✅")
	}

	app.InfoLog.Println("Analytics gathered for all projects ✅")
}
