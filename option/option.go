package option

type Option[T any] struct {
	value   T
	isEmpty bool
}

func Pure[T any](value T) Option[T] {
	return Option[T]{
		value:   value,
		isEmpty: false,
	}
}

func Empty[T any]() Option[T] {
	return Option[T]{
		isEmpty: true,
	}
}

func Of[T any](value *T) Option[T] {
	if value != nil {
		return Pure(*value)
	}
	return Empty[T]()
}

func (opt Option[T]) Get() T {
	return opt.value
}

func IsEmpty[T any](opt Option[T]) bool {
	return opt.isEmpty
}

func (opt Option[T]) IsEmpty() bool {
	return IsEmpty(opt)
}

func IsPresent[T any](opt Option[T]) bool {
	return !opt.isEmpty
}

func (opt Option[T]) IsPresent() bool {
	return IsPresent(opt)
}

func GetOrElse[T any](opt Option[T], defaultValue T) T {
	return Fold(opt, func() T { return defaultValue }, func(t T) T { return t })
}

func (opt Option[T]) GetOrElse(defaultValue T) T {
	return GetOrElse(opt, defaultValue)
}

func OrElse[T any](opt Option[T], defaultValue Option[T]) Option[T] {
	return Fold(opt, func() Option[T] {
		return defaultValue
	}, func(t T) Option[T] {
		return Pure(t)
	})
}

func (opt Option[T]) OrElse(defaultValue Option[T]) Option[T] {
	return OrElse[T](opt, defaultValue)
}

func Fold[T any, R any](opt Option[T], fEmpty func() R, fPresent func(T) R) R {
	if opt.isEmpty {
		return fEmpty()
	}
	return fPresent(opt.value)
}

func Map[T any, R any](opt Option[T], f func(T) R) Option[R] {
	return Fold(opt, Empty[R], func(t T) Option[R] { return Pure(f(t)) })
}

func FlatMap[T any, R any](opt Option[T], f func(T) Option[R]) Option[R] {
	return Fold(opt, Empty[R], f)
}

func ForEach[T any](opt Option[T], f func(T)) {
	if IsPresent(opt) {
		f(opt.value)
	}
}

func (opt Option[T]) ForEach(f func(T)) {
	ForEach(opt, f)
}

func IfEmpty[T any](opt Option[T], f func()) {
	if opt.isEmpty {
		f()
	}
}

func (opt Option[T]) IfEmpty(f func()) {
	IfEmpty(opt, f)
}
