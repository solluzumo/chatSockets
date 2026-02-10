package dto

import (
	"chatsockets/internal/domain"
	"errors"
	"strconv"
	"strings"
)

type CreateMessageRequest struct {
	Text     string `json:"text"`
	UserName string `json:"username"`
	ChatID   string `json:"chatid"`
}

func (req *CreateMessageRequest) ToDomain() (*domain.MessageDomain, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	uChatId, err := strconv.Atoi(req.ChatID)
	if err != nil {
		return nil, err
	}

	return &domain.MessageDomain{
		Text:     req.Text,
		UserName: req.UserName,
		ChatID:   uChatId,
	}, nil
}

func (req *CreateMessageRequest) Validate() error {
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		return errors.New("text не может быть пустым")
	}
	if len(req.Text) > 5000 {
		return errors.New("text слишком длинный(5000 максимум)")
	}
	return nil
}
