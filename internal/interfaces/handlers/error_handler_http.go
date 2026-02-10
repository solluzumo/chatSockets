package httpHandlers

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/dto"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	apiLogger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{
		apiLogger: logger,
	}
}

func (er *ErrorHandler) handleDomainErrorWS(send chan<- any, err error) {
	var errString string
	switch {
	case errors.Is(err, domain.ErrChatNotFound):
		errString = "Чат не найден"
	case errors.Is(err, domain.ErrUserIsNotConnectedToChat):
		errString = "пользователь не участник чата"
	case errors.Is(err, domain.ErrBadJWT):
		errString = "неверный jwt токен"
	default:
		errString = "Внутренняя ошибка сервера"
	}
	er.respondErrorWS(send, errString)
}

func (er *ErrorHandler) respondErrorWS(send chan<- any, err string) {
	errorMsg := dto.WsErrorResponse{
		Type:   "error",
		Error:  "Не удалось подключиться к чату",
		Reason: err,
	}
	send <- errorMsg

}

// respondError отправляет ошибку в формате JSON
func (er *ErrorHandler) respondError(w http.ResponseWriter, message string, code int, err error) {
	if err != nil {
		er.apiLogger.Warn(message, zap.Error(err))
	} else {
		er.apiLogger.Warn(message)
	}
	er.respondJSON(w, code, map[string]string{"error": message})
}

// handleDomainError маппит ошибки домена на HTTP коды
func (er *ErrorHandler) handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrChatNotFound):
		er.respondError(w, "чат не найден", http.StatusNotFound, err)
	case errors.Is(err, domain.ErrChatAlreadyExists):
		er.respondError(w, "чат с таким названием уже существует", http.StatusBadRequest, err)
	case errors.Is(err, domain.ErrUserIsNotConnectedToChat):
		er.respondError(w, "пользователь не участник чата", http.StatusBadRequest, err)
	case errors.Is(err, domain.ErrFieldIsNotAllowed):
		er.respondError(w, "не разрешенное для фильтрации поле", http.StatusBadRequest, err)
	case errors.Is(err, context.Canceled):
		er.apiLogger.Info("запрос отменён клиентом")
	default:
		er.apiLogger.Error("внутренняя ошибка сервера", zap.Error(err))
		er.respondError(w, "internal server error", http.StatusInternalServerError, nil)
	}
}

// parseID извлекает и валидирует ID из URL
func (er *ErrorHandler) parseID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, errors.New("chatID отсутствует")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		return 0, errors.New("chatID должен быть положительным числом")
	}
	return id, nil
}

// parseLimit парсит параметр limit или возвращает дефолтное значение
func (er *ErrorHandler) parseLimit(r *http.Request) int {
	defaultLimit := 20
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return defaultLimit
	}
	if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
		return parsed
	}
	return defaultLimit
}

// decodeJSON декодирует тело запроса
func (er *ErrorHandler) decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// respondJSON отправляет стандартизированный JSON ответ
func (er *ErrorHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			er.apiLogger.Error("ошибка при кодировании ответа", zap.Error(err))
		}
	}
}
