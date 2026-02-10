package models

import (
	"chatsockets/internal/domain"
	"time"
)

type Chat struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"size:200;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (c *Chat) ModelToDomain(message []*domain.MessageDomain) *domain.ChatDomain {

	return &domain.ChatDomain{
		ID:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
		Messages:  message,
	}
}

func DomainToModelChat(data *domain.ChatDomain) *Chat {
	return &Chat{
		ID:        data.ID,
		Title:     data.Title,
		CreatedAt: data.CreatedAt,
	}
}
