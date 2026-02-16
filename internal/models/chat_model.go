package models

import (
	"chatsockets/internal/domain"
	"time"
)

type Chat struct {
	ID         int               `gorm:"primaryKey;autoIncrement"`
	Title      string            `gorm:"size:200;not null"`
	CreatedAt  time.Time         `gorm:"not null"`
	ChatStatus domain.ChatStatus `gorm:"not null"`
}

func (c *Chat) ToDomain(message []*domain.MessageDomain) *domain.ChatDomain {
	return &domain.ChatDomain{
		ID:         c.ID,
		Title:      c.Title,
		CreatedAt:  c.CreatedAt,
		Messages:   message,
		ChatStatus: c.ChatStatus,
	}
}

func ToModelChat(data *domain.ChatDomain) *Chat {
	return &Chat{
		ID:         data.ID,
		Title:      data.Title,
		CreatedAt:  data.CreatedAt,
		ChatStatus: data.ChatStatus,
	}
}
