package set

import (
	"github.com/Sugther/go-structs/equal"
	"github.com/Sugther/go-structs/list"
	"github.com/Sugther/go-structs/option"
)

type Set[T equal.Equal] struct {
	list list.List[T]
}

func Pure[T equal.Equal](values []T) Set[T] {
	l := list.Pure(values)
	return Set[T]{
		list: list.Unique(l),
	}
}

func Of[T equal.Equal](values ...T) Set[T] {
	return Pure(values)
}

func Unique[T equal.Equal](l list.List[T]) Set[T] {
	return Pure(l.ToArray())
}

func Empty[T equal.Equal]() Set[T] {
	return Pure([]T{})
}

func Len[T equal.Equal](set Set[T]) int {
	return list.Len(set.list)
}

func (set Set[T]) Len() int {
	return Len(set)
}

func IsEmpty[T equal.Equal](set Set[T]) bool {
	return list.IsEmpty(set.list)
}

func (set Set[T]) IsEmpty() bool {
	return IsEmpty(set)
}

func NonEmpty[T equal.Equal](set Set[T]) bool {
	return list.NonEmpty(set.list)
}

func (set Set[T]) NonEmpty() bool {
	return NonEmpty(set)
}

func Append[T equal.Equal](set Set[T], values ...T) Set[T] {
	valuesList := list.Pure(values)
	return list.Fold(valuesList, set, func(result Set[T], value T) Set[T] {
		if list.Contains(result.list, value) {
			return result
		}
		return result.Append(value)
	})
}

func (set Set[T]) Append(values ...T) Set[T] {
	return Append(set, values...)
}

func AppendSet[T equal.Equal](set Set[T], set2 Set[T]) Set[T] {
	return AppendList(set, set2.list)
}

func (set Set[T]) AppendSet(set2 Set[T]) Set[T] {
	return AppendSet(set, set2)
}

func AppendList[T equal.Equal](set Set[T], l list.List[T]) Set[T] {
	return Append(set, l.ToArray()...)
}

func (set Set[T]) AppendList(l list.List[T]) Set[T] {
	return AppendList(set, l)
}

func Head[T equal.Equal](set Set[T]) option.Option[T] {
	return list.Head(set.list)
}

func (set Set[T]) Head() option.Option[T] {
	return Head(set)
}

func Fold[T equal.Equal, R equal.Equal](set Set[T], root R, f func(R, T) R) R {
	return list.Fold(set.list, root, f)
}

func FlatMap[T equal.Equal, R equal.Equal](set Set[T], f func(T) Set[R]) Set[R] {
	l := list.FlatMap(set.list, func(t T) list.List[R] {
		return f(t).list
	})
	return Pure(l.ToArray())
}

func Flatten[T equal.Equal](set Set[Set[T]]) Set[T] {
	return FlatMap(set, func(t Set[T]) Set[T] { return t })
}

func Map[T equal.Equal, R equal.Equal](set Set[T], f func(T) R) Set[R] {
	return FlatMap(set, func(t T) Set[R] { return Of(f(t)) })
}

func Filter[T equal.Equal](set Set[T], f func(T) bool) Set[T] {
	return Pure(list.Filter(set.list, f).ToArray())
}

func (set Set[T]) Filter(f func(T) bool) Set[T] {
	return Filter(set, f)
}

func Find[T equal.Equal](set Set[T], f func(T) bool) option.Option[T] {
	return list.Find(set.list, f)
}

func (set Set[T]) Find(f func(T) bool) option.Option[T] {
	return Find(set, f)
}

func AnyMatch[T equal.Equal](set Set[T], f func(T) bool) bool {
	return list.AnyMatch(set.list, f)
}

func (set Set[T]) AnyMatch(f func(T) bool) bool {
	return AnyMatch(set, f)
}

func Forall[T equal.Equal](set Set[T], f func(T) bool) bool {
	return list.Forall(set.list, f)
}

func (set Set[T]) Forall(f func(T) bool) bool {
	return Forall(set, f)
}

func Remove[T equal.Equal](set Set[T], f func(T) bool) Set[T] {
	return Set[T]{
		list: list.Remove(set.list, f),
	}
}

func (set Set[T]) Remove(f func(T) bool) Set[T] {
	return Remove(set, f)
}

func Copy[T equal.Equal](set Set[T]) Set[T] {
	return Set[T]{
		list: list.Copy(set.list),
	}
}

func (set Set[T]) Copy() Set[T] {
	return Copy(set)
}

func Sort[T equal.Equal](set Set[T], isInOrder func(T, T) bool) Set[T] {
	return Set[T]{
		list: list.Sort(set.list, isInOrder),
	}
}

func (set Set[T]) Sort(isInOrder func(T, T) bool) Set[T] {
	return Sort(set, isInOrder)
}

func toList[T equal.Equal](set Set[T]) list.List[T] {
	return set.list
}

func (set Set[T]) toList() list.List[T] {
	return toList(set)
}

func Contains[T equal.Equal](set Set[T], value T) bool {
	return list.Contains(set.list, value)
}

func (set Set[T]) Equals(other interface{}) bool {
	if o, ok := other.(Set[T]); ok {
		return Forall(set, func(t T) bool {
			return Contains(o, t)
		})
	}
	return false
}
