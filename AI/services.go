package ai

import (
	"DBHS/config"
	"DBHS/utils"
	"context"
	"encoding/json"
	"github.com/Database-Hosting-Services/AI-Agent/RAG"
)

func getReport(projectUUID string, userID int, analytics Analytics, AI RAG.RAGmodel) (string, error) {
	// get project name and connection
	_, userDb, err := utils.ExtractDb(context.Background(), projectUUID, userID, config.DB)
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
