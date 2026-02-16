package domain

import "fmt"

type ChatStatus string

const (
	PrivateChat ChatStatus = "Приватный"
	PublicChat  ChatStatus = "Публичный"
	ChannelChan ChatStatus = "Канал"
)

func ParseChatStatus(status string) (ChatStatus, error) {
	switch status {
	case string(PrivateChat):
		return PrivateChat, nil
	case string(PublicChat):
		return PublicChat, nil
	case string(ChannelChan):
		return ChannelChan, nil
	default:
		return "", fmt.Errorf("неизвестная роль: %s", status)
	}
}
