package scrapper

type Observer interface {
	Notify()
}

type Observable interface {
	AddObserver(Observer)
	RemoveObserver(Observer)
}
