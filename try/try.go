package try

import (
	"github.com/Sugther/go-structs/either"
	"github.com/Sugther/go-structs/equal"
	"github.com/Sugther/go-structs/option"
)

/*
Try is a container for a value of type T that may or may not have been successfully computed.
It contains an Either value with Left holding an error and Right holding a value of type T,
and a finallyFunction that is executed when the computation is done.
*/
type Try[T any] struct {
	either          either.Either[error, T]
	finallyFunction func()
}

/*
Success creates a Try value containing a successful computation result of type T.
Example: Success[int](42) returns Try{either: Either{Right: Option(42), Left: Option()}}
*/
func Success[T any](value T) Try[T] {
	return Try[T]{
		either:          either.Right[error, T](value),
		finallyFunction: func() {},
	}
}

/*
Fail creates a Try value containing a failed computation with an error.
Example: Fail[int](errors.New("error")) returns Try{either: Either{Right: Option(), Left: Option(error)}}
*/
func Fail[T any](err error) Try[T] {
	return Try[T]{
		either:          either.Left[error, T](err),
		finallyFunction: func() {},
	}
}

/*
Pure creates a Try value with a given value and an error.
If the error is nil, the Try value will contain a successful computation result.
If the error is not nil, the Try value will contain a failed computation with an error.
Example: Pure[int](42, nil) returns Success[int](42)
*/
func Pure[T any](value T, err error) Try[T] {
	if err == nil {
		return Success(value)
	}
	return Fail[T](err)
}

/*
IsSuccess checks if a Try value contains a successful computation result.
Returns true if the Try value contains a Right value, false otherwise.
Examples:
IsSuccess(Success[int](42)) returns true
IsSuccess(Fail[int](error)) returns false
*/
func IsSuccess[T any](try Try[T]) bool {
	return either.IsRight(try.either)
}

func (try Try[T]) IsSuccess() bool {
	return IsSuccess(try)
}

/*
IsFail checks if a Try value contains a failed computation with an error.
Returns true if the Try value contains a Left value, false otherwise.
Examples:
IsFail(Fail[int](error)) returns true
IsFail(Success[int](43)) returns false
*/
func IsFail[T any](try Try[T]) bool {
	return either.IsLeft(try.either)
}

func (try Try[T]) IsFail() bool {
	return IsFail(try)
}

/*
GetOrElse returns the successful computation result of a Try value if it exists,
or a default value if the Try value contains a failed computation.
Example:
GetOrElse(Success[int](42), 0) returns 42
GetOrElse(Fail[int](error), 0) returns 0
*/
func GetOrElse[T any](try Try[T], defaultValue T) T {
	return either.GetOrElse(try.either, defaultValue)
}

func (try Try[T]) GetOrElse(defaultValue T) T {
	return GetOrElse(try, defaultValue)
}

/*
Finally registers a function to be executed when the Try value is finalized by calling the End method.
The function f takes the successful computation result of type T as input.
The finallyFunction is stored in the Try value and is executed when End is called.
*/
func Finally[T any](try Try[T], f func(T)) Try[T] {
	return Try[T]{
		either: try.either,
		finallyFunction: func() {
			either.ForEach(try.either, f)
			End(try)
		},
	}
}

func (try Try[T]) Finally(f func(T)) Try[T] {
	return Finally(try, f)
}

/*
End executes the finallyFunction of a Try value and returns the Try value without the finallyFunction.
Example: End(Finally(Success[int](42), func(v int) { fmt.Println("Ended:", v) })) returns Success[int](42)
*/
func End[T any](try Try[T]) Try[T] {
	try.finallyFunction()
	return Try[T]{
		either:          try.either,
		finallyFunction: func() {},
	}
}

func (try Try[T]) End() Try[T] {
	return End(try)
}

/*
Fold applies fFail if the Try value contains a failed computation, or fSuccess if the Try value contains a successful computation.
Examples:
Fold(Success[int](42), func(err error) string { return "Error" }, func(value int) string { return strconv.Itoa(value) }) returns "42"
Fold(Fail[int](error), func(err error) string { return "Error" }, func(value int) string { return strconv.Itoa(value) }) returns "Error"
*/
func Fold[T any, R any](try Try[T], fFail func(error) R, fSuccess func(T) R) R {
	return either.Fold(try.either, fFail, fSuccess)
}

/*
FlatMap applies the function f to the successful computation result of a Try value, returning a new Try value of a different type.
Examples:
FlatMap(Success[int](2), func(value int) Try[string] { return Success[strconv.Itoa(value * 2)] }) returns Success[string]("4")
FlatMap(Fail[int](error), func(value int) Try[string] { return Success[strconv.Itoa(value * 2)] }) returns Fail[string](error)
*/
func FlatMap[T any, R any](try Try[T], f func(T) Try[R]) Try[R] {
	return Fold(try, func(err error) Try[R] {
		return Finally(Fail[R](err), func(ignored R) { try.finallyFunction() })
	}, func(t T) Try[R] {
		r := f(t)
		return Try[R]{
			either: r.either,
			finallyFunction: func() {
				End(r)
				End(try)
			},
		}
	})
}

