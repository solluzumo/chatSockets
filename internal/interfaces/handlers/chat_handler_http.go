package httpHandlers

import (
	"chatsockets/internal/dto"
	"chatsockets/internal/services"
	"net/http"

	"go.uber.org/zap"
)

type ChatHandler struct {
	*ErrorHandler
	chatService *services.ChatService
}

func NewChatAPIHTTP(mService *services.ChatService, appLogger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		ErrorHandler: NewErrorHandler(appLogger),
		chatService:  mService,
	}
}

// Создать чат
func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateChatRequest
	if err := ch.decodeJSON(r, &req); err != nil {
		ch.handleDomainError(w, err)
		return
	}

	chatDomain, err := req.ToDomain()
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	result, err := ch.chatService.CreateChat(r.Context(), chatDomain)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	response := dto.ToCreateChatResponse(result)

	ch.respondJSON(w, http.StatusCreated, response)
}

// Получить чат и limit сообщений
func (ch *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {

	id, err := ch.parseID(r)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	limit := ch.parseLimit(r)

	result, err := ch.chatService.GetChatById(r.Context(), id, limit)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	messagesResponse := dto.ToGetMessageResponseSlice(result.Messages)
	response := dto.ToGetChatResponse(messagesResponse, result)

	ch.respondJSON(w, http.StatusOK, response)
}

// Удаление чата
func (ch *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	id, err := ch.parseID(r)
	if err != nil {
		ch.handleDomainError(w, err)
		return
	}

	if err := ch.chatService.DeleteChatByID(r.Context(), id); err != nil {
		ch.handleDomainError(w, err)
		return
	}
	response := dto.ToDeleteChatResponse("чат успешно удалён", http.StatusOK)
	ch.respondJSON(w, http.StatusOK, response)
}
