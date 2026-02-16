package repo

import (
	"chatsockets/internal/domain"
	"context"
)

type LinkRepository interface {
	LinkExists(ctx context.Context, linkDomain *domain.UserChatLinkDomain) (bool, error)
	GetLink(ctx context.Context, linkDomain *domain.UserChatLinkDomain) error
	CreateLink(ctx context.Context, linkDomain *domain.UserChatLinkDomain) error
	UpdateLink(ctx context.Context, linkDomain *domain.UserChatLinkDomain) error
}
