package dto

import (
	"chatsockets/internal/domain"

	"fmt"
	"strconv"
)

type UpdateUserRoleRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func (req *UpdateUserRoleRequest) ToDomain(chatID int) (*domain.UserChatLinkDomain, error) {
	uUserID, err := strconv.Atoi(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить user id в int:%w", err)
	}

	if err := req.Validate(uUserID); err != nil {
		return nil, fmt.Errorf("ошибка user id:%w", err)
	}

	parsedRole, err := domain.ParseUserRole(req.Role)
	if err != nil {
		return nil, fmt.Errorf("неправильно указана роль: %w", err)
	}

	return &domain.UserChatLinkDomain{
		ChatID:   chatID,
		UserID:   uUserID,
		UserRole: parsedRole,
	}, nil
}

func (req *UpdateUserRoleRequest) Validate(uUserID int) error {
	if uUserID < 0 {
		return fmt.Errorf("id не может быть отрицательным")
	}
	return nil
}
