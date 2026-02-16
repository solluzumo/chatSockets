package events

import (
	"chatsockets/internal/domain"
	"context"
)

type EventBus struct {
	MessageCreated chan MessageCreatedEvent
}

func NewEventBus() *EventBus {
	return &EventBus{
		MessageCreated: make(chan MessageCreatedEvent, 100),
	}
}

func (e *EventBus) PublicTask(ctx context.Context, data *domain.MessageDomain) {
	messageTask := MessageCreatedEvent{
		Ctx:  ctx,
		Data: data,
	}
	select {
	case e.MessageCreated <- messageTask:
	case <-ctx.Done():
		return
	}
}
