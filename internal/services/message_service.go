package services

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/dto"
	"chatsockets/internal/repo"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type MessageHub struct {
	RWmutex *sync.RWMutex
	Clients map[*websocket.Conn]bool
	WG      *sync.WaitGroup
	sLogger *zap.Logger
	mRepo   repo.MessageRepostiory
	cRepo   repo.ChatRepostiory
}

func NewMessageHub(rwmut *sync.RWMutex, clients map[*websocket.Conn]bool, WG *sync.WaitGroup, appLogger *zap.Logger, mRepo repo.MessageRepostiory, cRepo repo.ChatRepostiory) *MessageHub {
	return &MessageHub{
		RWmutex: rwmut,
		Clients: clients,
		WG:      WG,
		sLogger: appLogger.Named("message_service"),
		mRepo:   mRepo,
		cRepo:   cRepo,
	}
}

func (mh *MessageHub) Unregister(oldClient *websocket.Conn) {
	mh.sLogger.Info("клиент отключён")
	mh.RWmutex.Lock()
	delete(mh.Clients, oldClient)
	mh.RWmutex.Unlock()
}

func (mh *MessageHub) Register(newClient *websocket.Conn) {
	mh.sLogger.Info("зарегестрирован новый клиент", zap.Int("№", len(mh.Clients)))
	mh.RWmutex.Lock()
	mh.Clients[newClient] = true
	mh.RWmutex.Unlock()
}

func (mh *MessageHub) ChatExists(ctx context.Context, param repo.FilterParam) (bool, error) {
	return mh.cRepo.ChatExists(ctx, param)
}

func (mh *MessageHub) MessageWorker(id int, ctx context.Context, msgChannel chan domain.MessageTask) {
	defer mh.WG.Done()
	log.Printf("Воркер %d запущен\n", id)

	for {
		select {
		case <-ctx.Done():
			mh.sLogger.Info("Завершается по контексту", zap.Int("Воркер", id))
			return

		case msg, ok := <-msgChannel:
			if !ok {
				mh.sLogger.Info("Канал закрыт, выходим", zap.Int("Воркер", id))
				return
			}

			if err := mh.Broadcast(ctx, msg.Data); err != nil {
				mh.sLogger.Error("не удалось забродкастить сообщение", zap.Error(err))
			}

			mh.sLogger.Info("Обработал сообщение", zap.Int("Воркер", id))

		}
	}
}

func (mh *MessageHub) Broadcast(ctx context.Context, messageDomain *domain.MessageDomain) error {

	//Проверяем связан ли пользователь и чат
	if err := mh.cRepo.IsUserConnectedToChat(ctx, messageDomain.ChatID, messageDomain.UserID); err != true {
		return fmt.Errorf("пользователь не является участником чата")
	}

	msgResponse, err := dto.ToGetMessageResponse(messageDomain).ToBytes()
	if err != nil {
		return fmt.Errorf("не удалось преобразовать дто в слайс байтов: %w", err)
	}

	// Блокируем чтение мапы, чтобы никто не менял её в этот момент и бродкастим сообщение
	mh.RWmutex.RLock()
	for client := range mh.Clients {
		err := client.WriteMessage(websocket.TextMessage, msgResponse)
		if err != nil {
			mh.sLogger.Warn("ошибка отправки клиенту", zap.Error(err))
		}
	}
	mh.RWmutex.RUnlock()

	//Создаем запись в бд
	_, err = mh.mRepo.CreateMessage(ctx, messageDomain)
	if err != nil {
		return fmt.Errorf("не удалось создать сообщение: %w", err)
	}

	return nil
}

// func (mh *MessageHub) LoadAllMessages(ctx context.Context, client *websocket.Conn, chatID int, limit int) error {
// 	chatDomain := &domain.ChatDomain{
// 		ID: chatID,
// 	}

// 	err := mh.cRepo.FindChatById(ctx, chatDomain)
// 	if err != nil {
// 		return domain.ErrChatNotFound
// 	}

// 	messages := mh.mRepo.GetMessagesByChatWithLimit(ctx, chatID, limit)

// 	chatDomain.Messages = messages

// 	//Отправляем сначала название чата
// 	err = client.WriteMessage(websocket.TextMessage, []byte(chatDomain.Title))
// 	if err != nil {
// 		log.Printf("ошибка отправки клиенту: %w", err)
// 	}

// 	wsMessages := make([]dto.GetMessageResponseWS, len(messages))

// 	for idx := range messages {
// 		wsMessages[idx].ID = messages[idx].ID
// 		wsMessages[idx].Text = messages[idx].Text
// 	}

// 	getChatResponse := dto.GetChatResonseWebsocket{
// 		ChatID:    chatDomain.ID,
// 		Title:     chatDomain.Title,
// 		CreatedAt: chatDomain.CreatedAt,
// 		Messages:  wsMessages,
// 	}

// 	jsonData, _ := json.Marshal(getChatResponse)

// 	//Отправляем клиенту сообщение
// 	mh.RWmutex.RLock()
// 	err = client.WriteMessage(websocket.TextMessage, jsonData)
// 	mh.RWmutex.RUnlock()

// 	if err != nil {
// 		log.Printf("ошибка отправки клиенту: %w", err)
// 	}

// 	return nil
// }
