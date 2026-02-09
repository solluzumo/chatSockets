package httpHandlers

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/dto"
	"chatsockets/internal/services"
	"net/http"

	"go.uber.org/zap"
)

type ChatHandler struct {
	*ErrorHandler
	chatService *services.ChatService
	apiLogger   *zap.Logger
}

func NewChatAPIHTTP(mService *services.ChatService, appLogger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		ErrorHandler: NewErrorHandler(appLogger),
		chatService:  mService,
		apiLogger:    appLogger.Named("chat_api_http"),
	}
}

// Создать чат
func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateChatRequest
	if err := ch.decodeJSON(r, &req); err != nil {
		ch.respondError(w, "некорректный JSON", http.StatusBadRequest, err)
		return
	}

	// валидация title
	if err := req.Validate(); err != nil {
		ch.respondError(w, "ошибка валидации title", http.StatusBadRequest, err)
		return
	}

	chat := &domain.ChatDomain{Title: req.Title}
	result, err := ch.chatService.CreateChat(r.Context(), chat)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	ch.respondJSON(w, http.StatusCreated, &dto.CreateChatResponse{
		ID:        result.ID,
		Title:     result.Title,
		CreatedAt: result.CreatedAt,
	})
}

// Получить чат и limit сообщений
func (ch *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {

	id, err := ch.parseID(r)
	if err != nil {
		ch.respondError(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	limit := ch.parseLimit(r)

	result, err := ch.chatService.GetChatById(r.Context(), id, limit)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	messagesResponse := make([]dto.GetMessageResponse, len(result.Messages))

	for i, msg := range result.Messages {
		messagesResponse[i] = dto.GetMessageResponse{
			ID:        msg.ID,
			Text:      msg.Text,
			UserName:  msg.UserName,
			CreatedAt: msg.CreatedAt,
		}
	}

	ch.respondJSON(w, http.StatusOK, &dto.GetChatResonse{
		ChatID:    result.ID,
		Title:     result.Title,
		CreatedAt: result.CreatedAt,
		Messages:  messagesResponse,
	})
}

// Удаление чата
func (ch *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	id, err := ch.parseID(r)
	if err != nil {
		ch.respondError(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if err := ch.chatService.DeleteChatByID(r.Context(), id); err != nil {
		ch.handleDomainError(w, err)
		return
	}

	ch.respondJSON(w, http.StatusOK, &dto.DeleteChatResponse{
		Content:    "чат успешно удалён",
		StatusCode: http.StatusOK, // Обычно No Content (204) не возвращает тело, но оставил как у вас
	})
}

// Отправка сообщения
func (ch *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	id, err := ch.parseID(r)
	if err != nil {
		ch.respondError(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	var req dto.CreateMessageRequest
	if err := ch.decodeJSON(r, &req); err != nil {
		ch.respondError(w, "некорректный JSON", http.StatusBadRequest, err)
		return
	}

	// валидация text
	if err := req.Validate(); err != nil {
		ch.respondError(w, "ошибка валидации text", http.StatusBadRequest, err)
		return
	}

	msgDomain := &domain.MessageDomain{
		ChatID: id,
		Text:   req.Text,
	}

	result, err := ch.chatService.SendMessage(r.Context(), msgDomain)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	ch.respondJSON(w, http.StatusOK, result)
}
