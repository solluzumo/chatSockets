package domain

import "errors"

var (
	//УЖЕ СУЩЕСТВУЕТ
	ErrChatAlreadyExists = errors.New("чат уже существует")

	//НЕ СУЩЕСТВУЕТ
	ErrChatNotFound    = errors.New("чат не найден")
	ErrMessageNotFound = errors.New("сообщение не существует")
	ErrUserNotFound    = errors.New("пользователь не найден")

	// ОСТАЛЬНОЕ
	ErrFieldIsNotAllowed        = errors.New("не разрешенное для фильтрации поле")
	ErrUserIsNotConnectedToChat = errors.New("пользователь и чат не связаны")
	ErrBadJWT                   = errors.New("токен не валидный")
	ErrForeignKeyViolation      = errors.New("ошибка внешего ключа")
)
