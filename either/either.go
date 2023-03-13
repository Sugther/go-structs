package either

import (
	"github.com/Sugther/go-structs/option"
)

type Either[L any, R any] struct {
	Right option.Option[R]
	Left  option.Option[L]
}

func Right[L any, R any](r R) Either[L, R] {
	return Either[L, R]{
		Right: option.Pure(r),
		Left:  option.Empty[L](),
	}
}

func Left[L any, R any](l L) Either[L, R] {
	return Either[L, R]{
		Right: option.Empty[R](),
		Left:  option.Pure(l),
	}
}

func IsRight[L any, R any](either Either[L, R]) bool {
	return !either.Right.IsEmpty()
}

func (either Either[L, R]) IsRight() bool {
	return IsRight(either)
}

func IsLeft[L any, R any](either Either[L, R]) bool {
	return !either.Left.IsEmpty()
}

func (either Either[L, R]) IsLeft() bool {
	return IsLeft(either)
}

func GetOrElse[L any, R any](either Either[L, R], defaultValue R) R {
	return either.Right.GetOrElse(defaultValue)
}

func (either Either[L, R]) GetOrElse(defaultValue R) R {
	return GetOrElse(either, defaultValue)
}

func Fold[L any, R any, T any](either Either[L, R], fLeft func(L) T, fRight func(R) T) T {
	if IsRight(either) {
		return fRight(either.Right.Get())
	}
	return fLeft(either.Left.Get())
}

func FlatMap[L any, R any, T any](either Either[L, R], f func(R) Either[L, T]) Either[L, T] {
	return Fold(either, Left[L, T], f)
}

func Map[L any, R any, T any](either Either[L, R], f func(R) T) Either[L, T] {
	return FlatMap(either, func(r R) Either[L, T] { return Right[L, T](f(r)) })
}

func ForEach[L any, R any](either Either[L, R], f func(R)) {
	if IsRight(either) {
		f(either.Right.Get())
	}
}

func (either Either[L, R]) ForEach(f func(R)) {
	ForEach(either, f)
}

func IfLeft[L any, R any](either Either[L, R], f func(L)) {
	if IsLeft(either) {
		f(either.Left.Get())
	}
}

func (either Either[L, R]) IfLeft(f func(L)) {
	IfLeft(either, f)
}

func FlatMapLeft[L any, R any, T any](either Either[L, R], f func(L) Either[T, R]) Either[T, R] {
	return Fold(either, f, Right[T, R])
}

func MapLeft[L any, R any, T any](either Either[L, R], f func(L) T) Either[T, R] {
	return FlatMapLeft(either, func(l L) Either[T, R] { return Left[T, R](f(l)) })
}

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

func ToOption[L any, R any](either Either[L, R]) option.Option[R] {
	return either.Right
}

func (either Either[L, R]) ToOption() option.Option[R] {
	return ToOption(either)
}
