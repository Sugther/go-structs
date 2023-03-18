package either

import (
	"github.com/Sugther/go-structs/equal"
	"github.com/Sugther/go-structs/option"
)

/*
Either is a container for a value that can be one of two types Left or Right (L or R).
*/
type Either[L any, R any] struct {
	Right option.Option[R]
	Left  option.Option[L]
}

/*
Right creates an Either value containing a value of type R.
Example: Right[string, int](42) returns Either{Right: Option(42), Left: Option()}.
*/
func Right[L any, R any](r R) Either[L, R] {
	return Either[L, R]{
		Right: option.Pure(r),
		Left:  option.Empty[L](),
	}
}

/*
Left creates an Either value containing a value of type L.
Example: Left[string, int]("error") returns Either{Right: Option(), Left: Option("error")}.
*/
func Left[L any, R any](l L) Either[L, R] {
	return Either[L, R]{
		Right: option.Empty[R](),
		Left:  option.Pure(l),
	}
}

/*
IsRight checks if the Either value contains a value of type R (Right).
Examples:
IsRight(Right[string, int](42)) returns true
IsRight(Left[string, int]("error")) returns false.
*/
func IsRight[L any, R any](either Either[L, R]) bool {
	return !either.Right.IsEmpty()
}

func (either Either[L, R]) IsRight() bool {
	return IsRight(either)
}

/*
IsLeft checks if the Either value contains a value of type L (Left).
Examples:
IsLeft(Right[string, int](42)) returns false
IsLeft(Left[string, int]("error")) returns true.
*/
func IsLeft[L any, R any](either Either[L, R]) bool {
	return !either.Left.IsEmpty()
}

func (either Either[L, R]) IsLeft() bool {
	return IsLeft(either)
}

/*
GetOrElse retrieves the value of type R stored within the Either.
If the Either contains a Right value, it returns the value of type R.
If the Either contains a Left value, it returns the provided default value.
Example: GetOrElse(Left[string, int]("error"), 42) returns 42.
*/
func GetOrElse[L any, R any](either Either[L, R], defaultValue R) R {
	return either.Right.GetOrElse(defaultValue)
}

func (either Either[L, R]) GetOrElse(defaultValue R) R {
	return GetOrElse(either, defaultValue)
}

/*
Fold applies one of two functions depending on the state of the Either.
If the Either contains a Left value, it applies the fLeft function to the value.
If the Either contains a Right value, it applies the fRight function to the value.
Examples:
Fold(Left[string, int]("error"), func(l string) int { return len(l) }, func(r int) int { return r * 2 }) returns 5
Fold(Right[string, int](42), func(l string) int { return len(l) }, func(r int) int { return r * 2 }) returns 84.
*/
func Fold[L any, R any, T any](either Either[L, R], fLeft func(L) T, fRight func(R) T) T {
	if IsRight(either) {
		return fRight(either.Right.Get())
	}
	return fLeft(either.Left.Get())
}

/*
FlatMap applies a given function f to the Right value stored in the Either and returns a new Either of type L and T.
The function f should accept a value of type R and return an Either[L, T].
If the Either contains a Left value, it returns a new Either with the same Left value.
Examples:
FlatMap(Right[string, int](42), func(r int) Either[string, float64] { return Right[string, float64](float64(r) / 2) }) returns Right(21.0)
FlatMap(Left[string, int]("error"), func(r int) Either[string, float64] { return Right[string, float64](float64(r) / 2) }) returns Left("error").
*/
func FlatMap[L any, R any, T any](either Either[L, R], f func(R) Either[L, T]) Either[L, T] {
	return Fold(either, Left[L, T], f)
}

/*
Map applies a given function f to the Right value stored in the Either and returns a new Either of type L and T.
The function f should accept a value of type R and return a value of type T.
If the Either contains a Left value, it returns a new Either with the same Left value.
Examples:
Map(Right[string, int](42), func(r int) float64 { return float64(r) / 2 }) returns Right(21.0)
Map(Left[string, int]("error"), func(r int) float64 { return float64(r) / 2 }) returns Left("error").
*/
func Map[L any, R any, T any](either Either[L, R], f func(R) T) Either[L, T] {
	return FlatMap(either, func(r R) Either[L, T] { return Right[L, T](f(r)) })
}

