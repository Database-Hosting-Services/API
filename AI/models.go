package ai

var SENDER_TYPE_AI string = "ai"
var SENDER_TYPE_USER string = "user"


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