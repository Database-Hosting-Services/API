package ai

import (
	"DBHS/config"
	"DBHS/utils"
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Database-Hosting-Services/AI-Agent/RAG"
)

func getReport(projectUUID string, userID int64, analytics Analytics, AI RAG.RAGmodel) (string, error) {
	// get project name and connection
	_, userDb, err := utils.ExtractDb(context.Background(), projectUUID, userID, config.DB)
	if err != nil {
		return "", err
	}

	// get database schema
	databaseSchema, err := utils.ExtractDatabaseSchema(context.Background(), userDb)
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

func AgentQuery(projectUUID string, userID int64, prompt string, AI RAG.RAGmodel) (*RAG.AgentResponse, error) {
	// get project name and connection
	_, userDb, err := utils.ExtractDb(context.Background(), projectUUID, userID, config.DB)
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

	// remove the json code segment with all the code from the response
	response.Response = removeJSONCodeSegments(response.Response)

	// add the schema changes to the cache
	err = config.VerifyCache.Set("schema-changes:"+projectUUID, response.SchemaDDL, 10*time.Minute)
	if err != nil {
		config.App.ErrorLog.Println("Error adding schema changes to cache:", err)
		return nil, err
	}
	return response, nil
}

func AgentExec(projectUUID string, userID int64, AI RAG.RAGmodel) error {
	// get project name and connection
	_, userDb, err := utils.ExtractDb(context.Background(), projectUUID, userID, config.DB)
	if err != nil {
		config.App.ErrorLog.Println("Error extracting database connection:", err)
		return err
	}

	// get the DDL from the cache
	ddl, err := config.VerifyCache.Get("schema-changes:"+projectUUID, nil)
	if err != nil {
		config.App.ErrorLog.Println("Error getting schema changes from cache:", err)
		return err
	}

	if ddl == nil {
		config.App.ErrorLog.Println("No schema changes found in cache for project:", projectUUID)
		return errors.New("changes expired or not found")
	}

	// execute the DDL
	tx, err := userDb.Begin(context.Background())
	if err != nil {
		config.App.ErrorLog.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), ddl.(string))
	if err != nil {
		config.App.ErrorLog.Println("Error executing DDL:", err)
		return err
	}
	if err := tx.Commit(context.Background()); err != nil {
		config.App.ErrorLog.Println("Error committing transaction:", err)
		return err
	}
	return nil
}

func ClearCacheForProject(projectUUID string) error {
	// clear the cache for the project
	err := config.VerifyCache.Delete("schema-changes:" + projectUUID)
	if err != nil {
		config.App.ErrorLog.Println("Error clearing cache for project:", projectUUID, "Error:", err)
		return err
	}
	return nil
}

func removeJSONCodeSegments(text string) string {
	// Remove JSON code blocks (```json...``` or ```...``` containing JSON)
	jsonCodeBlockRegex := regexp.MustCompile("(?s)```(?:json)?\\s*\\{.*?\\}```")
	text = jsonCodeBlockRegex.ReplaceAllString(text, "")

	// Remove any remaining triple backticks blocks that might contain JSON
	codeBlockRegex := regexp.MustCompile("(?s)```[^`]*```")
	text = codeBlockRegex.ReplaceAllString(text, "")

	// Remove inline JSON code segments (single backticks containing JSON-like content)
	inlineJSONRegex := regexp.MustCompile("`[^`]*\\{[^}]*\\}[^`]*`")
	text = inlineJSONRegex.ReplaceAllString(text, "")

	return text
}
