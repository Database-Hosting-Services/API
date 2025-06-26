package workers

import (
	"DBHS/analytics"
	"DBHS/config"
	"context"
)

const (
	STORAGE_ANALYTICS = "Storage"
	EXECUTION_TIME_ANALYTICS = "ExecutionTimeStats"
	DATABASE_USAGE_ANALYTICS = "DatabaseUsageStats"
	DATABASE_USAGE_COST = "DatabaseUsageCost"
)

const (
	GET_ALL_PROJECTS = `
		SELECT id, oid, owner_id FROM projects
	`
	INSERT_ANALYTICS = `
		INSERT INTO analytics ("projectId", "type", "data") VALUES ($1, $2, $3)
	`
)

func GatherAnalytics(app *config.Application) {
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
		context := context.WithValue(context.Background(), "user-id", ownerId)
		// get the storage of the project
		storage, apiErr := analytics.GetDatabaseStorage(context, config.DB, projectOid)
		if apiErr.Error() != nil {
			app.ErrorLog.Println(apiErr.Error())
			continue
		}
		// get the execution time stats of the project
		executionTimeStats, apiErr := analytics.GetExecutionTimeStats(context, config.DB, projectOid)
		if apiErr.Error() != nil {
			app.ErrorLog.Println(apiErr.Error())
			continue
		}
		// get the database usage stats of the project
		databaseUsageStats, apiErr := analytics.GetDatabaseUsageStats(context, config.DB, projectOid)
		if apiErr.Error() != nil {
			app.ErrorLog.Println(apiErr.Error())
			continue
		}

		// insert the analytics to the database
		_, err = config.DB.Exec(context, INSERT_ANALYTICS, projectOid, STORAGE_ANALYTICS, storage)
		if err != nil {
			app.ErrorLog.Println(err)
			continue
		}

		_, err = config.DB.Exec(context, INSERT_ANALYTICS, projectOid, EXECUTION_TIME_ANALYTICS, executionTimeStats)
		if err != nil {
			app.ErrorLog.Println(err)
			continue
		}

		_, err = config.DB.Exec(context, INSERT_ANALYTICS, projectOid, DATABASE_USAGE_ANALYTICS, databaseUsageStats)
		if err != nil {
			app.ErrorLog.Println(err)
			continue
		}

		// log the analytics
		app.InfoLog.Println("Analytics gathered for project", projectOid, "✅")
	}

	app.InfoLog.Println("Analytics gathered for all projects ✅")
}