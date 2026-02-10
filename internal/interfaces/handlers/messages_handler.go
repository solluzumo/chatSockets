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
}

func NewMessageHandler(upgrader websocket.Upgrader, msgChan chan domain.MessageTask, mhub *services.MessageHub, appLogger *zap.Logger) *MessageHandler {
	return &MessageHandler{
		ErrorHandler: NewErrorHandler(appLogger),
		upgrader:     upgrader,
		msgChannel:   msgChan,
		messageHub:   *mhub,
	}
}

func (mh *MessageHandler) ConnectToChatWS(w http.ResponseWriter, r *http.Request) {
	ws, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		mh.apiLogger.Error("не удалось сменить протокол на websocket", zap.Error(err))
		return
	}
	client := NewClient(ws, mh.msgChannel, mh.apiLogger)
	//Регестрируем новое подключение
	mh.messageHub.Register(ws)

	go client.writePump(ws)
	go client.readPump(ws, r)
}
