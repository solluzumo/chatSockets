package events

import (
	"chatsockets/internal/domain"
	"context"
)

type MessageCreatedEvent struct {
	Ctx  context.Context
	Data *domain.MessageDomain
}
