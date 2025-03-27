package validators

type Checker[T any, S any] interface {
	Check(T) S
}
