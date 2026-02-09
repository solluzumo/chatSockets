package dto

import (
	"time"
)

type GetChatResonse struct {
	ChatID    int                  `json:"chat_id"`
	Title     string               `json:"chat_title"`
	CreatedAt time.Time            `json:"created_at"`
	Messages  []GetMessageResponse `json:"messages"`
}
