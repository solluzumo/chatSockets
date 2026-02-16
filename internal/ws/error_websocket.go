package websockets

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/dto"
	"errors"

	"go.uber.org/zap"
)

type ErrorHandlerWS struct {
	apiLogger *zap.Logger
}

func NewErrorHandlerWS(logger *zap.Logger) *ErrorHandlerWS {
	return &ErrorHandlerWS{
		apiLogger: logger,
	}
}

func (er *ErrorHandlerWS) handleDomainErrorWS(send chan<- any, err error) {
	var errString string
	switch {
	case errors.Is(err, domain.ErrChatNotFound):
		errString = "Чат не найден"
	case errors.Is(err, domain.ErrUserIsNotConnectedToChat):
		errString = "пользователь не участник чата"
	case errors.Is(err, domain.ErrBadJWT):
		errString = "неверный jwt токен"
	case errors.Is(err, domain.ErrNotEnoughStatus):
		errString = "пользователь не имеет достаточно прав"
	default:
		errString = "Внутренняя ошибка сервера"
	}
	er.respondErrorWS(send, errString)
}

func (er *ErrorHandlerWS) respondErrorWS(send chan<- any, err string) {
	errorMsg := dto.WsErrorResponse{
		Type:   "error",
		Error:  "Не удалось подключиться к чату",
		Reason: err,
	}
	send <- errorMsg

}
