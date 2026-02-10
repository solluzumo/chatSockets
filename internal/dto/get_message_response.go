package dto

import (
	"chatsockets/internal/domain"
	"encoding/json"
	"time"
)

type GetMessageResponse struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

func (res *GetMessageResponse) ToBytes() ([]byte, error) {
	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return bytes, nil

}

func ToGetMessageResponse(data *domain.MessageDomain) *GetMessageResponse {
	return &GetMessageResponse{
		ID:        data.ID,
		UserName:  data.UserName,
		Text:      data.Text,
		CreatedAt: data.CreatedAt,
	}
}

func ToGetMessageResponseSlice(data []*domain.MessageDomain) []GetMessageResponse {

	messagesResponse := make([]GetMessageResponse, len(data))

	for i, msg := range data {
		messagesResponse[i] = GetMessageResponse{
			ID:        msg.ID,
			Text:      msg.Text,
			UserName:  msg.UserName,
			CreatedAt: msg.CreatedAt,
		}
	}
	return messagesResponse
}
