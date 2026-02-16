package domain

import (
	"time"
)

type ChatDomain struct {
	ID         int
	Title      string
	CreatedAt  time.Time
	ChatStatus ChatStatus
	Messages   []*MessageDomain
	Users      []*UserChatLinkDomain
}
