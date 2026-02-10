package services

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/repo"
	"context"
	"fmt"
)

type ChatService struct {
	chatRepo    repo.ChatRepostiory
	messageRepo repo.MessageRepostiory
}

func NewChatService(mRepo repo.MessageRepostiory, cRepo repo.ChatRepostiory) *ChatService {
	return &ChatService{
		chatRepo:    cRepo,
		messageRepo: mRepo,
	}
}

func (cs *ChatService) ChatExists(ctx context.Context, param repo.FilterParam) (bool, error) {
	return cs.chatRepo.ChatExists(ctx, param)
}

func (cs *ChatService) CreateChat(ctx context.Context, chat *domain.ChatDomain) (*domain.ChatDomain, error) {

	chat, err := cs.chatRepo.CreateChat(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать чат: %w", err)
	}

	return chat, nil
}

func (cs *ChatService) GetChatById(ctx context.Context, chatID int, limit int) (*domain.ChatDomain, error) {
	chatDomain := &domain.ChatDomain{
		ID: chatID,
	}
	if err := cs.chatRepo.FindChatById(ctx, chatDomain); err != nil {
		return nil, fmt.Errorf("не удалось найти чат: %w", domain.ErrChatNotFound)
	}

	messages, err := cs.messageRepo.GetMessagesByChatWithLimit(ctx, chatDomain.ID, limit)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить сообщения чата: %w", err)
	}

	chatDomain.Messages = messages

	return chatDomain, nil
}

func (cs *ChatService) DeleteChatByID(ctx context.Context, chatID int) error {

	if err := cs.chatRepo.DeleteChat(ctx, chatID); err != nil {
		return fmt.Errorf("не удалось удалить чат по id: %w", err)
	}

	return nil
}

func (cs *ChatService) SendMessage(ctx context.Context, message *domain.MessageDomain) (*domain.MessageDomain, error) {

	//создаём сообщение
	message, err := cs.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("не удалось отправить сообщение: %w", err)
	}

	return message, nil
}
