package internal

import (
	"fmt"
	"math"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

// CheckAssignableTo checks if the given type t implements the ReverseAssignable interface and if
// so, returns the value of calling its AssignableTo method with the other type as an
// argument. Otherwise, this goFunc returns false
func CheckAssignableTo(guard dgo.RecursionGuard, t, other dgo.Type) bool {
	if lt, ok := t.(dgo.ReverseAssignable); ok {
		if guard != nil {
			guard = guard.Swap()
		}
		return lt.AssignableTo(guard, other)
	}
	return false
}

// TypeFromReflected returns the dgo.Type that represents the given reflected type
func TypeFromReflected(vt reflect.Type) dgo.Type {
	if pt, ok := wellKnownTypes[vt]; ok {
		return pt
	}

	kind := vt.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		return ArrayType(TypeFromReflected(vt.Elem()), 0, math.MaxInt64)
	case reflect.Map:
		return MapType(TypeFromReflected(vt.Key()), TypeFromReflected(vt.Elem()), 0, math.MaxInt64)
	case reflect.Ptr:
		return OneOfType([]dgo.Type{TypeFromReflected(vt.Elem()), DefaultNilType})
	case reflect.Func:
		return exactFunctionType{vt}
	case reflect.Interface:
		switch vt.Name() {
		case ``:
			return DefaultAnyType
		case `error`:
			return DefaultErrorType
		}
	default:
		if pt, ok := primitivePTypes[vt.Kind()]; ok {
			return pt
		}
	}
	return &nativeType{vt}
}

// Generic returns the generic form of the given type. All non exact types are considered generic
// and will be returned directly. Exact types will loose information about what instance they represent
// and also range and size information. Nested types will return a generic version of the contained
// types as well.
func Generic(t dgo.Type) dgo.Type {
	if et, ok := t.(dgo.ExactType); ok {
		return et.Generic()
	}
	return t
}

func typeAsType(v dgo.Value) dgo.Type {
	return v.(dgo.Type)
}

func valueAsType(v dgo.Value) dgo.Type {
	return v.Type()
}

func illegalArgument(name, expected string, args []interface{}, argno int) error {
	return fmt.Errorf(`illegal argument %d for %s with %d arguments. Expected %s, got %tst`, argno+1, name, len(args), expected, args[argno])
}

var primitivePTypes = map[reflect.Kind]dgo.Type{
	reflect.String:  DefaultStringType,
	reflect.Int:     DefaultIntegerType,
	reflect.Int8:    IntegerRangeType(math.MinInt8, math.MaxInt8, true),
	reflect.Int16:   IntegerRangeType(math.MinInt16, math.MaxInt16, true),
	reflect.Int32:   IntegerRangeType(math.MinInt32, math.MaxInt32, true),
	reflect.Int64:   DefaultIntegerType,
	reflect.Uint:    IntegerRangeType(0, math.MaxInt64, true),
	reflect.Uint8:   IntegerRangeType(0, math.MaxInt8, true),
	reflect.Uint16:  IntegerRangeType(0, math.MaxInt16, true),
	reflect.Uint32:  IntegerRangeType(0, math.MaxInt32, true),
	reflect.Uint64:  IntegerRangeType(0, math.MaxInt64, true),
	reflect.Float32: FloatRangeType(-math.MaxFloat32, math.MaxFloat32, true),
	reflect.Float64: DefaultFloatType,
	reflect.Bool:    DefaultBooleanType,
}
