package dgo

type (
	// Doer is performs some task on behalf of a caller
	Doer func()

	// Actor is performs some task with a value on behalf of a caller
	Actor func(value Value)

	// Mapper maps a value to another value
	Mapper func(value Value) interface{}

	// Producer is a function that produces a value
	Producer func() Value

	// Callable is implemented by data types that can be called such as named or
	// unnamed functions and methods.
	Callable interface {
		Value

		// Call calls the callable with the given input and returns a result
		Call(args ...interface{}) []Value
	}
)
