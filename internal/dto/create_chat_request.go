package dto

import (
	"errors"
	"strings"
)

type CreateChatRequest struct {
	Title string `json:"title"`
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
