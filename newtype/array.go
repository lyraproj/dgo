package newtype

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Array returns a type that represents an Array value
func Array(args ...interface{}) dgo.ArrayType {
	return internal.ArrayType(args...)
}

// Tuple returns a type that represents an Array value with a specific set of typed elements
func Tuple(types ...dgo.Type) dgo.TupleType {
	return internal.TupleType(types)
}

// VariadicTuple returns a type that represents an Array value with a variadic number of elements. Each
// given type determines the type of a corresponding element in an array except for the last one which
// must be an ArrayType that determines the remaining elements.
func VariadicTuple(types ...dgo.Type) dgo.TupleType {
	return internal.VariadicTupleType(types)
}
