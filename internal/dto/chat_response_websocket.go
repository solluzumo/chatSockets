package dto

import (
	"time"
)

type GetChatResonseWebsocket struct {
	ChatID    int                        `json:"chat_id"`
	Title     string                     `json:"chat_title"`
	CreatedAt time.Time                  `json:"created_at"`
	Messages  []MessageResponseWebsocket `json:"messages"`
}
