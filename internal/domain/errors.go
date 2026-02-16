package domain

import "errors"

var (
	//ОШИБКИ ЗАПРОСА
	ErrBadRequest = errors.New("неправильный запрос")

	//УЖЕ СУЩЕСТВУЕТ
	ErrChatAlreadyExists     = errors.New("чат уже существует")
	ErrUserAlreadySubscribed = errors.New("пользователь уже состоит в чате")

	//НЕ СУЩЕСТВУЕТ
	ErrChatNotFound    = errors.New("чат не найден")
	ErrMessageNotFound = errors.New("сообщение не существует")
	ErrUserNotFound    = errors.New("пользователь не найден")

	//ОШИБКИ ДОСТУПА
	ErrUserBlocked              = errors.New("пользователь заблокирован")
	ErrNotEnoughStatus          = errors.New("недостаточно прав для этого действия")
	ErrUserIsNotConnectedToChat = errors.New("пользователь не состоит в чате")

	// ОСТАЛЬНОЕ
	ErrFieldIsNotAllowed   = errors.New("не разрешенное для фильтрации поле")
	ErrBadJWT              = errors.New("токен не валидный")
	ErrForeignKeyViolation = errors.New("ошибка внешего ключа")
	ErrNoChanges           = errors.New("изменений нет")
)
