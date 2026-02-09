package services

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/dto"
	"chatsockets/internal/middleware"
	"chatsockets/internal/repo"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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
	mh.RWmutex.Lock()
	delete(mh.Clients, oldClient)
	mh.RWmutex.Unlock()
}

func (mh *MessageHub) Register(newClient *websocket.Conn) {
	mh.RWmutex.Lock()
	mh.Clients[newClient] = true
	mh.RWmutex.Unlock()
}

func (mh *MessageHub) MessageWorker(id int, ctx context.Context, msgChannel chan domain.MessageTask) {
	defer mh.WG.Done()
	log.Printf("Воркер %d запущен\n", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Воркер %d завершается по контексту", id)
			return

		case msg, ok := <-msgChannel:
			if !ok {
				log.Printf("Воркер %d: канал закрыт, вTыходим", id)
				return
			}

			//Достаём userID из контекста
			userID, ok := middleware.UserIDFromContext(msg.Ctx)
			if !ok {
				mh.sLogger.Error("ошибка получения id пользователя из контекста")
			}

			if err := mh.Broadcast(ctx, msg.Data, userID); err != nil {
				mh.sLogger.Error("не удалось забродкастить сообщение", zap.Error(err))
			}

			log.Printf("Воркер %d обработал сообщение", id)

		}
	}
}

func (mh *MessageHub) LoadAllMessages(ctx context.Context, client *websocket.Conn, chatID int, limit int) error {
	chatDomain := &domain.ChatDomain{
		ID: chatID,
	}

	err := mh.cRepo.FindChatById(ctx, chatDomain)
	if err != nil {
		return domain.ErrChatNotFound
	}

	messages := mh.mRepo.GetMessagesByChatWithLimit(ctx, chatID, limit)

	chatDomain.Messages = messages

	//Отправляем сначала название чата
	err = client.WriteMessage(websocket.TextMessage, []byte(chatDomain.Title))
	if err != nil {
		log.Printf("ошибка отправки клиенту: %v", err)
	}

	wsMessages := make([]dto.MessageResponseWebsocket, len(messages))

	for idx := range messages {
		wsMessages[idx].ID = messages[idx].ID
		wsMessages[idx].Text = messages[idx].Text
	}

	getChatResponse := dto.GetChatResonseWebsocket{
		ChatID:    chatDomain.ID,
		Title:     chatDomain.Title,
		CreatedAt: chatDomain.CreatedAt,
		Messages:  wsMessages,
	}

	jsonData, _ := json.Marshal(getChatResponse)

	//Отправляем клиенту сообщение
	mh.RWmutex.RLock()
	err = client.WriteMessage(websocket.TextMessage, jsonData)
	mh.RWmutex.RUnlock()

	if err != nil {
		log.Printf("ошибка отправки клиенту: %v", err)
	}

	return nil
}

func (mh *MessageHub) ChatExists(ctx context.Context, param repo.FilterParam) (bool, error) {
	return mh.cRepo.ChatExists(ctx, param)
}

func (mh *MessageHub) Broadcast(ctx context.Context, msg []byte, userID int) error {
	var messageSendRequest dto.MessageSendRequest
	if err := json.Unmarshal(msg, &messageSendRequest); err != nil {
		return err
	}

	fmt.Println("запрос на отправку сообщения: ", messageSendRequest)
	chatID, err := strconv.Atoi(messageSendRequest.ChatID)
	if err != nil {
		return err
	}

	messageDomain := &domain.MessageDomain{
		Text:   messageSendRequest.Text,
		ChatID: chatID,
		UserID: userID,
	}

	chatExists, err := mh.ChatExists(ctx, repo.FilterParam{Field: "id", Value: messageDomain.ChatID})

	if err != nil {
		return domain.ErrFieldIsNotAllowed
	}

	if !chatExists {
		return domain.ErrChatNotFound
	}

	//Проверяем связан ли пользователь и чат
	if err := mh.cRepo.IsUserConnectedToChat(ctx, messageDomain.ChatID, userID); err != true {
		return domain.ErrUserIsNotConnectedToChat
	}

	// Блокируем чтение мапы, чтобы никто не менял её в этот момент и бродкастим сообщение
	mh.RWmutex.RLock()
	for client := range mh.Clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("ошибка отправки клиенту: %v", err)
		}
	}
	mh.RWmutex.RUnlock()

	//Создаем запись в бд
	mh.mRepo.CreateMessage(ctx, messageDomain)

	return nil
}
