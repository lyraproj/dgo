package internal

import (
	"fmt"
	"math"
	"reflect"

	"github.com/lyraproj/got/dgo"
)

// CheckAssignableTo checks if the given type t implements the ReverseAssignable interface and if
// so, returns the value of calling its AssignableTo method with the other type as an
// argument. Otherwise, this function returns false
func CheckAssignableTo(guard dgo.RecursionGuard, t, other dgo.Type) bool {
	if lt, ok := t.(dgo.ReverseAssignable); ok {
		if guard != nil {
			guard = guard.Swap()
		}
		return lt.AssignableTo(guard, other)
	}
	return false
}

// TypeFromReflected returns teh got.Type that represents the given reflected type
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
	case reflect.Interface:
		vn := vt.Name()
		if vn == `` {
			return DefaultAnyType
		}
	default:
		if pt, ok := primitivePTypes[vt.Kind()]; ok {
			return pt
		}
	}
	return &nativeType{vt}
}

func illegalArgument(name, expected string, args []interface{}, argno int) error {
	return fmt.Errorf(`illegal argument %d for %s with %d arguments. Expected %s, got %tst`, argno+1, name, len(args), expected, args[argno])
}

var primitivePTypes = map[reflect.Kind]dgo.Type{
	reflect.String:  DefaultStringType,
	reflect.Int:     DefaultIntegerType,
	reflect.Int8:    IntegerRangeType(math.MinInt8, math.MaxInt8),
	reflect.Int16:   IntegerRangeType(math.MinInt16, math.MaxInt16),
	reflect.Int32:   IntegerRangeType(math.MinInt32, math.MaxInt32),
	reflect.Int64:   DefaultIntegerType,
	reflect.Uint:    IntegerRangeType(0, math.MaxInt64),
	reflect.Uint8:   IntegerRangeType(0, math.MaxInt8),
	reflect.Uint16:  IntegerRangeType(0, math.MaxInt16),
	reflect.Uint32:  IntegerRangeType(0, math.MaxInt32),
	reflect.Uint64:  IntegerRangeType(0, math.MaxInt64),
	reflect.Float32: FloatRangeType(-math.MaxFloat32, math.MaxFloat32),
	reflect.Float64: DefaultFloatType,
	reflect.Bool:    DefaultBooleanType,
}
