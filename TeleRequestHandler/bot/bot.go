package bot

type Bot[T any, S any, U any] interface {
	GetUpdates(timeout int) T
	SendMessage(id U, msg S) error
}
