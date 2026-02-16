package domain

import (
	"time"
)

type UserChatLinkDomain struct {
	UserID      int
	ChatID      int
	UserRole    UserRole
	UserBlocked bool
	CreatedAt   time.Time
}
