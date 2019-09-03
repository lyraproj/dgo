package newtype

import (
	"github.com/lyraproj/got/dgo"
	"github.com/lyraproj/got/internal"
)

// AllOf returns a type that represents all values that matches all of the included types
func AllOf(types ...dgo.Type) dgo.Type {
	return internal.AllOfType(types)
}

// AnyOf returns a type that represents all values that matches at least one of the included types
func AnyOf(types ...dgo.Type) dgo.Type {
	return internal.AnyOfType(types)
}

// OneOf returns a type that represents all values that matches exactly one of the included types
func OneOf(types ...dgo.Type) dgo.Type {
	return internal.OneOfType(types)
}

// Not returns a type that represents all values that are not represented by the given type
func Not(t dgo.Type) dgo.Type {
	return internal.NotType(t)
}
