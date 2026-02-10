package postgres

import (
	"chatsockets/internal/domain"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// если внешнего ключа не существует
// - отсутсвует чат для отправки сообщени
// - отсутствует пользователь отправляющий сообщение
// - и т.д.
func IsForeignKeyViolation(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil, false
	}

	if pgErr.Code != "23503" {
		return nil, false
	}

	return pgErr, true
}

func MapFKConstraint(constraint string) error {
	switch constraint {
	case "messages_chat_id_fkey":
		return domain.ErrChatNotFound
	case "messages_user_id_fkey":
		return domain.ErrUserNotFound
	default:
		return domain.ErrForeignKeyViolation
	}
}

// если нарушается уникальность поля
// - чаты с одинаковыми названиями
// - пользователи с одинаковыми юзернеймами
// - и т.д.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
