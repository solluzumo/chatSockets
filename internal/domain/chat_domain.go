package domain

import (
	"time"
)

type ChatDomain struct {
	ID        int
	Title     string
	CreatedAt time.Time
	Messages  []*MessageDomain
}
