package internal_test

import (
	"math"
	"regexp"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestParse_default(t *testing.T) {
	require.Equal(t, typ.Any, newtype.Parse(`any`))
	require.Equal(t, typ.Boolean, newtype.Parse(`bool`))
	require.Equal(t, typ.Integer, newtype.Parse(`int`))
	require.Equal(t, typ.True, newtype.Parse(`true`))
	require.Equal(t, typ.False, newtype.Parse(`false`))
	require.Equal(t, typ.Float, newtype.Parse(`float`))
	require.Equal(t, typ.String, newtype.Parse(`string`))
	require.Equal(t, typ.Binary, newtype.Parse(`binary`))
	require.Equal(t, typ.Nil, newtype.Parse(`nil`))
	require.Equal(t, typ.Array, newtype.Parse(`[]any`))
	require.Equal(t, typ.Map, newtype.Parse(`map[any]any`))
}

func TestParse_exact(t *testing.T) {
	require.Equal(t, vf.Value(32).Type(), newtype.Parse(`32`))
	require.Equal(t, vf.Value(3.14).Type(), newtype.Parse(`3.14`))
	require.Equal(t, vf.Value("pelle").Type(), newtype.Parse(`"pelle"`))
	require.Equal(t, newtype.Tuple(vf.Value("pelle").Type(), vf.Value(3.14).Type(), typ.Boolean), newtype.Parse(`{"pelle", 3.14, bool}`))

	st := newtype.Parse(`{"a":1..5,"b"?:2}`)
	require.Equal(t, newtype.Struct(
		newtype.StructEntry(`a`, newtype.IntegerRange(1, 5), true),
		newtype.StructEntry(`b`, vf.Value(2).Type(), false)), st)

	require.Equal(t, `{"a":1..5,"b"?:2}`, st.String())
}

func TestParse_sized(t *testing.T) {
	require.Equal(t, newtype.String(1), newtype.Parse(`string[1]`))
	require.Equal(t, newtype.String(1, 10), newtype.Parse(`string[1,10]`))
	require.Equal(t, newtype.Array(1), newtype.Parse(`[1]any`))
	require.Equal(t, newtype.Map(1), newtype.Parse(`map[any,1]any`))
	require.Equal(t, newtype.Pattern(regexp.MustCompile(`a.*`)), newtype.Parse(`/a.*/`))
}

func TestParse_nestedSized(t *testing.T) {
	require.Equal(t, newtype.Array(newtype.String(1), 1), newtype.Parse(`[1]string[1]`))
	require.Equal(t, newtype.Array(newtype.String(1, 10), 2, 5), newtype.Parse(`[2,5]string[1,10]`))
	require.Equal(t, newtype.Map(typ.String, newtype.String(1), 1), newtype.Parse(`map[string,1]string[1]`))
	require.Equal(t, newtype.Map(typ.String, newtype.String(1, 10), 2, 5), newtype.Parse(`map[string,2,5]string[1,10]`))
	require.Equal(t, newtype.Map(newtype.Map(typ.String, typ.Integer), newtype.String(1, 10), 2, 5),
		newtype.Parse(`map[map[string]int,2,5]string[1,10]`))
}

func TestParse_range(t *testing.T) {
	require.Equal(t, newtype.IntegerRange(1, 10), newtype.Parse(`1..10`))
	require.Equal(t, newtype.IntegerRange(1, math.MaxInt64), newtype.Parse(`1..`))
	require.Equal(t, newtype.IntegerRange(math.MinInt64, 0), newtype.Parse(`..0`))
	require.Equal(t, newtype.FloatRange(1, 10), newtype.Parse(`1.0..10`))
	require.Equal(t, newtype.FloatRange(1, 10), newtype.Parse(`1..10.0`))
	require.Equal(t, newtype.FloatRange(1, 10), newtype.Parse(`1.0..10.0`))
	require.Equal(t, newtype.FloatRange(1, math.MaxFloat64), newtype.Parse(`1.0..`))
	require.Equal(t, newtype.FloatRange(-math.MaxFloat64, 0), newtype.Parse(`..0.0`))
}

func TestParse_unary(t *testing.T) {
	require.Equal(t, newtype.Not(typ.String), newtype.Parse(`!string`))
}

func TestParse_ternary(t *testing.T) {
	require.Equal(t, newtype.AnyOf(typ.String, typ.Integer, typ.Float), newtype.Parse(`string|int|float`))
	require.Equal(t, newtype.OneOf(typ.String, typ.Integer, typ.Float), newtype.Parse(`string^int^float`))
	require.Equal(t, newtype.AllOf(typ.String, typ.Integer, typ.Float), newtype.Parse(`string&int&float`))
	require.Equal(t, newtype.AnyOf(typ.String, newtype.OneOf(
		newtype.AllOf(typ.Integer, typ.Float), typ.Boolean)), newtype.Parse(`string|int&float^bool`))
	require.Equal(t, newtype.AnyOf(typ.String,
		newtype.AllOf(typ.Integer, newtype.OneOf(typ.Float, typ.Boolean))), newtype.Parse(`string|int&(float^bool)`))
	require.Equal(t, newtype.Enum(`a`, `b`, `c`), newtype.Parse(`"a"|"b"|"c"`))
}

func TestParse_errors(t *testing.T) {
	require.Panic(t, func() { newtype.Parse(`[1 23]`) }, `expected one of ',' or '\]', got 23: \(column: 4\)`)
	require.Panic(t, func() { newtype.Parse(`map{}`) }, `expected '\[', got '\{': \(column: 4\)`)
	require.Panic(t, func() { newtype.Parse(`map[string}string`) }, `expected '\]', got '\}': \(column: 11\)`)
	require.Panic(t, func() { newtype.Parse(`map[string](string|int]`) }, `expected '\)', got '\]': \(column: 23\)`)
	require.Panic(t, func() { newtype.Parse(`.."b"`) }, `expected an literal integer or a float, got "b"`)
	require.Panic(t, func() { newtype.Parse(`{1 23}`) }, `expected one of ',' or '\}', got 23`)
	require.Panic(t, func() { newtype.Parse(`[}int`) }, `expected a type expression, got '\}'`)
	require.Panic(t, func() { newtype.Parse(`[]string[1][2]`) }, `expected end of expression, got '\['`)
	require.Panic(t, func() { newtype.Parse(`[]string[1`) }, `expected one of ',' or '\]', got EOF`)
	require.Panic(t, func() { newtype.Parse(`apple`) }, `unknown identifier 'apple'`)
}

func TestParseFile_errors(t *testing.T) {
	require.Panic(t, func() { newtype.ParseFile(`foo.dgo`, `[1 2]`) }, `expected one of ',' or '\]', got 2: \(file: foo\.dgo, line: 1, column: 4\)`)
}
