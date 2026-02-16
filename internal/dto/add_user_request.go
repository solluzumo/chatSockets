package dto

import (
	"chatsockets/internal/domain"
	"fmt"
	"strconv"
)

type AddUserRequest struct {
	UserID string `json:"user_id"`
}

func (req *AddUserRequest) ToDomain() (*domain.UserChatLinkDomain, error) {

	uUserID, err := strconv.Atoi(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать user id в int: %w", domain.ErrBadRequest)
	}
	return &domain.UserChatLinkDomain{
		UserID: uUserID,
	}, nil
}
