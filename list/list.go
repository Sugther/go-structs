package list

import (
	"github.com/Sugther/go-structs/equal"
	"github.com/Sugther/go-structs/option"
	"sort"
)

/*
List is a generic struct representing a list of values of type T.
*/
type List[T any] struct {
	values []T
}

/*
Pure creates a new List containing the given values.
Example: Pure([]int{1, 2, 3}) returns List[int]([1,2,3])
*/
func Pure[T any](values []T) List[T] {
	return List[T]{
		values: values,
	}
}

/*
Of creates a new List containing the given values.
Example: Of(1, 2, 3) returns List[int]([1,2,3])
*/
func Of[T any](values ...T) List[T] {
	return Pure(values)
}

/*
Empty creates a new empty List.
Example: Empty[int]() returns List[int]([])
*/
func Empty[T any]() List[T] {
	return Pure([]T{})
}

/*
Len returns the length of the given list.
Example: Len(Of(1, 2, 3)) returns 3
*/
func Len[T any](list List[T]) int {
	return len(list.values)
}

func (list List[T]) Len() int {
	return Len(list)
}

/*
IsEmpty returns true if the given list is empty, false otherwise.
Examples:
IsEmpty(Of(1, 2, 3)) returns false
IsEmpty(Empty[int)()) returns true
*/
func IsEmpty[T any](list List[T]) bool {
	return Len(list) == 0
}

func (list List[T]) IsEmpty() bool {
	return IsEmpty(list)
}

/*
NonEmpty returns true if the given list is not empty, false otherwise.
Examples:
NonEmpty(Of(1, 2, 3)) returns true
NonEmpty(Empty[int]()) returns false
*/
func NonEmpty[T any](list List[T]) bool {
	return !IsEmpty(list)
}

func (list List[T]) NonEmpty() bool {
	return NonEmpty(list)
}

/*
Append returns a new List with the given values appended to the original list.
Example: Append(Of(1, 2, 3), 4, 5) returns List[int]([1,2,3,4,5])
*/
func Append[T any](list List[T], values ...T) List[T] {
	return Pure(append(list.values, values...))
}

func (list List[T]) Append(values ...T) List[T] {
	return Append(list, values...)
}

/*
AppendList returns a new List with the elements of list2 appended to the original list.
Example: AppendList(Of(1, 2, 3), Of(4, 5, 6)) returns List[int]([1,2,3,4,5,6])
*/
func AppendList[T any](list List[T], list2 List[T]) List[T] {
	return Append(list, list2.values...)
}

func (list List[T]) AppendList(l List[T]) List[T] {
	return AppendList(list, l)
}

/*
Head returns the first element of the list wrapped in an Option.
If the list is empty, it returns an empty Option.
Examples:
Head(Of(1, 2, 3)) returns Option[int](1)
Head(Empty[int]()) returns Option[int]{isEmpty: true}
*/
func Head[T any](list List[T]) option.Option[T] {
	if IsEmpty(list) {
		return option.Empty[T]()
	}
	return option.Pure(list.values[0])
}

func (list List[T]) Head() option.Option[T] {
	return Head(list)
}

/*
Tail returns a new List containing all elements except the first one.
Example: Tail(Of(1, 2, 3)) returns List[int]([2,3])
*/
func Tail[T any](list List[T]) List[T] {
	return Pure(list.values[1:])
}

func (list List[T]) Tail() List[T] {
	return Tail(list)
}

/*
Fold applies a function to the elements of the list in a cumulative way, starting from the given root value.
Examples:
Fold(Of(1, 2, 3, 4), 0, func(a int, b int) int { return a + b }) returns 10
Fold(Empty[int](), 10, func(a int, b int) int { return a + b }) returns 10
*/
func Fold[T any, R any](list List[T], root R, f func(R, T) R) R {
	result := root
	values := list.values
	for i := 0; i < len(values)-1; i++ {
		result = f(result, values[i])
	}
	return result
}

