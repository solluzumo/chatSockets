package services

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/repo"
	"context"
	"fmt"
)

type LinkService struct {
	lRepo             repo.LinkRepository
	permissionService *PermissionService
}

func NewLinkService(lRepo repo.LinkRepository, permissionService *PermissionService) *LinkService {
	return &LinkService{
		lRepo:             lRepo,
		permissionService: permissionService,
	}
}

// func (cs *ChatService) LinkExists(ctx context.Context, userID int, chatID int) (bool, error) {
// 	//Проверяем связан ли пользователь и чат
// 	exists, err := cs.cRepo.LinkExists(ctx, &domain.UserChatLinkDomain{UserID: userID, ChatID: chatID})

// 	if err != nil {
// 		return false, fmt.Errorf("не удалось проверить связь пользователь-чат: %w", err)
// 	}

// 	if exists != true {
// 		return false, fmt.Errorf("пользователь не является участником чата")
// 	}

// 	return true, nil
// }

func (ls *LinkService) AddUserToChat(ctx context.Context, linkDomain *domain.UserChatLinkDomain, reqAuthorID int) error {
	_, err := ls.permissionService.CanUserThis(ctx, reqAuthorID, linkDomain.ChatID, domain.AddUser)
	if err != nil {
		return fmt.Errorf("не удалось добавить пользователя: %w", err)
	}

	linkDomain.UserRole = domain.MemberRole
	linkDomain.UserBlocked = false

	if err := ls.lRepo.CreateLink(ctx, linkDomain); err != nil {
		return fmt.Errorf("не удалось добавить пользователя: %w", err)
	}

	return nil
}

func (ls *LinkService) UpdateUserRole(ctx context.Context, oldLinkDomain *domain.UserChatLinkDomain, reqAuthorID int) error {
	_, err := ls.permissionService.CanUserThis(ctx, reqAuthorID, oldLinkDomain.ChatID, domain.UpdateRole)
	if err != nil {
		return fmt.Errorf("пользователь не может изменить роль другого пользователя: %w", err)
	}

	if err := ls.lRepo.GetLink(ctx, oldLinkDomain); err != nil {
		return fmt.Errorf("пользователь для обновления не привязан к этому чату: %w", err)
	}

	if err := ls.lRepo.UpdateLink(ctx, oldLinkDomain); err != nil {
		return fmt.Errorf("не удалось обновить роль пользователя: %w", err)
	}
	return nil
}

func (ls *LinkService) SubscribeToChat(ctx context.Context, chatID int, userID int) error {
	linkDomain := &domain.UserChatLinkDomain{
		UserID:      userID,
		ChatID:      chatID,
		UserBlocked: false,
		UserRole:    domain.GuestRole,
	}

	if err := ls.lRepo.CreateLink(ctx, linkDomain); err != nil {
		return fmt.Errorf("не удалось привязать пользователя к чату: %w", err)
	}

	return nil
}
