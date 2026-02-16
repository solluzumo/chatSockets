package httpHandlers

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/dto"
	"chatsockets/internal/middleware"
	"chatsockets/internal/services"
	"net/http"

	"go.uber.org/zap"
)

type LinkHandler struct {
	*ErrorHandler
	linkService *services.LinkService
}

func NewLinkAPIHTTP(lService *services.LinkService, appLogger *zap.Logger) *LinkHandler {
	return &LinkHandler{
		ErrorHandler: NewErrorHandler(appLogger),
		linkService:  lService,
	}
}

func (lh *LinkHandler) SubscribeToChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chatID, err := lh.parseID(r)
	if err != nil {
		lh.handleDomainError(w, err)
		return
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		lh.apiLogger.Warn("ошибка получения id пользователя из контекста")
		lh.handleDomainError(w, domain.ErrBadJWT)
		return
	}

	if err := lh.linkService.SubscribeToChat(ctx, chatID, userID); err != nil {
		lh.handleDomainError(w, err)
		return
	}

	lh.respondJSON(w, http.StatusOK, nil)

}

func (lh *LinkHandler) AddUserToChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chatID, err := lh.parseID(r)
	if err != nil {
		lh.handleDomainError(w, err)
		return
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		lh.apiLogger.Warn("ошибка получения id пользователя из контекста")
		lh.handleDomainError(w, domain.ErrBadJWT)
		return
	}

	var req dto.AddUserRequest
	if err := lh.decodeJSON(r, &req); err != nil {
		lh.handleDomainError(w, err)
		return
	}

	linkDomain, err := req.ToDomain()
	if err != nil {
		lh.apiLogger.Warn("ошибка валидации")
		lh.handleDomainError(w, domain.ErrBadJWT)
		return
	}

	linkDomain.ChatID = chatID

	if err := lh.linkService.AddUserToChat(ctx, linkDomain, userID); err != nil {
		lh.handleDomainError(w, err)
		return
	}

	lh.respondJSON(w, http.StatusOK, nil)
}

func (lh *LinkHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.UpdateUserRoleRequest
	if err := lh.decodeJSON(r, &req); err != nil {
		lh.handleDomainError(w, err)
		return
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		lh.apiLogger.Warn("ошибка получения id пользователя из контекста")
		lh.handleDomainError(w, domain.ErrBadJWT)
		return
	}

	chatID, err := lh.parseID(r)
	if err != nil {
		lh.handleDomainError(w, err)
		return
	}

	linkDomain, err := req.ToDomain(chatID)
	if err != nil {
		lh.apiLogger.Warn("ошибка валидации")
		lh.handleDomainError(w, domain.ErrBadRequest)
		return
	}

	if err := lh.linkService.UpdateUserRole(ctx, linkDomain, userID); err != nil {
		lh.handleDomainError(w, err)
		return
	}

	lh.respondJSON(w, http.StatusOK, nil)
}
