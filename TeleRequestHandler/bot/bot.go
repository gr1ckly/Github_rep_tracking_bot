package bot

type Bot[T any, S any] interface {
	GetUpdatesChannel(timeout int) T
	SendMessage(msg S) error
}
