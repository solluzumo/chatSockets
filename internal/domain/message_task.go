package domain

import "context"

type MessageTask struct {
	Ctx  context.Context
	Data []byte
}
