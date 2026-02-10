package dto

import (
	"chatsockets/internal/domain"
	"time"
)

type GetChatResponse struct {
	ChatID    int                  `json:"chat_id"`
	Title     string               `json:"chat_title"`
	CreatedAt time.Time            `json:"created_at"`
	Messages  []GetMessageResponse `json:"messages"`
}

func ToGetChatResponse(messages []GetMessageResponse, data *domain.ChatDomain) *GetChatResponse {
	return &GetChatResponse{
		ChatID:    data.ID,
		Title:     data.Title,
		CreatedAt: data.CreatedAt,
		Messages:  messages,
	}
}