/*
FlatMap applies a function that returns a List for each element of the input list, then concatenates the resulting lists.
Example:
FlatMap(Of(1, 2, 3), func(n int) List[int] { return Of(n, n * 2) }) returns List[int]([1,2,2,4,3,6])
FlatMap(Empty[int](), func(n int) List[int] { return Of(n, n * 2) }) returns List[int]([])
*/
func FlatMap[T any, R any](list List[T], f func(T) List[R]) List[R] {
	return Fold(list, Empty[R](), func(l List[R], t T) List[R] {
		return AppendList(l, f(t))
	})
}

/*
Flatten returns a new List with all elements of the input List flattened into a single List.
Example:
Flatten(Of(Of(1, 2), Of(3, 4))) returns List[int]([1,2,3,4])
*/
func Flatten[T any](list List[List[T]]) List[T] {
	return FlatMap(list, func(t List[T]) List[T] { return t })
}

/*
Map applies a function to each element of the list and returns a new List with the results.
Examples:
Map(Of(1, 2, 3), func(n int) int { return n * n }) returns List[int]([1,4,9])
Map(Empty[int](), func(n int) int { return n * n }) returns List[int]([])
*/
func Map[T any, R any](list List[T], f func(T) R) List[R] {
	return FlatMap(list, func(t T) List[R] { return Of(f(t)) })
}

/*
Filter returns a new List containing only the elements that satisfy the given predicate function.
Example: Filter(Of(1, 2, 3, 4, 5), func(n int) bool { return n % 2 == 0 }) returns List[int]([2,4])
*/
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

/*
Find returns the first element that satisfies the given predicate function, wrapped in an Option.
If no elements satisfy the predicate, it returns an empty Option.
Examples:
Find(Of(1, 2, 3, 4, 5), func(n int) bool { return n % 2 == 0 }) returns Option[int](2)
Find(myList, func(n int) bool { return n < 0 }) returns Option[int]()
*/
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

/*
AnyMatch returns true if any element in the list satisfies the given predicate function, false otherwise.
Examples:
AnyMatch(Of(1, 2, 3, 4, 5), func(n int) bool { return n % 2 == 0 }) returns true
AnyMatch(Of(1, 2, 3, 4, 5), func(n int) bool { return n < 0 }) returns false
*/
func AnyMatch[T any](list List[T], f func(T) bool) bool {
	return Find(list, f).IsPresent()
}

func (list List[T]) AnyMatch(f func(T) bool) bool {
	return AnyMatch(list, f)
}

/*
Forall returns true if all elements in the input List satisfy the given predicate function, false otherwise.
Examples:
Forall(Of(1, 2, 3), func(n int) bool { return n > 0 }) returns true
Forall(Of(1, 2, 3), func(n int) bool { return n > 1 }) returns false
*/
func Forall[T any](list List[T], f func(T) bool) bool {
	return !AnyMatch(list, func(t T) bool { return !f(t) })
}

func (list List[T]) Forall(f func(T) bool) bool {
	return Forall(list, f)
}

/*
Remove returns a new List with all elements of the input List that do not satisfy the given predicate function.
Examples:
Remove(Of(1, 2, 3, 4, 5), func(n int) bool { return n % 2 == 0 }) returns List[int]([1,3,5])
Remove(Of(), func(n int) bool { return n % 2 == 0 }) returns List[int]([])
*/
func Remove[T any](list List[T], f func(T) bool) List[T] {
	return Filter(list, func(t T) bool { return !f(t) })
}

func (list List[T]) Remove(f func(T) bool) List[T] {
	return Remove(list, f)
}

/*
Copy returns a new List with all elements of the input List copied.
Example:
Copy(Of(1, 2, 3)) returns List[int]([1,2,3])
*/
func Copy[T any](list List[T]) List[T] {
	copyValues := make([]T, len(list.values))
	copy(copyValues, list.values)
	return Pure(copyValues)
}

func (list List[T]) Copy() List[T] {
	return Copy(list)
}

/*
Sort returns a new List with all elements of the input List sorted according to the given comparison function.
Example:
Sort(Of(3, 1, 4, 1, 5, 9), func(a int, b int) bool { return a < b }) returns List[int]([1,1,3,4,5,9])
*/
func Sort[T any](list List[T], isInOrder func(T, T) bool) List[T] {
	copyValues := make([]T, len(list.values))
	copy(copyValues, list.values)
	sort.Slice(copyValues, func(i, j int) bool {
		return isInOrder(copyValues[i], copyValues[j])
	})
	return Pure(copyValues)
}