/*
ForEach applies a given function f to the Right value stored in the Either if it exists.
The function f should accept a value of type R.
If the Either contains a Left value, the function does nothing.
Examples:
ForEach(Right[string, int](42), func(r int) { fmt.Println(r) }) prints 42
ForEach(Left[string, int]("error"), func(r int) { fmt.Println(r) }) does nothing.
*/
func ForEach[L any, R any](either Either[L, R], f func(R)) {
	if IsRight(either) {
		f(either.Right.Get())
	}
}

func (either Either[L, R]) ForEach(f func(R)) {
	ForEach(either, f)
}

/*
IfLeft applies a given function f to the Left value stored in the Either if it exists.
The function f should accept a value of type L.
If the Either contains a Right value, the function does nothing.
Examples:
IfLeft(Left[string, int]("error"), func(l string) { fmt.Println(l) }) prints "error"
IfLeft(Right[string, int](42), func(l string) { fmt.Println(l) }) does nothing.
*/
func IfLeft[L any, R any](either Either[L, R], f func(L)) {
	if IsLeft(either) {
		f(either.Left.Get())
	}
}

func (either Either[L, R]) IfLeft(f func(L)) {
	IfLeft(either, f)
}

/*
FlatMapLeft applies a given function f to the Left value stored in the Either and returns a new Either of type T and R.
The function f should accept a value of type L and return an Either[T, R].
If the Either contains a Right value, it returns a new Either with the same Right value.
Examples:
FlatMapLeft(Left[string, int]("error"), func(l string) Either[int, int] { return Left[int, int](len(l)) }) returns Left(5)
FlatMapLeft(Right[string, int](42), func(l string) Either[int, int] { return Left[int, int](len(l)) }) returns Right(42).
*/
func FlatMapLeft[L any, R any, T any](either Either[L, R], f func(L) Either[T, R]) Either[T, R] {
	return Fold(either, f, Right[T, R])
}

/*
MapLeft applies a given function f to the Left value stored in the Either and returns a new Either of type T and R.
The function f should accept a value of type L and return a value of type T.
If the Either contains a Right value, it returns a new Either with the same Right value.
Examples:
MapLeft(Left[string, int]("error"), func(l string) int { return len(l) }) returns Left(5)
MapLeft(Right[string, int](42), func(l string) int { return len(l) }) returns Right(42).
*/
func MapLeft[L any, R any, T any](either Either[L, R], f func(L) T) Either[T, R] {
	return FlatMapLeft(either, func(l L) Either[T, R] { return Left[T, R](f(l)) })
}

/*
BiForEach applies a given function fLeft to the Left value stored in the Either and a given function fRight to the Right value stored in the Either if they exist.
The functions fLeft and fRight should accept values of types L and R, respectively.
Examples:
BiForEach(Left[string, int]("error"), func(l string) { fmt.Println(l) }, func(r int) { fmt.Println(r) }) prints "error"
BiForEach(Right[string, int](42), func(l string) { fmt.Println(l) }, func(r int) { fmt.Println(r) }) prints 42.
*/
func BiForEach[L any, R any](either Either[L, R], fLeft func(L), fRight func(R)) {
	if IsRight(either) {
		fRight(either.Right.Get())
	} else {
		fLeft(either.Left.Get())
	}
}

func (either Either[L, R]) BiForEach(fLeft func(L), fRight func(R)) {
	BiForEach(either, fLeft, fRight)
}

/*
ToOption converts an Either value to an Option value containing the Right value if it exists.
Examples:
ToOption(Left[string, int]("error")) returns Option()
ToOption(Right[string, int](42)) returns Option(42).
*/
func ToOption[L any, R any](either Either[L, R]) option.Option[R] {
	return either.Right
}

func (either Either[L, R]) ToOption() option.Option[R] {
	return ToOption(either)
}

func (either Either[L, R]) Equals(other interface{}) bool {
	if oe, ok := other.(Either[L, R]); ok {
		return equal.Equals(either.Right, oe.Right) && equal.Equals(either.Left, oe.Left)
	}
	return false
}
