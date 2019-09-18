package vf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Array returns a frozen dgo.Array that represents a copy of the given value. The value can be
// a slice or an Iterable
func Array(value interface{}) dgo.Array {
	return internal.Array(value)
}

// MutableArray creates a new mutable array that wraps the given slice. Unset entries in the
// slice will be replaced by Nil.
func MutableArray(typ dgo.ArrayType, slice []dgo.Value) dgo.Array {
	return internal.MutableArray(typ, slice)
}

// Values returns a frozen dgo.Array that represents the given values
func Values(values ...interface{}) dgo.Array {
	return internal.Values(values)
}

// MutableValues returns a frozen dgo.Array that represents the given values
func MutableValues(typ dgo.ArrayType, values ...interface{}) dgo.Array {
	return internal.MutableValues(typ, values)
}

// Strings returns a frozen dgo.Array that represents the given strings
func Strings(values ...string) dgo.Array {
	return internal.Strings(values)
}

// Integers returns a frozen dgo.Array that represents the given ints
func Integers(values ...int) dgo.Array {
	return internal.Integers(values)
}
