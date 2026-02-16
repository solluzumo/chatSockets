package services

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/repo"
	"context"
	"fmt"
)

type PermissionService struct {
	lRepo repo.LinkRepository
}

func NewPermissionService(lRepo repo.LinkRepository) *PermissionService {
	return &PermissionService{
		lRepo: lRepo,
	}
}

func (ps *PermissionService) CheckPermissions(userRole domain.UserRole, action domain.UserAction) bool {
	return roleMatrix[userRole][action]
}

func (ps *PermissionService) CheckChatPermissions(chatStatus domain.ChatStatus, action domain.UserAction) bool {
	return chatMatrix[chatStatus][action]
}

func (ps *PermissionService) CanUserThis(ctx context.Context, userID int, chatID int, action domain.UserAction) (*domain.UserChatLinkDomain, error) {

	link := &domain.UserChatLinkDomain{
		UserID: userID,
		ChatID: chatID,
	}

	err := ps.lRepo.GetLink(ctx, link)
	if err != nil {
		return nil, fmt.Errorf("пользователь не является участником чата: %w", err)
	}

	switch {
	case link.UserBlocked:
		return nil, fmt.Errorf("пользователь заблокирован: %w", domain.ErrUserBlocked)
	case !ps.CheckPermissions(link.UserRole, action):
		return nil, fmt.Errorf("недостаточно прав для выполнения действия %s: %w", string(action), domain.ErrNotEnoughStatus)
	}
	return link, nil
}
