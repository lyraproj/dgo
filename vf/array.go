package vf

import (
	"github.com/lyraproj/got/dgo"
	"github.com/lyraproj/got/internal"
)

// Array returns a frozen got.Array that represents a copy of the given slice
func Array(slice []dgo.Value) dgo.Array {
	return internal.Array(slice)
}

// MutableArray returns a mutable got.Array that wraps the given slice
func MutableArray(typ dgo.ArrayType, slice []dgo.Value) dgo.Array {
	return internal.MutableArray(typ, slice)
}

// Values returns a frozen got.Array that represents the given values
func Values(values ...interface{}) dgo.Array {
	return internal.Values(values)
}

// MutableValues returns a frozen got.Array that represents the given values
func MutableValues(typ dgo.ArrayType, values ...interface{}) dgo.Array {
	return internal.MutableValues(typ, values)
}

// Strings returns a frozen got.Array that represents the given strings
func Strings(values ...string) dgo.Array {
	return internal.Strings(values)
}

// Integers returns a frozen got.Array that represents the given ints
func Integers(values ...int) dgo.Array {
	return internal.Integers(values)
}
