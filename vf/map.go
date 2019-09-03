package vf

import (
	"reflect"

	"github.com/lyraproj/got/dgo"
	"github.com/lyraproj/got/internal"
)

// Map creates an immutable got.Map from the given argument which must be a go map.
func Map(m interface{}) dgo.Map {
	return internal.Map(m)
}

// MutableMap creates an empty got.Map with the given capacity. The map can be optionally constrained
// by the given type which can be nil, the zero value of a go map, or a got.MapType
func MutableMap(capacity int, typ interface{}) dgo.Map {
	return internal.MutableMap(capacity, typ)
}

// MapFromReflected creates a Map from a reflected map. If frozen is true, the created Map will be
// immutable and the type will reflect exactly that map and nothing else. If frozen is false, the
// created Map will be mutable and its type will be derived from the reflected map.
func MapFromReflected(rm reflect.Value, frozen bool) dgo.Map {
	return internal.MapFromReflected(rm, frozen)
}
