package set

import (
	"github.com/Sugther/go-structs/list"
	"github.com/Sugther/go-structs/option"
)

/*
Set is a generic struct representing a list of unique values of type T.
The elements has to be Equal or comparable
It doesn't necessarily have to respect order to be equal.
*/
type Set[T any] struct {
	list list.List[T]
}

func pure[T any](values []T) Set[T] {
	l := list.Pure(values)
	return pureList(l)
}

func pureList[T any](l list.List[T]) Set[T] {
	return Set[T]{
		list: l,
	}
}

/*
Pure creates a new Set containing the given values.
Example: Pure([]int{1, 2, 3, 3}) returns Set[int]([1,2,3])
*/
func Pure[T any](values []T) Set[T] {
	l := list.Pure(values).Distinct()
	return pureList(l)
}

/*
Of creates a new Seq containing the given values.
Example: Of(1, 2, 3, 3) returns List[int]([1,2,3])
*/
func Of[T any](values ...T) Set[T] {
	return Pure(values)
}

/*
Distinct returns a new Set with all duplicate elements removed from the input List.
It uses the Equals method of the elements in the List to compare for equality.
Example:
Distinct(Of(1, 2, 3, 3)) returns List[int]([1, 2, 3])
*/
func Distinct[T any](l list.List[T]) Set[T] {
	return Pure(l.ToArray())
}

/*
Empty creates a new empty Set.
Example: Empty[int]() returns Set[int]([])
*/
func Empty[T any]() Set[T] {
	return pure([]T{})
}

/*
Len returns the length of the given set.
Example: Len(Of(1, 2, 3, 3)) returns 3
*/
func Len[T any](set Set[T]) int {
	return list.Len(set.list)
}

func (set Set[T]) Len() int {
	return Len(set)
}

/*
IsEmpty returns true if the given set is empty, false otherwise.
Examples:
IsEmpty(Of(1, 2, 3)) returns false
IsEmpty(Empty[int)()) returns true
*/
func IsEmpty[T any](set Set[T]) bool {
	return list.IsEmpty(set.list)
}

func (set Set[T]) IsEmpty() bool {
	return IsEmpty(set)
}

/*
NonEmpty returns true if the given set is not empty, false otherwise.
Examples:
NonEmpty(Of(1, 2, 3)) returns true
NonEmpty(Empty[int]()) returns false
*/
func NonEmpty[T any](set Set[T]) bool {
	return list.NonEmpty(set.list)
}

func (set Set[T]) NonEmpty() bool {
	return NonEmpty(set)
}

