package dto

import "time"

type GetMessageResponse struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}
