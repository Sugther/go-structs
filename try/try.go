package try

import (
	"github.com/Sugther/go-structs/either"
	"github.com/Sugther/go-structs/option"
)

type Try[T any] struct {
	either          either.Either[error, T]
	finallyFunction func()
}

func Success[T any](value T) Try[T] {
	return Try[T]{
		either:          either.Right[error, T](value),
		finallyFunction: func() {},
	}
}

func Fail[T any](err error) Try[T] {
	return Try[T]{
		either:          either.Left[error, T](err),
		finallyFunction: func() {},
	}
}

func Pure[T any](value T, err error) Try[T] {
	if err == nil {
		return Success(value)
	}
	return Fail[T](err)
}

func IsSuccess[T any](try Try[T]) bool {
	return either.IsRight(try.either)
}

func (try Try[T]) IsSuccess() bool {
	return IsSuccess(try)
}

func IsFail[T any](try Try[T]) bool {
	return either.IsLeft(try.either)
}

func (try Try[T]) IsFail() bool {
	return IsFail(try)
}

func GetOrElse[T any](try Try[T], defaultValue T) T {
	return either.GetOrElse(try.either, defaultValue)
}

func (try Try[T]) GetOrElse(defaultValue T) T {
	return GetOrElse(try, defaultValue)
}

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

func Fold[T any, R any](try Try[T], fFail func(error) R, fSuccess func(T) R) R {
	return either.Fold(try.either, fFail, fSuccess)
}

func failWithFinally[T any](err error, f func()) Try[T] {
	return Finally(Fail[T](err), func(ignored T) { f() })
}

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

func Map[T any, R any](try Try[T], f func(T) R) Try[R] {
	return FlatMap(try, func(t T) Try[R] {
		return Success(f(t))
	})
}

func ForEach[T any](try Try[T], f func(T)) {
	try.either.ForEach(f)
}

func (try Try[T]) ForEach(f func(T)) {
	ForEach(try, f)
}

func IfFail[T any](try Try[T], f func(error)) {
	try.either.IfLeft(f)
}

func (try Try[T]) IfFail(f func(error)) {
	IfFail(try, f)
}

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

func MapLeft[T any](try Try[T], f func(err error) error) Try[T] {
	return FlatMapFail(try, func(e error) Try[T] {
		return Fail[T](f(e))
	})
}

func BiForEach[T any](try Try[T], fFail func(err error), fSuccess func(T)) {
	either.BiForEach(try.either, fFail, fSuccess)
}

func (try Try[T]) BiForEach(fFail func(err error), fSuccess func(T)) {
	BiForEach(try, fFail, fSuccess)
}

func ToOption[T any](try Try[T]) option.Option[T] {
	return try.either.ToOption()
}

func (try Try[T]) ToOption() option.Option[T] {
	return ToOption(try)
}

func ToEither[T any](try Try[T]) either.Either[error, T] {
	return try.either
}

func (try Try[T]) ToEither() either.Either[error, T] {
	return try.either
}
