package tf

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Named returns the type with the given name from the global type registry. The function returns
// nil if no type has been registered under the given name.
func Named(name string) dgo.NamedType {
	return internal.NamedType(name)
}

// ExactNamed returns the exact NamedType that represents the given value.
//
// This is the function that the value.Type() method of a named type instance uses to
// obtain the actual type.
func ExactNamed(typ dgo.NamedType, value dgo.Value) dgo.NamedType {
	return internal.ExactNamedType(typ, value)
}

// NewNamed registers a new named type under the given name with the global type registry. The method panics
// if a type has already been registered with the same name.
func NewNamed(name string, ctor dgo.Constructor, extractor dgo.InitArgExtractor, implType, ifdType reflect.Type) dgo.NamedType {
	return internal.NewNamedType(name, ctor, extractor, implType, ifdType)
}

// NamedFromReflected returns the named type for the reflected implementation type from the global type
// registry. The function returns nil if no such type has been registered.
func NamedFromReflected(rt reflect.Type) dgo.NamedType {
	return internal.NamedTypeFromReflected(rt)
}

// RemoveNamed removes a named type from the global type registry. It is primarily intended for
// testing purposes.
func RemoveNamed(name string) {
	internal.RemoveNamedType(name)
}