func (list List[T]) Sort(isInOrder func(T, T) bool) List[T] {
	return Sort(list, isInOrder)
}

/*
ToArray returns a new slice with all elements of the input List.
Example:
ToArray(Of(1, 2, 3)) returns []int{1, 2, 3}
*/
func ToArray[T any](list List[T]) []T {
	return list.values
}

func (list List[T]) ToArray() []T {
	return ToArray(list)
}

/*
Contains returns true if the given value is present in the input List, false otherwise.
It uses the Equals method of the elements in the List to compare for equality.
Example:
Contains(Of[int](1, 2, 3), 3) returns true
Contains(Of[int](1, 2, 3), 4) returns false
*/
func Contains[T any](list List[T], value T) bool {
	return AnyMatch(list, func(t T) bool {
		return equal.Equals(t, value)
	})
}

/*
Distinct returns a new List with all duplicate elements removed from the input List.
It uses the Equals method of the elements in the List to compare for equality.
Example:
Distinct(Of(1, 2, 3, 3)) returns List[int]([1, 2, 3])
*/
func Distinct[T any](list List[T]) List[T] {
	return Fold(list, Empty[T](), func(uniqueList List[T], value T) List[T] {
		if Contains(uniqueList, value) {
			return uniqueList
		}
		return Append(uniqueList, value)
	})
}

func (list List[T]) Distinct() List[T] {
	return Distinct(list)
}

/*
Intersection returns a new List containing the elements that are common between two input Lists.
Example:
list1 := Of[int](1, 2, 3)
list2 := Of[int](2, 3, 4)
Intersection(list1, list2) returns list[int]([2,3])
*/
func Intersection[T any](list1 List[T], list2 List[T]) List[T] {
	return Filter(list1, func(t T) bool {
		return Contains(list2, t)
	})
}

/*
Difference returns a new List containing the elements that are present in the first input List but not in the second input List.
Example:
list1 := Of[int](1, 2, 3)
list2 := Of[int](2, 3, 4)
Difference(list1, list2) returns List[int]([1])
*/
func Difference[T any](list1 List[T], list2 List[T]) List[T] {
	return Filter(list1, func(t T) bool {
		return !Contains(list2, t)
	})
}

/*
IsSublistOf returns true if all elements of the first input List are present in the second input List, false otherwise.
Example:
list1 := Of[int](1, 2, 3)
list2 := Of[int](1, 2, 3, 4)
IsSublistOf(list1, list2) returns true
*/
func IsSublistOf[T any](list1 List[T], list2 List[T]) bool {
	return Forall(list1, func(t T) bool {
		return Contains(list2, t)
	})
}

/*
IsSuperListOf returns true if all elements of the second input List are present in the first input List, false otherwise.
Example:
list1 := Of[int](1, 2, 3, 4)
list2 := Of[int](1, 2, 3)
IsSuperListOf(list1, list2) returns true
*/
func IsSuperListOf[T any](list1 List[T], list2 List[T]) bool {
	return IsSublistOf(list2, list1)
}

/*
Reverse returns a new list with the elements of the original list in reverse order.
Example: Reverse(Of(1, 2, 3)) returns List[int]([3,2,1])
*/
func Reverse[T any](list List[T]) List[T] {
	reversed := make([]T, len(list.values))
	for i, j := len(list.values)-1, 0; i >= 0; i, j = i-1, j+1 {
		reversed[j] = list.values[i]
	}
	return Pure(reversed)
}

func (list List[T]) Reverse() List[T] {
	return Reverse(list)
}

func (list List[T]) Equals(other interface{}) bool {
	if ol, ok := other.(List[T]); ok {
		if len(list.values) != len(ol.values) {
			return false
		}
		for i := range list.values {
			if equal.Equals(list.values[i], ol.values[i]) {
				return false
			}
		}
		return true
	}
	return false
}