/*
Append returns a new Set with the given values appended to the original set.
Example: Append(Of(1, 2, 3), 3, 4, 5) returns List[int]([1,2,3,4,5])
*/
func Append[T any](set Set[T], values ...T) Set[T] {
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

/*
AppendSet returns a new Set with the elements of set2 appended to the original set.
Example: AppendSet(Of(1, 2, 3), Of(3, 4, 5, 6)) returns List[int]([1,2,3,4,5,6])
*/
func AppendSet[T any](set Set[T], set2 Set[T]) Set[T] {
	return AppendList(set, set2.list)
}

func (set Set[T]) AppendSet(set2 Set[T]) Set[T] {
	return AppendSet(set, set2)
}

/*
AppendList returns a new Set with the elements of list appended to the original set.
Example: AppendList(Of(1, 2, 3), Of(3, 4, 5, 6)) returns List[int]([1,2,3,4,5,6])
*/
func AppendList[T any](set Set[T], l list.List[T]) Set[T] {
	return Append(set, l.ToArray()...)
}

func (set Set[T]) AppendList(l list.List[T]) Set[T] {
	return AppendList(set, l)
}

/*
Head returns the first element of the set wrapped in an Option.
If the set is empty, it returns an empty Option.
Examples:
Head(Of(1, 2, 3)) returns Option[int](1)
Head(Empty[int]()) returns Option[int]{isEmpty: true}
*/
func Head[T any](set Set[T]) option.Option[T] {
	return list.Head(set.list)
}

func (set Set[T]) Head() option.Option[T] {
	return Head(set)
}

/*
Tail returns a new Set containing all elements except the first one.
Example: Tail(Of(1, 1, 2, 3)) returns List[int]([2,3])
*/
func Tail[T any](set Set[T]) Set[T] {
	return pureList(list.Tail(set.list))
}

func (set Set[T]) Tail() Set[T] {
	return Tail(set)
}

/*
Fold applies a function to the elements of the set in a cumulative way, starting from the given root value.
Examples:
Fold(Of(1, 2, 3, 4), 0, func(a int, b int) int { return a + b }) returns 10
Fold(Empty[int](), 10, func(a int, b int) int { return a + b }) returns 10
*/
func Fold[T any, R any](set Set[T], root R, f func(R, T) R) R {
	return list.Fold(set.list, root, f)
}

/*
FlatMap applies a function that returns a Set for each element of the input set, then concatenates the resulting sets.
Example:
FlatMap(Of(1, 2, 3), func(n int) Set[int] { return Of(n, n * 2) }) returns Set[int]([1,2,2,4,3,6])
FlatMap(Empty[int](), func(n int) Set[int] { return Of(n, n * 2) }) returns Set[int]([])
*/
func FlatMap[T any, R any](set Set[T], f func(T) Set[R]) Set[R] {
	l := list.FlatMap(set.list, func(t T) list.List[R] {
		return f(t).list
	})
	return Pure(l.ToArray())
}

/*
Map applies a function to each element of the set and returns a new Set with the results.
Examples:
Map(Of(1, 2, 3), func(n int) int { return n * n }) returns Set[int]([1,4,9])
Map(Empty[int](), func(n int) int { return n * n }) returns Set[int]([])
*/
func Map[T any, R any](set Set[T], f func(T) R) Set[R] {
	return FlatMap(set, func(t T) Set[R] { return pure([]R{f(t)}) })
}

/*
Filter returns a new Set containing only the elements that satisfy the given predicate function.
Example: Filter(Of(1, 2, 3, 4, 5), func(n int) bool { return n % 2 == 0 }) returns Set[int]([2,4])
*/
func Filter[T any](set Set[T], f func(T) bool) Set[T] {
	return pureList(list.Filter(set.list, f))
}

func (set Set[T]) Filter(f func(T) bool) Set[T] {
	return Filter(set, f)
}

/*
Find returns the first element that satisfies the given predicate function, wrapped in an Option.
If no elements satisfy the predicate, it returns an empty Option.
Examples:
Find(Of(1, 2, 3, 4, 5), func(n int) bool { return n % 2 == 0 }) returns Option[int](2)
Find(myList, func(n int) bool { return n < 0 }) returns Option[int]()
*/
func Find[T any](set Set[T], f func(T) bool) option.Option[T] {
	return list.Find(set.list, f)
}

func (set Set[T]) Find(f func(T) bool) option.Option[T] {
	return Find(set, f)
}

/*
AnyMatch returns true if any element in the set satisfies the given predicate function, false otherwise.
Examples:
AnyMatch(Of(1, 2, 3, 4, 5), func(n int) bool { return n % 2 == 0 }) returns true
AnyMatch(Of(1, 2, 3, 4, 5), func(n int) bool { return n < 0 }) returns false
*/
func AnyMatch[T any](set Set[T], f func(T) bool) bool {
	return list.AnyMatch(set.list, f)
}

func (set Set[T]) AnyMatch(f func(T) bool) bool {
	return AnyMatch(set, f)
}

/*
Forall returns true if all elements in the input Set satisfy the given predicate function, false otherwise.
Examples:
Forall(Of(1, 2, 3), func(n int) bool { return n > 0 }) returns true
Forall(Of(1, 2, 3), func(n int) bool { return n > 1 }) returns false
*/
func Forall[T any](set Set[T], f func(T) bool) bool {
	return list.Forall(set.list, f)
}

func (set Set[T]) Forall(f func(T) bool) bool {
	return Forall(set, f)
}

/*
Remove returns a new Set with all elements of the input Set that do not satisfy the given predicate function.
Examples:
Remove(Of(1, 2, 3, 4, 5), func(n int) bool { return n % 2 == 0 }) returns Set[int]([1,3,5])
Remove(Of(), func(n int) bool { return n % 2 == 0 }) returns Set[int]([])
*/
func Remove[T any](set Set[T], f func(T) bool) Set[T] {
	return pureList(list.Remove(set.list, f))
}

func (set Set[T]) Remove(f func(T) bool) Set[T] {
	return Remove(set, f)
}

/*
Copy returns a new Set with all elements of the input set copied.
Example:
Copy(Of(1, 2, 3)) returns Set[int]([1,2,3])
*/
func Copy[T any](set Set[T]) Set[T] {
	return pureList(list.Copy(set.list))
}

func (set Set[T]) Copy() Set[T] {
	return Copy(set)
}

/*
Sort returns a new Set with all elements of the input Set sorted according to the given comparison function.
Example:
Sort(Of(3, 1, 4, 1, 5, 9), func(a int, b int) bool { return a < b }) returns Set[int]([1,3,4,5,9])
*/
func Sort[T any](set Set[T], isInOrder func(T, T) bool) Set[T] {
	return pureList(list.Sort(set.list, isInOrder))
}

func (set Set[T]) Sort(isInOrder func(T, T) bool) Set[T] {
	return Sort(set, isInOrder)
}

/*
toList converts a Set to a list.List containing the same elements as the input set.
It maintains the order of elements in the original set.

Example:
toList(Of(1, 2, 3, 3)) returns list.List[int](1,2,3)
*/
func toList[T any](set Set[T]) list.List[T] {
	return set.list
}

func (set Set[T]) toList() list.List[T] {
	return toList(set)
}

/*
Contains returns true if the given value is present in the input Set, false otherwise.
Example:
Contains(Of[int](1, 2, 3), 3) returns true
Contains(Of[int](1, 2, 3), 4) returns false
*/
func Contains[T any](set Set[T], value T) bool {
	return list.Contains(set.list, value)
}

/*
Intersection returns a new Set containing the elements that are common between two input Sets.
Example:
set1 := Of[int](1, 2, 3)
set2 := Of[int](2, 3, 4)
Intersection(set1, set2) returns Set[int]([2,3])
*/
func Intersection[T any](set1 Set[T], set2 Set[T]) Set[T] {
	return Filter(set1, func(t T) bool {
		return Contains(set2, t)
	})
}

func (set Set[T]) Intersection(set2 Set[T]) Set[T] {
	return Intersection(set, set2)
}

/*
Difference returns a new Set containing the elements that are present in the first input Set but not in the second input Set.
Example:
set1 := Of[int](1, 2, 3)
set2 := Of[int](2, 3, 4)
Difference(set1, set2) returns Set[int]([1])
*/
func Difference[T any](set1 Set[T], set2 Set[T]) Set[T] {
	return Filter(set1, func(t T) bool {
		return !Contains(set2, t)
	})
}

func (set Set[T]) Difference(set2 Set[T]) Set[T] {
	return Difference(set, set2)
}

/*
Equals compares two Sets for equality by checking if all elements of the input set are present in the other set.
Returns true if both sets have the same elements, false otherwise.

Example:
set1 := Of[int](1, 2, 3)
set2 := Of[int](3, 2, 1)
set1.Equals(set2) returns true

set1 := Of[int](1, 2)
set2 := Of[int](1, 2, 3)
set1.Equals(set2) returns false
*/
func (set Set[T]) Equals(other interface{}) bool {
	if os, ok := other.(Set[T]); ok {
		if Len(set) != Len(set) {
			return false
		}
		return Fold(set, false, func(result bool, value T) bool {
			return result && Contains(os, value)
		})
	}
	return false
}
