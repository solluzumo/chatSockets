package services

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/events"
	"chatsockets/internal/repo"
	"context"
	"fmt"
)

type MessageService struct {
	permissionService *PermissionService
	mRepo             repo.MessageRepostiory
	cRepo             repo.ChatRepostiory
	lRepo             repo.LinkRepository
	bus               *events.EventBus
}

func NewMessageService(mRepo repo.MessageRepostiory, cRepo repo.ChatRepostiory, lRepo repo.LinkRepository, bus *events.EventBus, permissionService *PermissionService) *MessageService {
	return &MessageService{
		mRepo:             mRepo,
		cRepo:             cRepo,
		lRepo:             lRepo,
		bus:               bus,
		permissionService: permissionService,
	}
}

func (ms *MessageService) SendMessage(ctx context.Context, message *domain.MessageDomain) error {

	_, err := ms.permissionService.CanUserThis(ctx, message.UserID, message.ChatID, domain.SendMessage)
	if err != nil {
		return fmt.Errorf("пользователь не может отправить сообщение: %w", err)
	}

	//создаём сообщение
	message, err = ms.mRepo.CreateMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("не удалось отправить сообщение: %w", err)
	}

	ms.bus.PublicTask(ctx, message)

	return nil
}
