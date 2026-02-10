package repo

import (
	"chatsockets/internal/domain"
	"context"
)

type MessageRepostiory interface {
	CreateMessage(ctx context.Context, data *domain.MessageDomain) (*domain.MessageDomain, error)
	GetMessagesByChatWithLimit(ctx context.Context, chatID int, limit int) ([]*domain.MessageDomain, error)
	DeleteMessages(ctx context.Context, chatID int) error
	Count(ctx context.Context) int64
}
