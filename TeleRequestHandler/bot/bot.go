package bot

type Bot[T any, S any] interface {
	GetUpdates(timeout int) T
	Send(msg S) error
}
