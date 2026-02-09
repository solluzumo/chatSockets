package repo

import (
	"chatsockets/internal/domain"
	"context"
)

type ChatRepostiory interface {
	CreateChat(ctx context.Context, data *domain.ChatDomain) (*domain.ChatDomain, error)
	FindChatById(ctx context.Context, data *domain.ChatDomain) error
	ChatExists(ctx context.Context, param FilterParam) (bool, error)
	DeleteChat(ctx context.Context, chatID int) error
	IsUserConnectedToChat(ctx context.Context, chatID int, userID int) bool
	Count(ctx context.Context) int64
}
