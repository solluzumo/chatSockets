package dto

import (
	"chatsockets/internal/domain"
	"time"
)

type CreateChatResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func ToCreateChatResponse(data *domain.ChatDomain) *CreateChatResponse {
	return &CreateChatResponse{
		ID:        data.ID,
		Title:     data.Title,
		CreatedAt: data.CreatedAt,
	}
}
