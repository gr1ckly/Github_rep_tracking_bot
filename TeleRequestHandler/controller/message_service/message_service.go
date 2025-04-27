package message_service

import "context"

type MessageService[T any] interface {
	ProcessMessages(ctx context.Context, ch T)
}
