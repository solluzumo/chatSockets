package dto

import (
	"errors"
	"strings"
)

type CreateMessageRequest struct {
	Text string `json:"text"`
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
