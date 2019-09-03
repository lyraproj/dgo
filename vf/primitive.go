// Package vf (Value Factory) contains all factory methods for creating values
package vf

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// True is the dgo.Value representation of true
const True = internal.True

// False is the dgo.Value representation of false
const False = internal.False

// Nil is the dgo.Value representation of nil
const Nil = internal.Nil

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

// String returns the given string as a dgo.String
func String(string string) dgo.String {
	return internal.String(string)
}

// Value converts the given value into an immutable dgo.Value
func Value(v interface{}) dgo.Value {
	return internal.Value(v)
}

// ValueFromReflected converts the given reflected value into an immutable dgo.Value
func ValueFromReflected(v reflect.Value) dgo.Value {
	return internal.ValueFromReflected(v)
}

// SameInstance returns true if the two arguments represent the same object instance.
func SameInstance(a, b dgo.Value) bool {
	return internal.SameInstance(a, b)
}
