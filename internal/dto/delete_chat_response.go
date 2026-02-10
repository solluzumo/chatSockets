package dto

type DeleteChatResponse struct {
	Content    string `json:"content"`
	StatusCode int    `json:"status_code"`
}

func ToDeleteChatResponse(c string, sc int) *DeleteChatResponse {
	return &DeleteChatResponse{
		Content:    c,
		StatusCode: sc,
	}
}
