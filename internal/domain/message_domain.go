package domain

import (
	"time"
)

type MessageDomain struct {
	ID        int
	Text      string
	ChatID    int
	UserID    int
	UserName  string
	CreatedAt time.Time
}
