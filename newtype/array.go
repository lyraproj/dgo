package newtype

import (
	"github.com/lyraproj/got/dgo"
	"github.com/lyraproj/got/internal"
)

// Array returns a type that represents an Array value
func Array(args ...interface{}) dgo.ArrayType {
	return internal.ArrayType(args...)
}

// Tuple returns a type that represents an Array value with a specific set of typed elements
func Tuple(types ...dgo.Type) dgo.TupleType {
	return internal.TupleType(types)
}
