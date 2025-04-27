package notification_service

import "context"

type NotificationWaiter interface {
	WaitNotification(ctx context.Context) error
	Close() error
}
