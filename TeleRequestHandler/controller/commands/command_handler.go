package commands

type CommandHandler[T any] interface {
	Execute(T) error
}
