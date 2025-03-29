package LinkTracker

type Observer[T any] interface {
	Notify(T)
}

type Observable[T any] interface {
	NotifyAll(T)
	AddObserver(Observer[T])
	RemoveObserver(Observer[T])
}
