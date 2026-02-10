package dto

import (
	"chatsockets/internal/domain"
	"time"
)

type CreateMessageResponse struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func ToCreateMessageResponse(data *domain.MessageDomain) *CreateMessageResponse {
	return &CreateMessageResponse{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
	}
}
