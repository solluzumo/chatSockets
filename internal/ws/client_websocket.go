package websockets

import (
	"chatsockets/internal/dto"
	"chatsockets/internal/services"
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	pongWaitS = 10 * time.Second
	pingWaitS = (pongWaitS * 9) / 10
	writeWait = 10 * time.Second
)

type Client struct {
	*ErrorHandlerWS
	UserID    int
	ChatID    int
	conn      *websocket.Conn
	send      chan any //в этот канал пишем всё что отправляем клиенту
	apiLogger *zap.Logger
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewClient(conn *websocket.Conn, apiLogger *zap.Logger, userID, chatID int, ctx context.Context, cancel context.CancelFunc) *Client {

	return &Client{
		ErrorHandlerWS: NewErrorHandlerWS(apiLogger.Named("client_websocket")),
		conn:           conn,
		send:           make(chan any, 5),
		apiLogger:      apiLogger.Named("websocket_client"),
		UserID:         userID,
		ChatID:         chatID,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (c *Client) ReadPump(mService *services.MessageService, mh *MessageHub) {
	defer func() {
		c.apiLogger.Info("Закрываем соединение")
		c.cancel()
		mh.Unregister <- c
		c.conn.Close()
	}()

	//ставим дедлайн для проверки подключения
	c.conn.SetReadDeadline(time.Now().Add(pongWaitS))

	//обновляем счётчик после получения понга
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWaitS))
		return nil
	})

	for {
		// Читаем сообщение из сокета
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.apiLogger.Error("ошибка чтения из сокета: ", zap.Error(err))
			c.handleDomainErrorWS(c.send, err)
			break
		}
		//Анмаршелим запрос
		var messageSendRequest dto.CreateMessageRequest
		if err := json.Unmarshal(message, &messageSendRequest); err != nil {
			c.apiLogger.Error("ошибка анмаршелинга: ", zap.Error(err))
			c.handleDomainErrorWS(c.send, err)
			continue
		}

		//Преобразовывам в домен
		messageDomain, err := messageSendRequest.ToDomain(c.UserID)
		if err != nil {
			c.apiLogger.Error("плохой запрос: ", zap.Error(err))
			c.handleDomainErrorWS(c.send, err)
			continue
		}

		if err := mService.SendMessage(c.ctx, messageDomain); err != nil {
			c.apiLogger.Error("не удалось создать сообщение", zap.Error(err))
			c.handleDomainErrorWS(c.send, err)
			continue
		}
	}

}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingWaitS)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}

			if err := c.conn.WriteJSON(msg); err != nil {
				return
			}

		case <-ticker.C:
			c.apiLogger.Info("Отправлен пинг")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}

	}
}
