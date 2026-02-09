package httpHandlers

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/services"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	pongWaitS = 10 * time.Second
	pingWaitS = (pongWaitS * 9) / 10
	writeWait = 10 * time.Second
)

type MessageHandler struct {
	*ErrorHandler
	upgrader   websocket.Upgrader
	msgChannel chan domain.MessageTask
	messageHub services.MessageHub
	apiLogger  *zap.Logger
}

func NewMessageHandler(upgrader websocket.Upgrader, msgChan chan domain.MessageTask, mhub *services.MessageHub, appLogger *zap.Logger) *MessageHandler {
	return &MessageHandler{
		ErrorHandler: NewErrorHandler(appLogger),
		upgrader:     upgrader,
		msgChannel:   msgChan,
		messageHub:   *mhub,
		apiLogger:    appLogger.Named("message_websocket_api_http"),
	}
}

func (mh *MessageHandler) ConnectToChatWS(w http.ResponseWriter, r *http.Request) {
	ws, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		mh.ErrorHandler.handleDomainError(w, err)
		return
	}

	//Регестрируем новое подключение
	mh.messageHub.Register(ws)

	go mh.writePump(ws)
	go mh.readPump(ws, r)
}

func (mh *MessageHandler) readPump(conn *websocket.Conn, r *http.Request) {
	defer conn.Close()

	//ставим дедлайн для проверки подключения
	conn.SetReadDeadline(time.Now().Add(pongWaitS))

	//обновляем счётчик после получения понга
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWaitS))
		return nil
	})

	for {
		// Читаем сообщение из сокета
		_, message, err := conn.ReadMessage()
		if err != nil {
			mh.apiLogger.Error("ошибка чтения из сокета: ", zap.Error(err))
			break
		}

		messageTask := domain.MessageTask{
			Ctx:  r.Context(),
			Data: message,
		}

		mh.msgChannel <- messageTask
	}

}

func (mh *MessageHandler) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingWaitS)

	defer func() {
		ticker.Stop()
		conn.Close()
		mh.messageHub.Unregister(conn)
	}()

	for range ticker.C {
		if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
			mh.apiLogger.Error("ошибка отправки пинга: ", zap.Error(err))

			return
		}
	}
}
