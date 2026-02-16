package models

import (
	"chatsockets/internal/domain"
	"time"
)

type UserChatLink struct {
	UserID      int             `gorm:"primaryKey"`
	ChatID      int             `gorm:"primaryKey"`
	UserBlocked bool            `gorm:"not null"`
	UserRole    domain.UserRole `gorm:"not null"`
	CreatedAt   time.Time       `gorm:"not null"`
}

func (u *UserChatLink) ToDomain() *domain.UserChatLinkDomain {
	return &domain.UserChatLinkDomain{
		UserID:      u.UserID,
		ChatID:      u.ChatID,
		UserRole:    u.UserRole,
		UserBlocked: u.UserBlocked,
		CreatedAt:   u.CreatedAt,
	}
}

func ToUserChatLink(link *domain.UserChatLinkDomain) *UserChatLink {
	return &UserChatLink{
		UserID:      link.UserID,
		ChatID:      link.ChatID,
		UserBlocked: link.UserBlocked,
		UserRole:    link.UserRole,
	}
}
func (m *UserChatLink) UpdateUserChatLinkDomain(d *domain.UserChatLinkDomain) {
	d.UserID = m.UserID
	d.ChatID = m.ChatID
	d.UserRole = m.UserRole
	d.UserBlocked = m.UserBlocked
	d.CreatedAt = m.CreatedAt
}
