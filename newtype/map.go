package newtype

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Map returns a type that represents an Map value
func Map(args ...interface{}) dgo.MapType {
	return internal.MapType(args...)
}

// StructEntry returns a new StructEntry initiated with the given parameters
func StructEntry(key interface{}, value interface{}, required bool) dgo.StructEntry {
	return internal.StructEntry(key, value, required)
}

// Struct returns a new Struct type built from the given MapEntryTypes. If
// additional is true, the struct will allow additional unconstrained entries
func Struct(additional bool, entries ...dgo.StructEntry) dgo.StructType {
	return internal.Struct(additional, entries)
}

// StructFromMap returns a new type built from a map[string](dgo|type|{type:dgo|type,required?:bool,...})
func StructFromMap(additional bool, entries dgo.Map) dgo.StructType {
	return internal.StructFromMap(additional, entries)
}
