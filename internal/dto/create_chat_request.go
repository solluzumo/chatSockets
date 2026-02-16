package dto

import (
	"chatsockets/internal/domain"
	"errors"
	"fmt"
	"strings"
)

type CreateChatRequest struct {
	Title      string `json:"title"`
	ChatStatus string `json:"chat_status"`
}

func (req *CreateChatRequest) ToDomain() (*domain.ChatDomain, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	parsedStatus, err := domain.ParseChatStatus(req.ChatStatus)
	if err != nil {
		return nil, fmt.Errorf("неправильно указана роль: %w", err)
	}
	return &domain.ChatDomain{
		Title:      req.Title,
		ChatStatus: parsedStatus,
	}, nil
}

func (req *CreateChatRequest) Validate() error {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return errors.New("title не может быть пустым")
	}
	if len(req.Title) > 200 {
		return errors.New("title слишком длинный(200 максимум)")
	}
	return nil
}
