package models

import (
	"chatsockets/internal/domain"
	"time"
)

type ChatWithUsers struct {
	ID         int
	Title      string
	CreatedAt  time.Time
	ChatStatus domain.ChatStatus
	Users      []*UserChatLink `gorm:"foreignKey:ChatID"`
}

func (c *ChatWithUsers) ToDomain() *domain.ChatDomain {
	users := make([]*domain.UserChatLinkDomain, len(c.Users))

	for idx, u := range c.Users {
		users[idx] = u.ToDomain()
	}

	return &domain.ChatDomain{
		ID:         c.ID,
		Title:      c.Title,
		CreatedAt:  c.CreatedAt,
		ChatStatus: c.ChatStatus,
		Users:      users,
	}
}
