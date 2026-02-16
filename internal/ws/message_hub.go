package websockets

import (
	"chatsockets/internal/events"
	"context"
	"log"
	"sync"

	"go.uber.org/zap"
)

type MessageHub struct {
	clients    map[int]map[*Client]bool
	WG         *sync.WaitGroup
	sLogger    *zap.Logger
	Register   chan *Client
	Unregister chan *Client
	bus        *events.EventBus
}

func NewMessageHub(WG *sync.WaitGroup, appLogger *zap.Logger, bus *events.EventBus) *MessageHub {
	return &MessageHub{
		clients:    make(map[int]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		bus:        bus,
		WG:         WG,
		sLogger:    appLogger.Named("message_hub"),
	}
}

func (h *MessageHub) addClient(client *Client) {
	h.sLogger.Info("добавляем клиента в пул")
	if h.clients[client.ChatID] == nil {
		h.clients[client.ChatID] = make(map[*Client]bool)
	}

	h.clients[client.ChatID][client] = true
}

func (h *MessageHub) removeClient(client *Client) {
	h.sLogger.Info("удаляем клиента из пула")
	delete(h.clients[client.ChatID], client)
	if len(h.clients[client.ChatID]) == 0 {
		delete(h.clients, client.ChatID)
	}
	close(client.send)
}

func (mh *MessageHub) MessageWorker(id int, ctx context.Context) {
	defer mh.WG.Done()
	log.Printf("Воркер %d запущен\n", id)

	for {
		select {
		case <-ctx.Done():
			mh.sLogger.Info("Завершается по контексту", zap.Int("Воркер", id))
			return

		case client := <-mh.Register:
			mh.addClient(client)

		case client := <-mh.Unregister:
			mh.removeClient(client)

		//Поймали событие создания сообщения
		case event, ok := <-mh.bus.MessageCreated:
			if !ok {
				mh.sLogger.Info("Канал закрыт, выходим", zap.Int("Воркер", id))
				return
			}

			mh.broadcast(event)

			mh.sLogger.Info("Обработал сообщение", zap.Int("Воркер", id))
		}
	}
}

func (mh *MessageHub) broadcast(event events.MessageCreatedEvent) {

	clients := mh.clients[event.Data.ChatID]

	for client := range clients {

		select {
		case client.send <- event:
		default:
			mh.removeClient(client)
		}
	}
}

// func (mh *MessageHub) Broadcast(ctx context.Context, messageDomain *domain.MessageDomain) error {

// 	msgResponse, err := dto.ToGetMessageResponse(messageDomain).ToBytes()
// 	if err != nil {
// 		return fmt.Errorf("не удалось преобразовать дто в слайс байтов: %w", err)
// 	}

// 	// Блокируем чтение мапы, чтобы никто не менял её в этот момент и бродкастим сообщение
// 	mh.RWmutex.RLock()
// 	for client := range mh.Clients {
// 		err := client.WriteMessage(websocket.TextMessage, msgResponse)
// 		if err != nil {
// 			mh.sLogger.Warn("ошибка отправки клиенту", zap.Error(err))
// 		}
// 	}
// 	mh.RWmutex.RUnlock()

// 	return nil
// }

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
