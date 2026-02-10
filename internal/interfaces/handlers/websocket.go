package httpHandlers

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/dto"
	"chatsockets/internal/middleware"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	*ErrorHandler
	conn       *websocket.Conn
	send       chan any
	msgChannel chan domain.MessageTask
	apiLogger  *zap.Logger
}

func NewClient(conn *websocket.Conn, mChan chan domain.MessageTask, apiLogger *zap.Logger) *Client {
	return &Client{
		ErrorHandler: NewErrorHandler(apiLogger),
		conn:         conn,
		send:         make(chan any, 5),
		msgChannel:   mChan,
		apiLogger:    apiLogger,
	}
}

func (c *Client) readPump(conn *websocket.Conn, r *http.Request) {
	defer conn.Close()

	//ставим дедлайн для проверки подключения
	conn.SetReadDeadline(time.Now().Add(pongWaitS))

	//обновляем счётчик после получения понга
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWaitS))
		return nil
	})

	for {
		ctx := r.Context()
		// Читаем сообщение из сокета
		_, message, err := conn.ReadMessage()
		if err != nil {
			c.apiLogger.Error("ошибка чтения из сокета: ", zap.Error(err))
			c.handleDomainErrorWS(c.send, err)
			break
		}

		var messageSendRequest dto.CreateMessageRequest
		if err := json.Unmarshal(message, &messageSendRequest); err != nil {
			c.apiLogger.Error("ошибка анмаршелинга: ", zap.Error(err))
			c.handleDomainErrorWS(c.send, err)
			continue
		}

		messageDomain, err := messageSendRequest.ToDomain()
		if err != nil {
			c.apiLogger.Error("плохой запрос: ", zap.Error(err))
			c.handleDomainErrorWS(c.send, err)
			continue
		}

		//Достаём userID из контекста
		userID, ok := middleware.UserIDFromContext(ctx)
		if !ok {
			c.apiLogger.Error("ошибка получения id пользователя из контекста")
			c.handleDomainErrorWS(c.send, domain.ErrBadJWT)
			continue
		}

		messageDomain.UserID = userID

		messageTask := domain.MessageTask{
			Ctx:  ctx,
			Data: messageDomain,
		}

		select {
		case c.msgChannel <- messageTask:
		case <-ctx.Done():
			return
		}
	}

}

func (c *Client) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingWaitS)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}

			if err := conn.WriteJSON(msg); err != nil {
				return
			}

		case <-ticker.C:
			if err := conn.WriteControl(
				websocket.PingMessage,
				[]byte{},
				time.Now().Add(writeWait),
			); err != nil {
				return
			}
		}
	}
}
