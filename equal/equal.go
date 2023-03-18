package equal

import "reflect"

/*
Equal is an interface that defines a single method `Equals`, which takes an `interface{}` value and returns a `bool`.
Types that implement this interface can be compared for equality using the `Equals` method.
*/
type Equal interface {
	Equals(value interface{}) bool
}

/*
Equals is a function that compares two values for equality.
If both values implement the `Equal` interface, the function uses the `Equals` method to compare the values.
Otherwise, the function uses the `comparableEquals` function to compare the values.
*/
func Equals(value1 interface{}, value2 interface{}) bool {
	v1, okV1 := value1.(Equal)
	v2, okV2 := value2.(Equal)
	if okV1 && okV2 {
		return v1.Equals(v2)
	}
	return comparableEquals(value1, value2)
}

func comparableEquals(value1 interface{}, value2 interface{}) bool {
	return reflect.DeepEqual(value1, value2)
}

/*
IsEqual is a function that checks whether a value can be compared for equality.
If the value implements the `Equal` interface, the function returns `true`.
Otherwise, the function uses reflection to check whether the type of the value is comparable, and returns `true` if it is.
*/
func IsEqual(i interface{}) bool {
	_, ok := i.(Equal)
	return ok || reflect.TypeOf(i).Comparable()
}
