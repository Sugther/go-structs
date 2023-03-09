package list

import (
	"go-structs/option"
)

type List[T any] struct {
	values []T
}

func Pure[T any](values []T) List[T] {
	return List[T]{
		values: values,
	}
}

func Of[T any](values ...T) List[T] {
	return Pure(values)
}

func Empty[T any]() List[T] {
	return Pure([]T{})
}

func Len[T any](list List[T]) int {
	return len(list.values)
}

func (list List[T]) Len() int {
	return Len(list)
}

func IsEmpty[T any](list List[T]) bool {
	return Len(list) == 0
}

func (list List[T]) IsEmpty() bool {
	return IsEmpty(list)
}

func NonEmpty[T any](list List[T]) bool {
	return !IsEmpty(list)
}

func (list List[T]) NonEmpty() bool {
	return NonEmpty(list)
}

func Append[T any](list List[T], values ...T) List[T] {
	return Pure(append(list.values, values...))
}

func (list List[T]) Append(values ...T) List[T] {
	return Append(list, values...)
}

func AppendList[T any](list List[T], list2 List[T]) List[T] {
	return Append(list, list2.values...)
}

func (list List[T]) AppendList(l List[T]) List[T] {
	return AppendList(list, l)
}

func Head[T any](list List[T]) option.Option[T] {
	if IsEmpty(list) {
		return option.Empty[T]()
	}
	return option.Pure(list.values[0])
}

func (list List[T]) Head() option.Option[T] {
	return Head(list)
}

func Tail[T any](list List[T]) List[T] {
	return Pure(list.values[1:])
}

func (list List[T]) Tail() List[T] {
	return Tail(list)
}

func Fold[T any, R any](list List[T], root R, f func(R, T) R) R {
	if IsEmpty(list) {
		return root
	}
	return Fold(Tail(list), f(root, Head(list).Get()), f)
}

func FlatMap[T any, R any](list List[T], f func(T) List[R]) List[R] {
	return Fold(list, Empty[R](), func(l List[R], t T) List[R] {
		return AppendList(l, f(t))
	})
}

func Flatten[T any](list List[List[T]]) List[T] {
	return FlatMap(list, func(t List[T]) List[T] { return t })
}

func Map[T any, R any](list List[T], f func(T) R) List[R] {
	return FlatMap(list, func(t T) List[R] { return Of(f(t)) })
}

func Filter[T any](list List[T], f func(T) bool) List[T] {
	return Fold(list, Empty[T](), func(l List[T], t T) List[T] {
		if f(t) {
			return Append(l, t)
		}
		return l
	})
}

func (list List[T]) Filter(f func(T) bool) List[T] {
	return Filter(list, f)
}

func Find[T any](list List[T], f func(T) bool) option.Option[T] {
	for _, v := range list.values {
		if f(v) {
			return option.Pure(v)
		}
	}
	return option.Empty[T]()
}

func (list List[T]) Find(f func(T) bool) option.Option[T] {
	return Find(list, f)
}

func AnyMatch[T any](list List[T], f func(T) bool) bool {
	return Find(list, f).IsPresent()
}

func (list List[T]) AnyMatch(f func(T) bool) bool {
	return AnyMatch(list, f)
}

func Forall[T any](list List[T], f func(T) bool) bool {
	return !AnyMatch(list, func(t T) bool { return !f(t) })
}

func (list List[T]) Forall(f func(T) bool) bool {
	return Forall(list, f)
}

func Remove[T any](list List[T], f func(T) bool) List[T] {
	return Filter(list, func(t T) bool { return !f(t) })
}

func (list List[T]) Remove(f func(T) bool) List[T] {
	return Remove(list, f)
}

func Copy[T any](list List[T]) List[T] {
	return Pure(list.values)
}

func (list List[T]) Copy() List[T] {
	return Copy(list)
}

func Sort[T any](list List[T], isInOrder func(T, T) bool) List[T] {
	values := list.values
	for i := 0; i < len(values)-1; i++ {
		for j := 0; j < len(values)-i-1; j++ {
			if isInOrder(values[j+1], values[j]) {
				values[j], values[j+1] = values[j+1], values[j]
			}
		}
	}
	return Pure(values)
}

func (list List[T]) Sort(isInOrder func(T, T) bool) List[T] {
	return Sort(list, isInOrder)
}
