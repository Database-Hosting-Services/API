package ai

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