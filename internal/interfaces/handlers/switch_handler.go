package httpHandlers

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/middleware"
	"chatsockets/internal/services"
	websockets "chatsockets/internal/ws"
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type MessageHandler struct {
	*ErrorHandler
	upgrader       websocket.Upgrader
	messageService *services.MessageService
	messageHub     *websockets.MessageHub
}

func NewMessageHandler(upgrader websocket.Upgrader, mService *services.MessageService, appLogger *zap.Logger, mHub *websockets.MessageHub) *MessageHandler {
	return &MessageHandler{
		ErrorHandler:   NewErrorHandler(appLogger),
		upgrader:       upgrader,
		messageService: mService,
		messageHub:     mHub,
	}
}

// Обслуживает подключение по http и смену протокола на websocket
func (mh *MessageHandler) ConnectToChatWS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Меняем протокол
	ws, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		mh.handleDomainError(w, fmt.Errorf("не удалось подключиться к чату"))
		mh.apiLogger.Warn("не удалось сменить протокол на websocket", zap.Error(err))
		return
	}

	// Получаем chat id из url
	chatID, err := mh.parseID(r)
	if err != nil {
		mh.handleDomainError(w, fmt.Errorf("не удалось получить chat id из url"))
		mh.apiLogger.Warn("не удалось получить chat id из url", zap.Error(err))
		return
	}

	// Получаем user id из jwt токена
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		mh.handleDomainError(w, domain.ErrBadJWT)
		mh.apiLogger.Error("не удалось получить id пользователя из токена", zap.Error(domain.ErrBadJWT))
		return
	}

	wsCtx, cancel := context.WithCancel(context.Background())

	//Создаем сущность клиента
	client := websockets.NewClient(ws, mh.apiLogger, userID, chatID, wsCtx, cancel)

	//Регестрируем новое подключение
	mh.messageHub.Register <- client

	go client.WritePump()
	go client.ReadPump(mh.messageService, mh.messageHub)
}
