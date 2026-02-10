package models

import (
	"chatsockets/internal/domain"
	"time"
)

type Message struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	ChatID    int       `gorm:"not null"`
	Text      string    `gorm:"size:5000;not null"`
	UserID    int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (m *Message) ModelToDomain() *domain.MessageDomain {
	return &domain.MessageDomain{
		ID:        m.ID,
		Text:      m.Text,
		ChatID:    m.ChatID,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
	}
}

func DomainToModelMessage(data *domain.MessageDomain) *Message {
	return &Message{
		ID:        data.ID,
		Text:      data.Text,
		ChatID:    data.ChatID,
		UserID:    data.UserID,
		CreatedAt: data.CreatedAt,
	}
}

func ModelSliceToDomainSlice(messageModels []*Message) []*domain.MessageDomain {
	var messageDomains []*domain.MessageDomain

	for _, el := range messageModels {
		messageDomains = append(messageDomains, &domain.MessageDomain{
			ID:        el.ID,
			ChatID:    el.ChatID,
			Text:      el.Text,
			CreatedAt: el.CreatedAt,
		})
	}

	return messageDomains
}
