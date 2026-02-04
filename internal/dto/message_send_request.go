package dto

type MessageSendRequest struct {
	Message   string `json:"message"`
	UserName  string `json:"username"`
	TimeStamp string `json:"timestamp"`
}
