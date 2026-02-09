package domain

import "errors"

var (
	ErrChatNotFound             = errors.New("чат не найден")
	ErrChatAlreadyExists        = errors.New("чат уже существует")
	ErrFieldIsNotAllowed        = errors.New("не разрешенное для фильтрации поле")
	ErrUserIsNotConnectedToChat = errors.New("пользователь и чат не связаны")
)
