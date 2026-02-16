package services

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/repo"
	"context"
	"fmt"
)

type ChatService struct {
	permissionService *PermissionService
	cRepo             repo.ChatRepostiory
	mRepo             repo.MessageRepostiory
	lRepo             repo.LinkRepository
}

func NewChatService(mRepo repo.MessageRepostiory, cRepo repo.ChatRepostiory, lRepo repo.LinkRepository, pService *PermissionService) *ChatService {
	return &ChatService{
		cRepo:             cRepo,
		mRepo:             mRepo,
		lRepo:             lRepo,
		permissionService: pService,
	}
}

func (cs *ChatService) CreateChat(ctx context.Context, chat *domain.ChatDomain, userID int) error {

	//Создаём чат
	err := cs.cRepo.CreateChat(ctx, chat)
	if err != nil {
		return fmt.Errorf("не удалось создать чат: %w", err)
	}

	//Создаём связь пользователь-чат
	linkDomain := &domain.UserChatLinkDomain{
		UserID:      userID,
		ChatID:      chat.ID,
		UserBlocked: false,
		UserRole:    domain.AdminRole,
	}

	if err := cs.lRepo.CreateLink(ctx, linkDomain); err != nil {
		return fmt.Errorf("не удалось привязать пользователя к чату: %w", err)
	}
	return nil
}

func (cs *ChatService) GetChatById(ctx context.Context, chatID int, limit int) (*domain.ChatDomain, error) {

	chatDomain, err := cs.cRepo.GetChatWithUsers(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти чат: %w", domain.ErrChatNotFound)
	}

	messages, err := cs.mRepo.GetMessagesByChatWithLimit(ctx, chatDomain.ID, limit)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить сообщения чата: %w", err)
	}

	chatDomain.Messages = messages

	return chatDomain, nil
}

func (cs *ChatService) DeleteChatByID(ctx context.Context, chatID int) error {

	if err := cs.cRepo.DeleteChat(ctx, chatID); err != nil {
		return fmt.Errorf("не удалось удалить чат по id: %w", err)
	}

	return nil
}
