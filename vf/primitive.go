// Package vf (Value Factory) contains all factory methods for creating values
package vf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// True is the dgo.Value representation of true
const True = internal.True

// False is the dgo.Value representation of false
const False = internal.False

// Nil is the dgo.Value representation of nil
const Nil = internal.Nil

// Binary creates a new Binary based on the given slice. If frozen is true, the
// binary will be immutable and contain a copy of the slice, otherwise the slice
// is simply wrapped and modifications to its elements will also modify the binary.
func Binary(bs []byte, frozen bool) dgo.Binary {
	return internal.Binary(bs, frozen)
}

// BinaryFromString creates a new Binary from the base64 encoded string
func BinaryFromString(base64 string) dgo.Binary {
	return internal.BinaryFromString(base64)
}

// Boolean returns a Boolean that represents the given bool
func Boolean(v bool) dgo.Boolean {
	if v {
		return True
	}
	return False
}

// Integer returns the given value as a dgo.Integer
func Integer(value int64) dgo.Integer {
	return internal.Integer(value)
}

// Float returns the given value as a dgo.Float
func Float(value float64) dgo.Float {
	return internal.Float(value)
}

// Sensitive creates a new Sensitive that wraps the given value
func Sensitive(v dgo.Value) dgo.Sensitive {
	return internal.Sensitive(v)
}

// String returns the given string as a dgo.String
func String(string string) dgo.String {
	return internal.String(string)
}
