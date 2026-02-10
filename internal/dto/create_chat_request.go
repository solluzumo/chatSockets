package dto

import (
	"chatsockets/internal/domain"
	"errors"
	"strings"
)

type CreateChatRequest struct {
	Title string `json:"title"`
}

func (req *CreateChatRequest) ToDomain() (*domain.ChatDomain, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &domain.ChatDomain{
		Title: req.Title,
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
