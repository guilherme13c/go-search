package set

type Set[T any] interface {
	Add(T)
	Find(T) bool
}
