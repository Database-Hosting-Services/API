package ai

import (
	"DBHS/config"
	"DBHS/tables"
	"DBHS/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/Database-Hosting-Services/AI-Agent/RAG"
)

func getReport(projectUUID string, userID int, analytics Analytics, AI RAG.RAGmodel) (string, error) {
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

func SaveChatAction(ctx context.Context, db utils.Querier, chatId, userID int, question string, answer string) error {
	// here i save the user prompt and the AI response together
	// the chat action is a combination of the user question and the AI answer
	if err := SaveUserChatMessage(ctx, db, chatId, question); err != nil {
		return err
	}
	if err := SaveAIChatMessage(ctx, db, chatId, answer); err != nil {
		return err
	}
	return nil
}

func GetProjectIDfromOID(ctx context.Context, db utils.Querier, projectOID string) (int, error) {
	var projectID int
	err := db.QueryRow(ctx, "SELECT id FROM projects WHERE oid = $1", projectOID).Scan(&projectID)
	if err != nil {
		return 0, err
	}
	return projectID, nil
}

func GetOrCreateChatData(ctx context.Context, db utils.Querier, userID, projectID int) (ChatData, error) {
	// i suppose return the chat history if needed, currently it just returns the chat data

	chat, err := GetUserChatForProject(ctx, db, userID, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			chat, err = CreateNewChat(ctx, db, utils.GenerateOID(), userID, projectID)
			if err != nil {
				return ChatData{}, err
			}
		} else {
			return ChatData{}, err
		}
	}
	return chat, nil
}
