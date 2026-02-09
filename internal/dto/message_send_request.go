package dto

type MessageSendRequest struct {
	Text     string `json:"text"`
	UserName string `json:"username"`
	ChatID   string `json:"chatid"`
}
