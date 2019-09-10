package vf

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Map creates an immutable dgo.Map from the given arguments. The arguments must either be
// one single argument which is a go map or an Array with an even number of elements, or an even number
// of arguments in which case  they will be considered a flat list of key, value [, key, value, ... ]
func Map(m ...interface{}) dgo.Map {
	return internal.Map(m)
}

// MutableMap creates an empty dgo.Map. The map can be optionally constrained
// by the given type which can be nil, the zero value of a go map, or a dgo.MapType
func MutableMap(typ interface{}) dgo.Map {
	return internal.MapWithCapacity(0, typ)
}

// MapWithCapacity creates an empty dgo.Map with the given capacity. The map can be optionally constrained
// by the given type which can be nil, the zero value of a go map, or a dgo.MapType
func MapWithCapacity(capacity int, typ interface{}) dgo.Map {
	return internal.MapWithCapacity(capacity, typ)
}

// MapFromReflected creates a Map from a reflected map. If frozen is true, the created Map will be
// immutable and the type will reflect exactly that map and nothing else. If frozen is false, the
// created Map will be mutable and its type will be derived from the reflected map.
func MapFromReflected(rm reflect.Value, frozen bool) dgo.Map {
	return internal.MapFromReflected(rm, frozen).(dgo.Map)
}
