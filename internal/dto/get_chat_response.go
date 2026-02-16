package dto

import (
	"chatsockets/internal/domain"
	"time"
)

type GetChatResponse struct {
	ChatID     int                          `json:"chat_id"`
	Title      string                       `json:"chat_title"`
	ChatStatus string                       `json:"chat_status"`
	CreatedAt  time.Time                    `json:"created_at"`
	Messages   []GetMessageResponse         `json:"messages"`
	Users      []*domain.UserChatLinkDomain `json:"users"`
}

func ToGetChatResponse(messages []GetMessageResponse, data *domain.ChatDomain) *GetChatResponse {
	return &GetChatResponse{
		ChatID:     data.ID,
		Title:      data.Title,
		ChatStatus: string(data.ChatStatus),
		CreatedAt:  data.CreatedAt,
		Messages:   messages,
		Users:      data.Users,
	}
}
