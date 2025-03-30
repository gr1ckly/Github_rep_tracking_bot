package LinkTracker

type Observer[T any] interface {
	Notify(T) error
}

type Observable[T any] interface {
	NotifyAll(T)
	AddObserver(Observer[T])
	RemoveObserver(Observer[T])
}
