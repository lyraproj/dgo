package internal_test

import (
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestType_String_floatExact(t *testing.T) {
	require.Equal(t, `1.0`, vf.Value(1.0).Type().String())
}

func TestType_String_stringExact(t *testing.T) {
	require.Equal(t, `"with \"quotes\" in it"`, vf.Value("with \"quotes\" in it").Type().String())
}

func TestType_String_arrayExact(t *testing.T) {
	require.Equal(t, `{1.0,2.3}`, vf.Values(1.0, 2.3).Type().String())
}

func TestType_String_arrayUnboundedElementInt(t *testing.T) {
	require.Equal(t, `[]int`, newtype.Array(typ.Integer).String())
}

func TestType_String_arrayBoundedElementInt(t *testing.T) {
	require.Equal(t, `[0,4]int`, newtype.Array(typ.Integer, 0, 4).String())
}

func TestType_String_arrayBoundedElementRange(t *testing.T) {
	require.Equal(t, `[0,4]3..8`, newtype.Array(newtype.IntegerRange(3, 8), 0, 4).String())
}

func TestType_String_arrayBoundedElementEnum(t *testing.T) {
	require.Equal(t, `[0,4]("a"|"b")`, newtype.Array(newtype.Enum(`a`, `b`), 0, 4).String())
}

func TestType_String_mapUnboundedEntryStringInt(t *testing.T) {
	require.Equal(t, `map[string]int`, newtype.Map(typ.String, typ.Integer).String())
}

func TestType_String_arrayBoundedEntryStringInt(t *testing.T) {
	require.Equal(t, `map[string,0,4]int`, newtype.Map(typ.String, typ.Integer, 0, 4).String())
}

func TestType_String_arrayBoundedEntryStringRange(t *testing.T) {
	require.Equal(t, `map[string,0,4]3..8`, newtype.Map(typ.String, newtype.IntegerRange(3, 8), 0, 4).String())
}

func TestType_String_arrayBoundedEntryStringEnum(t *testing.T) {
	require.Equal(t, `map[string,0,4]("a"|"b")`, newtype.Map(typ.String, newtype.Enum(`a`, `b`), 0, 4).String())
}

func TestType_String_arrayBoundedEntryIntRangeEnum(t *testing.T) {
	require.Equal(t, `map[3..8,0,4]("a"|"b")`, newtype.Map(newtype.IntegerRange(3, 8), newtype.Enum(`a`, `b`), 0, 4).String())
}

func TestType_String_priorities(t *testing.T) {
	var tp dgo.Type = newtype.Array(newtype.String(1), 2, 2)
	require.Equal(t, `[2,2]string[1]`, tp.String())

	tp = newtype.Array(vf.Value(regexp.MustCompile(`a`)).Type(), 2, 2)
	require.Equal(t, `[2,2]regexp["a"]`, tp.String())

	tp = newtype.Array(newtype.Not(newtype.String(1)), 2, 2)
	require.Equal(t, `[2,2]!string[1]`, tp.String())

	tp = newtype.Not(newtype.String(1))
	require.Equal(t, `!string[1]`, tp.String())

	tp = newtype.Array(newtype.Map(typ.Integer, typ.String, 1), 2, 2)
	require.Equal(t, `[2,2]map[int,1]string`, tp.String())

	tp = newtype.Array(newtype.Map(typ.Integer, newtype.String(1), 2), 3, 4)
	require.Equal(t, `[3,4]map[int,2]string[1]`, tp.String())

	tp = newtype.Array(newtype.AllOf(newtype.Enum(`a`, `b`), newtype.Enum(`b`, `c`)), 2, 2)
	require.Equal(t, `[2,2](("a"|"b")&("b"|"c"))`, tp.String())

	tp = newtype.Array(newtype.OneOf(typ.Integer, typ.String), 2, 2)
	require.Equal(t, `[2,2](int^string)`, tp.String())

	tp = newtype.Array(vf.Values(`a`, `b`).Type().(dgo.ArrayType).ElementType(), 2, 2)
	require.Equal(t, `[2,2]("a"&"b")`, tp.String())
}
