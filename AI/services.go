package ai

import (
	"DBHS/config"
	"DBHS/tables"
	"context"
	"encoding/json"
	"time"

	"github.com/Database-Hosting-Services/AI-Agent/RAG"
)

func getReport(projectUUID string, userID int64, analytics Analytics, AI RAG.RAGmodel) (string, error) {
	// get project name and connection
	_, userDb, err := tables.ExtractDb(context.Background(), projectUUID, userID, config.DB)
	if err != nil {
		return "", err
	}

	// get database schema
	databaseSchema, err := ExtractDatabaseSchema(context.Background(), userDb)
	if err != nil {
		return "", err
	}

	// convert analytics to string
	analyticsString, err := json.Marshal(analytics)
	if err != nil {
		return "", err
	}

	// get report
	report, err := AI.Report(databaseSchema, string(analyticsString))
	if err != nil {
		return "", err
	}

	return report, nil
}

func AgentQuery(projectUUID string, userID int64, prompt string, AI RAG.RAGmodel) (*RAG.AgentResponse, error) {
	// get project name and connection
	_, userDb, err := tables.ExtractDb(context.Background(), projectUUID, userID, config.DB)
	if err != nil {
		config.App.ErrorLog.Println("Error extracting database connection:", err)
		return nil, err
	}

	// get database schema
	databaseSchema, err := ExtractDatabaseSchema(context.Background(), userDb)
	if err != nil {
		config.App.ErrorLog.Println("Error extracting database schema:", err)
		return nil, err
	}

	response, err := AI.QueryAgent("schemas-json", databaseSchema, prompt, 10)
	if err != nil {
		config.App.ErrorLog.Println("Error querying agent:", err)
		return nil, err
	}

	// add the schema changes to the cache
	err = config.VerifyCache.Set("schema-changes:"+projectUUID, response.SchemaDDL, 10*time.Minute)
	if err != nil {
		config.App.ErrorLog.Println("Error adding schema changes to cache:", err)
		return nil, err
	}
	return response, nil
}