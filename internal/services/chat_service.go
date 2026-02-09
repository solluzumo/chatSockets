package services

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/repo"
	"context"

	"go.uber.org/zap"
)

type ChatService struct {
	chatRepo      repo.ChatRepostiory
	messageRepo   repo.MessageRepostiory
	serviceLogger *zap.Logger
}

func NewChatService(mRepo repo.MessageRepostiory, cRepo repo.ChatRepostiory, appLogger *zap.Logger) *ChatService {
	serviceLogger := appLogger.Named("chat_service")
	return &ChatService{
		chatRepo:      cRepo,
		messageRepo:   mRepo,
		serviceLogger: serviceLogger,
	}
}

func (cs *ChatService) ChatExists(ctx context.Context, param repo.FilterParam) (bool, error) {
	return cs.chatRepo.ChatExists(ctx, param)
}

func (cs *ChatService) CreateChat(ctx context.Context, chat *domain.ChatDomain) (*domain.ChatDomain, error) {
	chatExists, err := cs.ChatExists(ctx, repo.FilterParam{Field: "title", Value: chat.Title})

	if chatExists {
		return nil, domain.ErrChatAlreadyExists
	}

	if err != nil {
		return nil, domain.ErrFieldIsNotAllowed
	}

	chat, err = cs.chatRepo.CreateChat(ctx, chat)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (cs *ChatService) GetChatById(ctx context.Context, chatID int, limit int) (*domain.ChatDomain, error) {
	chatDomain := &domain.ChatDomain{
		ID: chatID,
	}

	err := cs.chatRepo.FindChatById(ctx, chatDomain)
	if err != nil {
		return nil, domain.ErrChatNotFound
	}

	messages := cs.messageRepo.GetMessagesByChatWithLimit(ctx, chatDomain.ID, limit)
	chatDomain.Messages = messages

	return chatDomain, nil
}

func (cs *ChatService) DeleteChatByID(ctx context.Context, chatID int) error {

	chatExists, err := cs.ChatExists(ctx, repo.FilterParam{Field: "id", Value: chatID})

	if err != nil {
		return domain.ErrFieldIsNotAllowed
	}

	if !chatExists {
		return domain.ErrChatNotFound
	}

	if err := cs.chatRepo.DeleteChat(ctx, chatID); err != nil {
		cs.serviceLogger.Error("не удалось удалить чат", zap.Error(err))
		return err
	}

	//УДАЛЯЕМ ВСЕ СООБЩЕНИЯ, не обязательно, так как в бд уже есть ON_DELETE
	// if err := cs.messageRepo.DeleteMessages(ctx, chatID); err != nil {
	// 	cs.serviceLogger.Error("не удалось удалить чат", zap.Error(err))
	// 	return err
	// }

	return nil
}

func (cs *ChatService) SendMessage(ctx context.Context, message *domain.MessageDomain) (*domain.MessageDomain, error) {
	chatExists, err := cs.ChatExists(ctx, repo.FilterParam{Field: "id", Value: message.ChatID})

	if err != nil {
		return nil, domain.ErrFieldIsNotAllowed
	}

	if !chatExists {
		return nil, domain.ErrChatNotFound
	}

	//создаём сообщение
	message, err = cs.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		cs.serviceLogger.Error("не удалось создать сообщение", zap.Error(err))
		return nil, err
	}

	return message, nil
}
