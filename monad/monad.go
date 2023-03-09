package monad

type Monad[T any] struct {
	value T
}

func Pure[T any](value T) Monad[T] {
	return Monad[T]{
		value: value,
	}
}

func Get[T any](monad Monad[T]) T {
	return monad.value
}

func (monad Monad[T]) Get() T {
	return Get(monad)
}

func Fold[T any, R any](monad Monad[T], f func(T) R) R {
	return f(monad.value)
}

func FlatMap[T any, R any](monad Monad[T], f func(T) Monad[R]) Monad[R] {
	return Fold(monad, f)
}

func Map[T any, R any](monad Monad[T], f func(T) R) Monad[R] {
	return FlatMap(monad, func(t T) Monad[R] {
		return Pure(f(t))
	})
}

func ForEach[T any](monad Monad[T], f func(T)) {
	f(monad.value)
}

func (monad Monad[T]) ForEach(f func(T)) {
	ForEach(monad, f)
}
