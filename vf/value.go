package vf

import (
	"encoding/json"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"

	// ensure that stringer package is initialized prior to using this package
	_ "github.com/lyraproj/dgo/stringer"
)

// New creates an instance of the given type from the given argument
func New(typ dgo.Type, argument dgo.Value) dgo.Value {
	return internal.New(typ, argument)
}

// Value returns the dgo.Value representation of its argument. If the argument type
// is known, it will be more efficient to use explicit methods such as Float(), String(),
// Map(), etc.
func Value(v interface{}) dgo.Value {
	return internal.Value(v)
}

// ValueFromReflected converts the given reflected value into an immutable dgo.Value
func ValueFromReflected(v reflect.Value) dgo.Value {
	return internal.ValueFromReflected(v)
}

// ReflectTo assigns the given dgo.Value to the given reflect.Value
func ReflectTo(src dgo.Value, dest reflect.Value) {
	internal.ReflectTo(src, dest)
}

// FromJSONNumber converts the given json.Number to a dgo.Number
func FromJSONNumber(v json.Number) dgo.Number {
	return internal.FromJSONNumber(v)
}

// FromValue converts a dgo.Value into a go native value. The given `dest` must be a pointer
// to the expected native value.
func FromValue(src dgo.Value, dest interface{}) {
	internal.FromValue(src, dest)
}
