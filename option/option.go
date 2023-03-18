package option

import "github.com/Sugther/go-structs/equal"

/*
Option represents an optional value container of type T.
*/
type Option[T any] struct {
	value   T
	isEmpty bool
}

/*
Pure creates a new Option containing the given value of type T.
Example: Pure(42) returns Option(42, false).
*/
func Pure[T any](value T) Option[T] {
	return Option[T]{
		value:   value,
		isEmpty: false,
	}
}

/*
Empty creates an empty Option for the given type T.
Example: Empty[int]() returns Option{isEmpty: true}.
*/
func Empty[T any]() Option[T] {
	return Option[T]{
		isEmpty: true,
	}
}

/*
Of creates an Option from a pointer to a value of type T.
If the pointer is not nil, it returns an Option containing the value.
Otherwise, it returns an empty Option.
Examples:
Of(&42) returns Option(42, false)
Of(nil) returns Option{isEmpty: true}.
*/
func Of[T any](value *T) Option[T] {
	if value != nil {
		return Pure(*value)
	}
	return Empty[T]()
}

/*
Get retrieves the value of type T stored within the Option.
It does not check if the Option is empty, so use with caution.
Example: opt.Get(Pure(42)) returns 42.
*/
func Get[T any](opt Option[T]) T {
	return opt.value
}

func (opt Option[T]) Get() T {
	return Get(opt)
}

/*
IsEmpty checks if the Option is empty.
Examples:
IsEmpty(Empty[int]()) returns true
IsEmpty(Pure(42)) returns false.
*/
func IsEmpty[T any](opt Option[T]) bool {
	return opt.isEmpty
}

func (opt Option[T]) IsEmpty() bool {
	return IsEmpty(opt)
}

/*
IsPresent checks if the Option contains a value.
Examples:
IsPresent(Empty[int]()) returns false
IsPresent(Pure(42)) returns true.
*/
func IsPresent[T any](opt Option[T]) bool {
	return !opt.isEmpty
}

func (opt Option[T]) IsPresent() bool {
	return IsPresent(opt)
}

/*
GetOrElse retrieves the value within the Option if present,
or returns the provided default value if the Option is empty.
Examples:
GetOrElse(Empty[int](), 42) returns 42
GetOrElse(Pure(1), 42) returns 1.
*/
func GetOrElse[T any](opt Option[T], defaultValue T) T {
	return Fold(opt, func() T { return defaultValue }, func(t T) T { return t })
}

func (opt Option[T]) GetOrElse(defaultValue T) T {
	return GetOrElse(opt, defaultValue)
}

/*
OrElse returns the Option if it contains a value, or returns the provided default Option if the original Option is empty.
Examples:
OrElse(Empty[int](), Pure(42)) returns Option(42, false)
OrElse(Pure(1), Pure(42)) returns Option(1, false).
*/
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

/*
Fold applies one of two functions depending on the state of the Option.
If the Option is empty, it applies the fEmpty function.
If the Option contains a value, it applies the fPresent function to the value.
Examples:
Fold(Empty[int](), func() int { return 42 }, func(x int) int { return x * 2 }) returns 42
Fold(Pure(1), func() int { return 42 }, func(x int) int { return x * 2 }) returns 2.
*/
func Fold[T any, R any](opt Option[T], fEmpty func() R, fPresent func(T) R) R {
	if opt.isEmpty {
		return fEmpty()
	}
	return fPresent(opt.value)
}

/*
Map applies a given function f to the value stored in the Option and returns a new Option of type R.
The function f should accept a value of type T and return a value of type R.
If the original Option is empty, it returns an empty Option[R].
Examples:
Map(Empty[int](), func(x int) int { return x * 2 }) returns Option{isEmpty: true}
Map(Pure(1), func(x int) int { return x * 2 }) returns Option(2, false).
*/
func Map[T any, R any](opt Option[T], f func(T) R) Option[R] {
	return Fold(opt, Empty[R], func(t T) Option[R] { return Pure(f(t)) })
}

/*
FlatMap applies a given function f to the value stored in the Option and returns a new Option of type R.
The function f should accept a value of type T and return an Option[R].
If the original Option is empty, it returns an empty Option[R].
Examples:
FlatMap(Empty[int](), func(x int) Option[int] { return Pure(x * 2) }) returns Option{isEmpty: true}
FlatMap(Pure(1), func(x int) Option[int] { return Pure(x * 2) }) returns Option(2, false).
*/
func FlatMap[T any, R any](opt Option[T], f func(T) Option[R]) Option[R] {
	return Fold(opt, Empty[R], f)
}

/*
ForEach applies a given function f to the value stored in the Option for its side effects.
If the Option is empty, the function is not called.
Examples:
ForEach(Empty[int](), func(x int) { fmt.Println(x) }) does not print
ForEach(Pure(1), func(x int) { fmt.Println(x) }) prints "1".
*/
func ForEach[T any](opt Option[T], f func(T)) {
	if IsPresent(opt) {
		f(opt.value)
	}
}

func (opt Option[T]) ForEach(f func(T)) {
	ForEach(opt, f)
}

/*
BiForEach applies one of two functions depending on the state of the Option.
If the Option is empty, it applies the fEmpty function.
If the Option contains a value, it applies the fPresent function to the value.
Examples:
BiForEach(Empty[int](), func() { fmt.Println("empty") }, func(x int) { fmt.Println(x) }) prints "empty"
BiForEach(Pure(1), func() { fmt.Println("empty") }, func(x int) { fmt.Println(x) }) prints "1".
*/
func BiForEach[T any](opt Option[T], fEmpty func(), fPresent func(T)) {
	if IsPresent(opt) {
		fPresent(opt.value)
	} else {
		fEmpty()
	}
}

func (opt Option[T]) BiForEach(fEmpty func(), fPresent func(T)) {
	BiForEach(opt, fEmpty, fPresent)
}

/*
IfEmpty applies a given function f if the Option is empty.
Examples:
IfEmpty(Empty[int](), func() { fmt.Println("empty") }) prints "empty"
IfEmpty(Pure(1), func() { fmt.Println("empty") }) does not print.
*/
func IfEmpty[T any](opt Option[T], f func()) {
	if opt.isEmpty {
		f()
	}
}

func (opt Option[T]) IfEmpty(f func()) {
	IfEmpty(opt, f)
}

/*
Filter applies a predicate function to the value stored in the Option and
returns a new Option. If the predicate function returns true, the new Option
contains the same value as the original. If the predicate returns false or the
original Option is empty, the new Option is empty.
Examples:
Filter(Empty[int](), func(x int) bool { return x > 0 }) returns Option{isEmpty: true}
Filter(Pure(1), func(x int) bool { return x > 0 }) returns Option(1, false)
Filter(Pure(-1), func(x int) bool { return x > 0 }) returns Option{isEmpty: true}.
*/
func Filter[T any](opt Option[T], f func(T) bool) Option[T] {
	return FlatMap(opt, func(value T) Option[T] {
		if f(value) {
			return Pure(value)
		}
		return Empty[T]()
	})
}

func (opt Option[T]) Filter(predicate func(T) bool) Option[T] {
	return Filter(opt, predicate)
}

func (opt Option[T]) Equals(other interface{}) bool {
	if oo, ok := other.(Option[T]); ok {
		return (oo.isEmpty && opt.isEmpty) || equal.Equals(oo.value, oo.value)
	}
	return false
}
