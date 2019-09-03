// Package vf (Value Factory) contains all factory methods for creating values
package vf

import (
	"reflect"

	"github.com/lyraproj/got/dgo"
	"github.com/lyraproj/got/internal"
)

// True is the got.Value representation of true
const True = internal.True

// False is the got.Value representation of false
const False = internal.False

// Nil is the got.Value representation of nil
const Nil = internal.Nil

// Boolean returns a Boolean that represents the given bool
func Boolean(v bool) dgo.Boolean {
	if v {
		return True
	}
	return False
}

// Integer returns the given value as a got.Integer
func Integer(value int64) dgo.Integer {
	return internal.Integer(value)
}

// Float returns the given value as a got.Float
func Float(value float64) dgo.Float {
	return internal.Float(value)
}

// String returns the given string as a got.String
func String(string string) dgo.String {
	return internal.String(string)
}

// Value converts the given value into an immutable got.Value
func Value(v interface{}) dgo.Value {
	return internal.Value(v)
}

// ValueFromReflected converts the given reflected value into an immutable got.Value
func ValueFromReflected(v reflect.Value) dgo.Value {
	return internal.ValueFromReflected(v)
}

// SameInstance returns true if the two arguments represent the same object instance.
func SameInstance(a, b dgo.Value) bool {
	return internal.SameInstance(a, b)
}
