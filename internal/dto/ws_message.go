package dto

type WsErrorResponse struct {
	Type   string `json:"type"`
	Error  string `json:"error"`
	Reason string `json:"reason,omitempty"`
}
