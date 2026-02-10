package dto

// import (
// 	"chatsockets/internal/domain"
// 	"encoding/json"
// )

// type GetMessageResponseWS struct {
// 	ID       int    `json:"id"`
// 	UserName string `json:"username"`
// 	Text     string `json:"text"`
// }

// func (res *GetMessageResponseWS) ToBytes() ([]byte, error) {
// 	bytes, err := json.Marshal(res)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bytes, nil

// }

// func ToGetMessageResponseWS(data *domain.MessageDomain) *GetMessageResponseWS {
// 	return &GetMessageResponseWS{
// 		ID:       data.ID,
// 		UserName: data.UserName,
// 		Text:     data.Text,
// 	}
// }
