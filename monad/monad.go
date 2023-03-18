package monad

import (
	"github.com/Sugther/go-structs/equal"
)

/*
Monad represents a generic monadic container for a value of type T.
*/
type Monad[T any] struct {
	value T
}

/*
Pure creates a new Monad containing the given value of type T.
Example: Pure(42) returns Monad(42).
*/
func Pure[T any](value T) Monad[T] {
	return Monad[T]{
		value: value,
	}
}

/*
Get retrieves the value of type T stored within the Monad.
Example: Get(Monad(42)) returns 42.
*/
func Get[T any](monad Monad[T]) T {
	return monad.value
}

func (monad Monad[T]) Get() T {
	return Get(monad)
}

/*
Fold applies a given function f to the value stored in the Monad and returns a result of type R.
Example: Fold(Monad(1), func(x int) int { return x * 2 }) returns 2.
*/
func Fold[T any, R any](monad Monad[T], f func(T) R) R {
	return f(monad.value)
}

/*
FlatMap applies a given function f to the value stored in the Monad and returns a new Monad containing a value of type R.
The function f should accept a value of type T and return a Monad[R].
Example: FlatMap(Monad(1), func(x int) Monad[int] { return Pure(x * 2) }) returns Monad(2).
*/
func FlatMap[T any, R any](monad Monad[T], f func(T) Monad[R]) Monad[R] {
	return Fold(monad, f)
}

/*
Map applies a given function f to the value stored in the Monad and returns a new Monad containing a transformed value of type R.
The function f should accept a value of type T and return a value of type R.
Example: Map(Monad(1), func(x int) int { return x + 1 }) returns Monad(2).
*/
func Map[T any, R any](monad Monad[T], f func(T) R) Monad[R] {
	return FlatMap(monad, func(t T) Monad[R] {
		return Pure(f(t))
	})
}

/*
ForEach applies a given function f to the value stored in the Monad for its side effects.
The function f should accept a value of type T and return no value (void function).
Example: ForEach(Monad("Hello, World!"), func(s string) { fmt.Println(s) }) prints "Hello, World!".
*/
func ForEach[T any](monad Monad[T], f func(T)) {
	f(monad.value)
}

func (monad Monad[T]) ForEach(f func(T)) {
	ForEach(monad, f)
}

/*
Contains checks if the value stored in the Monad is equal to the given value.
The function returns true if the values are equal, false otherwise.
Example:
Contains(Monad(42), 42) returns true.
Contains(Monad(42), 43) returns false.
*/
func Contains[T any](monad Monad[T], value T) bool {
	return equal.Equals(monad.value, value)
}

func (monad Monad[T]) Contains(value T) bool {
	return Contains(monad, value)
}

func (monad Monad[T]) Equals(other interface{}) bool {
	if om, ok := other.(Monad[T]); ok {
		return equal.Equals(om.value, monad.value)
	}
	return false
}
