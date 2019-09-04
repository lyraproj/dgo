package newtype

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Map returns a type that represents an Map value
func Map(args ...interface{}) dgo.MapType {
	return internal.MapType(args...)
}

// StructEntry returns a new MapEntryType initiated with the given parameters
func StructEntry(key string, valueType dgo.Type, required bool) dgo.MapEntryType {
	return internal.StructEntry(key, valueType, required)
}

// Struct returns a new Struct type built from the given MapEntryTypes.
func Struct(entries ...dgo.MapEntryType) dgo.StructType {
	return internal.Struct(entries)
}
