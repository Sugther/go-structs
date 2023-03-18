package tuple

import "github.com/Sugther/go-structs/equal"

/*
Tuple is a generic struct that represents a pair of values with types T1 and T2
*/
type Tuple[T1 any, T2 any] struct {
	_1 T1
	_2 T2
}

/*
Pure creates a new Tuple containing the given values of types T1 and T2.
Example: Pure(1, "hello") returns Tuple{1, "hello"}.
*/
func Pure[T1 any, T2 any](_1 T1, _2 T2) Tuple[T1, T2] {
	return Tuple[T1, T2]{
		_1: _1,
		_2: _2,
	}
}

/*
Values returns the two values stored within the Tuple.
Example: Values(Tuple{1, "hello"}) returns (1, "hello").
*/
func Values[T1 any, T2 any](tuple Tuple[T1, T2]) (T1, T2) {
	return tuple._1, tuple._2
}

func (tuple Tuple[T1, T2]) Values() (T1, T2) {
	return Values(tuple)
}

/*
Get1 returns the first value (_1) stored within the Tuple.
Example: Get1(Tuple{1, "hello"}) returns 1.
*/
func Get1[T1 any, T2 any](tuple Tuple[T1, T2]) T1 {
	return tuple._1
}

func (tuple Tuple[T1, T2]) Get1() T1 {
	return Get1(tuple)
}

/*
Get2 returns the second value (_2) stored within the Tuple.
Example: Get2(Tuple{1, "hello"}) returns "hello".
*/
func Get2[T1 any, T2 any](tuple Tuple[T1, T2]) T2 {
	return tuple._2
}

func (tuple Tuple[T1, T2]) Get2() T2 {
	return Get2(tuple)
}

/*
Map1 applies a given function f to the first value stored in the Tuple and returns a new Tuple containing the transformed value.
The function f should accept a value of type T1 and return a value of type R1.
Example: Map1(Tuple{1, "hello"}, func(x int) int { return x + 1 }) returns Tuple{2, "hello"}.
*/
func Map1[T1 any, T2 any, R1 any](tuple Tuple[T1, T2], f func(T1) R1) Tuple[R1, T2] {
	return Pure(f(tuple._1), tuple._2)
}

/*
Map2 applies a given function f to the second value stored in the Tuple and returns a new Tuple containing the transformed value.
The function f should accept a value of type T2 and return a value of type R2.
Example: Map2(Tuple{1, "hello"}, func(s string) string { return s + " world" }) returns Tuple{1, "hello world"}.
*/
func Map2[T1 any, T2 any, R2 any](tuple Tuple[T1, T2], f func(T2) R2) Tuple[T1, R2] {
	return Pure(tuple._1, f(tuple._2))
}

/*
BiMap applies two given functions f1 and f2 to the first and second values stored in the Tuple, respectively, and returns a new Tuple containing the transformed values.
The functions f1 and f2 should accept values of types T1 and T2 and return values of types R1 and R2, respectively.
Example: BiMap(Tuple{1, "hello"}, func(x int) int { return x + 1 }, func(s string) string { return s + " world" }) returns Tuple{2, "hello world"}.
*/
func BiMap[T1 any, T2 any, R1 any, R2 any](tuple Tuple[T1, T2], f1 func(T1) R1, f2 func(T2) R2) Tuple[R1, R2] {
	return Pure(f1(tuple._1), f2(tuple._2))
}

/*
FlatMap applies a given function f to the Tuple and returns a new Tuple with transformed values.
The function f should accept a Tuple of types T1 and T2 and return a Tuple of types R1 and R2.
Example: FlatMap(Tuple{1, "hello"}, func(t Tuple[int, string]) Tuple[string, int] { return Pure(t._2, t._1) }) returns Tuple{"hello", 1}.
*/
func FlatMap[T1 any, T2 any, R1 any, R2 any](tuple Tuple[T1, T2], f func(Tuple[T1, T2]) Tuple[R1, R2]) Tuple[R1, R2] {
	return f(tuple)
}

/*
ForEach applies a given function f to the Tuple.
The function f should accept a Tuple of types T1 and T2.
Example: ForEach(Tuple{1, "hello"}, func(t Tuple[int, string]) { fmt.Println(t) }) prints "Tuple{1, hello}".
*/
func ForEach[T1 any, T2 any](tuple Tuple[T1, T2], f func(Tuple[T1, T2])) {
	f(tuple)
}

func (tuple Tuple[T1, T2]) ForEach(f func(Tuple[T1, T2])) {
	ForEach(tuple, f)
}

/*
ForEach1 applies a given function f to the first value (_1) stored within the Tuple.
The function f should accept a value of type T1.
Example: ForEach1(Tuple{1, "hello"}, func(x int) { fmt.Println(x) }) prints "1".
*/
func ForEach1[T1 any, T2 any](tuple Tuple[T1, T2], f func(T1)) {
	f(tuple._1)
}

func (tuple Tuple[T1, T2]) ForEach1(f func(T1)) {
	ForEach1(tuple, f)
}

/*
ForEach2 applies a given function f to the second value (_2) stored within the Tuple.
The function f should accept a value of type T2.
Example: ForEach2(Tuple{1, "hello"}, func(s string) { fmt.Println(s) }) prints "hello".
*/
func ForEach2[T1 any, T2 any](tuple Tuple[T1, T2], f func(T2)) {
	f(tuple._2)
}

func (tuple Tuple[T1, T2]) ForEach2(f func(T2)) {
	ForEach2(tuple, f)
}

/*
Swap returns a new Tuple with the first and second values swapped.
Example: Swap(Tuple{1, "hello"}) returns Tuple{"hello", 1}.
*/
func Swap[T1 any, T2 any](tuple Tuple[T1, T2]) Tuple[T2, T1] {
	return Pure(tuple._2, tuple._1)
}

func (tuple Tuple[T1, T2]) Swap() Tuple[T2, T1] {
	return Swap(tuple)
}

/*
Equals checks if the given interface (other) is a Tuple with the same values as the current Tuple.
Returns true if the values match, false otherwise.
Example: Tuple{1, "hello"}.Equals(Tuple{1, "hello"}) returns true.
*/
func (tuple Tuple[T1, T2]) Equals(other interface{}) bool {
	if ot, ok := other.(Tuple[T1, T2]); ok {
		return equal.Equals(ot._1, tuple._1) && equal.Equals(ot._2, tuple._2)
	}
	return false
}
