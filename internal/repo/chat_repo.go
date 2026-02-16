package repo

import (
	"chatsockets/internal/domain"
	"context"
)

type ChatRepostiory interface {
	CreateChat(ctx context.Context, chatDomain *domain.ChatDomain) error
	ChatExists(ctx context.Context, chatDomain *domain.ChatDomain) (bool, error)
	GetChatWithUsers(ctx context.Context, chatID int) (*domain.ChatDomain, error)
	DeleteChat(ctx context.Context, chatID int) error

	Count(ctx context.Context) int64
}