/*
Map applies the function f to the successful computation result of a Try value, returning a new Try value of a different type.
Examples:
Map(Success[int](2), func(value int) string { return strconv.Itoa(value * 2) }) returns Success[string]("4")
Map(Fail[int](error), func(value int) string { return strconv.Itoa(value * 2) }) returns Fail[string](error)
*/
func Map[T any, R any](try Try[T], f func(T) R) Try[R] {
	return FlatMap(try, func(t T) Try[R] {
		return Success(f(t))
	})
}

/*
ForEach applies the function f to the successful computation result of a Try value.
Examples:
ForEach(Success[int](2), func(value int) { fmt.Println(value) }) prints "2"
ForEach(Fail[int](error), func(value int) { fmt.Println(value) }) does nothing
*/

func ForEach[T any](try Try[T], f func(T)) {
	try.either.ForEach(f)
}

func (try Try[T]) ForEach(f func(T)) {
	ForEach(try, f)
}

/*
IfFail applies the function f to the error of a failed computation in a Try value.
Examples:
IfFail(Fail[int](errors.New("error")), func(err error) { fmt.Println(err) }) prints "error"
IfFail(Success[int](20), func(err error) { fmt.Println(err) }) does nothing
*/
func IfFail[T any](try Try[T], f func(error)) {
	try.either.IfLeft(f)
}

func (try Try[T]) IfFail(f func(error)) {
	IfFail(try, f)
}

/*
FlatMapFail applies the function f to the error of a failed computation in a Try value, returning a new Try value of the same type.
Example:
FlatMapFail(Fail[int](errors.New("error")), func(err error) Try[int] { return Success[int](0) }) returns Success[int](0)
FlatMapFail(Success[int](50), func(err error) Try[int] { return Success[int](0) }) returns Success[int](50)
*/
func FlatMapFail[T any](try Try[T], f func(error) Try[T]) Try[T] {
	return Fold(try, func(err error) Try[T] {
		r := f(err)
		return Try[T]{
			either: r.either,
			finallyFunction: func() {
				End(r)
				End(try)
			},
		}
	}, func(t T) Try[T] { return try })
}

func (try Try[T]) FlatMapFail(f func(error) Try[T]) Try[T] {
	return FlatMapFail(try, f)
}

/*
MapLeft applies the function f to the error of a failed computation in a Try value, returning a new Try value with the transformed error.
Examples:
MapLeft(Fail[int](error), func(err error) error { return fmt.Errorf("wrapped: %w", err) }) returns Fail[int](wrapped: error)
MapLeft(Success[int](50), func(err error) error { return fmt.Errorf("wrapped: %w", err) }) returns Success[int](50)
*/
func MapLeft[T any](try Try[T], f func(err error) error) Try[T] {
	return FlatMapFail(try, func(e error) Try[T] {
		return Fail[T](f(e))
	})
}

/*
BiForEach applies fFail to the error of a failed computation in a Try value, and fSuccess to the successful computation result of a Try value.
Examples:
BiForEach(Success[int](42), func(err error) { fmt.Println("Error:", err) }, func(value int) { fmt.Println("Value:", value) }) prints "Value: 42"
BiForEach(Fail[int](error), func(err error) { fmt.Println("Error:", err) }, func(value int) { fmt.Println("Value:", value) }) prints "Error: error"
*/
func BiForEach[T any](try Try[T], fFail func(err error), fSuccess func(T)) {
	either.BiForEach(try.either, fFail, fSuccess)
}

func (try Try[T]) BiForEach(fFail func(err error), fSuccess func(T)) {
	BiForEach(try, fFail, fSuccess)
}

/*
ToOption converts a Try value to an Option value, discarding the error information of a failed computation.
Examples:
ToOption(Success[int](42)) returns Option[int]{value: 42, isDefined: true}
ToOption(Fail[int](error)) returns Option[int]{isDefined: false}
*/
func ToOption[T any](try Try[T]) option.Option[T] {
	return try.either.ToOption()
}

func (try Try[T]) ToOption() option.Option[T] {
	return ToOption(try)
}

/*
ToEither returns the Either value contained in a Try value.
Examples:
ToEither(Success[int](42)) returns Either{Right: Option(42), Left: Option()}
ToEither(Fail[int](error)) returns Either{Right: Option(), Left: Option(error)}
*/
func ToEither[T any](try Try[T]) either.Either[error, T] {
	return try.either
}

func (try Try[T]) ToEither() either.Either[error, T] {
	return try.either
}

func (try Try[T]) Equals(other interface{}) bool {
	if ot, ok := other.(Try[T]); ok {
		return equal.Equals(ot.either, try.either)
	}
	return false
}
