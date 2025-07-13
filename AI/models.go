package ai

import (
	"DBHS/utils"
)

const SENDER_TYPE_AI = "ai"
const SENDER_TYPE_USER = "user"


type Analytics struct {}

type ChatBotRequest struct {
	Question string `json:"question"`
}

type ChatData struct {
	ID        int    `json:"id"`
	Oid       string `json:"oid"`
	OwnerID   int    `json:"owner_id"`
	ProjectID int    `json:"project_id"`
}

type Request struct {
	Prompt string `json:"prompt"`
}

type AgentResponse struct {
	Response      string  `json:"response"`
    SchemaChanges []utils.Table `json:"schema_changes"`
    SchemaDDL     string  `json:"schema_ddl"`
}