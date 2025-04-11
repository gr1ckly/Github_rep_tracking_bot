package state_machine

type State[T any] interface {
	Start() error
	Handle(T) error
}
