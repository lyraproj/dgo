package internal_test

import (
	"math"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"

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
	require.Equal(t, newtype.Struct(false,
		newtype.StructEntry(`a`, newtype.IntegerRange(1, 5), true),
		newtype.StructEntry(`b`, vf.Value(2).Type(), false)), st)

	st = newtype.Parse(`{"a":1..5,"b"?:2,...}`)
	require.Equal(t, newtype.Struct(true,
		newtype.StructEntry(`a`, newtype.IntegerRange(1, 5), true),
		newtype.StructEntry(`b`, vf.Value(2).Type(), false)), st)

	require.Equal(t, `{"a":1..5,"b"?:2,...}`, st.String())

	st = newtype.Parse(`{}`)
	require.Equal(t, newtype.Struct(false), st)
	require.Equal(t, `{}`, st.String())

	st = newtype.Parse(`{...}`)
	require.Equal(t, newtype.Struct(true), st)
	require.Equal(t, `{...}`, st.String())
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

func TestParse_aliasBad(t *testing.T) {
	require.Panic(t, func() { newtype.Parse(`"f"=map[string]int`) }, `expected end of expression, got '='`)
	require.Panic(t, func() { newtype.Parse(`m=map[string](int|<"m">)`) }, `expected an identifier, got "m"`)
	require.Panic(t, func() { newtype.Parse(`m=map[string](int|<m)`) }, `expected '>', got '\)'`)
	require.Panic(t, func() { newtype.Parse(`map[string](int|n)`) }, `reference to unresolved type 'n'`)
	require.Panic(t, func() { newtype.Parse(`m=map[string](int|n)`) }, `reference to unresolved type 'n'`)
	require.Panic(t, func() { newtype.Parse(`int=map[string](int|int)`) }, `attempt redeclare identifier 'int'`)
}

func TestParse_aliasInUnary(t *testing.T) {
	tp := newtype.Parse(`type[m=map[string](string|m)]`).(dgo.UnaryType)
	require.Equal(t, `type[map[string](string|<recursive self reference to map type>)]`, tp.String())

	tp = newtype.Parse(`!m=map[string](string|m)`).(dgo.UnaryType)
	require.Equal(t, `!map[string](string|<recursive self reference to map type>)`, tp.String())
}

func TestParse_aliasInAllOf(t *testing.T) {
	tp := newtype.Parse(`m=[](0..5&3..8&m)`)
	require.Equal(t, `[](0..5&3..8&<recursive self reference to slice type>)`, tp.String())
}

func TestParse_aliasInOneOf(t *testing.T) {
	tp := newtype.Parse(`m=[](0..5^3..8^m)`)
	require.Equal(t, `[](0..5^3..8^<recursive self reference to slice type>)`, tp.String())
}

func TestParse_range(t *testing.T) {
	require.Equal(t, newtype.IntegerRange(1, 10), newtype.Parse(`1..10`))
	require.Equal(t, newtype.IntegerRange(1, 9), newtype.Parse(`1...10`))
	require.Equal(t, newtype.IntegerRange(1, math.MaxInt64), newtype.Parse(`1..`))
	require.Equal(t, newtype.IntegerRange(math.MinInt64, 0), newtype.Parse(`..0`))
	require.Equal(t, newtype.IntegerRange(math.MinInt64, -1), newtype.Parse(`...0`))
	require.Equal(t, newtype.FloatRange(1, 10), newtype.Parse(`1.0..10`))
	require.Equal(t, newtype.FloatRange(1, 10), newtype.Parse(`1..10.0`))
	require.Equal(t, newtype.FloatRange(1, 10), newtype.Parse(`1.0..10.0`))
	require.Equal(t, newtype.FloatRange(1, math.MaxFloat64), newtype.Parse(`1.0..`))
	require.Equal(t, newtype.FloatRange(-math.MaxFloat64, 0), newtype.Parse(`..0.0`))

	require.Panic(t, func() { newtype.Parse(`1.0...10`) }, `float boundary on ... range: \(column: 4\)`)
	require.Panic(t, func() { newtype.Parse(`1...10.0`) }, `float boundary on ... range: \(column: 5\)`)
	require.Panic(t, func() { newtype.Parse(`1.0...10.0`) }, `float boundary on ... range: \(column: 4\)`)
	require.Panic(t, func() { newtype.Parse(`1.0...`) }, `float boundary on ... range: \(column: 4\)`)
	require.Panic(t, func() { newtype.Parse(`...0.0`) }, `float boundary on ... range: \(column: 4\)`)
	require.Panic(t, func() { newtype.Parse(`.."b"`) }, `expected an integer or a float, got "b"`)
	require.Panic(t, func() { newtype.Parse(`../a*/`) }, `expected an integer or a float, got /a\*/`)
	require.Panic(t, func() { newtype.Parse(`..."b"`) }, `expected an integer, got "b"`)
	require.Panic(t, func() { newtype.Parse(`.../a*/`) }, `expected an integer, got /a\*/`)
}

func TestParse_unary(t *testing.T) {
	require.Equal(t, newtype.Not(typ.String), newtype.Parse(`!string`))
	require.Equal(t, typ.String.Type(), newtype.Parse(`type[string]`))
	require.Equal(t, typ.Any.Type(), newtype.Parse(`type`))
	require.Panic(t, func() { newtype.Parse(`type[string`) }, `expected ']', got EOT`)
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

func TestParse_string(t *testing.T) {
	require.Equal(t, vf.String("\r").Type(), newtype.Parse(`"\r"`))
	require.Equal(t, vf.String("\n").Type(), newtype.Parse(`"\n"`))
	require.Equal(t, vf.String("\t").Type(), newtype.Parse(`"\t"`))
	require.Equal(t, vf.String("\"").Type(), newtype.Parse(`"\""`))

	require.Equal(t, vf.String("\"").Type(), newtype.Parse("`\"`"))

	require.Panic(t, func() { newtype.Parse(`"\`) }, `unterminated string`)
	require.Panic(t, func() { newtype.Parse(`"\"`) }, `unterminated string`)
	require.Panic(t, func() { newtype.Parse(`"\y"`) }, `illegal escape`)
	require.Panic(t, func() {
		newtype.Parse(`["
"`)
	}, `unterminated string`)

	require.Panic(t, func() { newtype.Parse("`x") }, `unterminated string`)
}

func TestParse_multiAliases(t *testing.T) {
	tp := newtype.Parse(`
{
  types: {
    ascii=1..127,
    slug=/^[a-z0-9-]+$/
  },
  x: map[slug]{token:ascii,value:string}
}`).(dgo.StructType)
	require.Equal(t, `map[/^[a-z0-9-]+$/]{"token":1..127,"value":string}`, tp.Get(`x`).Value().String())
}

func TestParse_errors(t *testing.T) {
	require.Panic(t, func() { newtype.Parse(`[1 23]`) }, `expected one of ',' or '\]', got 23: \(column: 4\)`)
	require.Panic(t, func() { newtype.Parse(`map{}`) }, `expected '\[', got '\{': \(column: 4\)`)
	require.Panic(t, func() { newtype.Parse(`map[string}string`) }, `expected '\]', got '\}': \(column: 11\)`)
	require.Panic(t, func() { newtype.Parse(`map[string](string|int]`) }, `expected '\)', got '\]': \(column: 23\)`)
	require.Panic(t, func() { newtype.Parse(`{1 23}`) }, `expected one of ',' or '\}', got 23`)
	require.Panic(t, func() { newtype.Parse(`[}int`) }, `expected a type expression, got '\}'`)
	require.Panic(t, func() { newtype.Parse(`[]string[1][2]`) }, `expected end of expression, got '\['`)
	require.Panic(t, func() { newtype.Parse(`[]string[1`) }, `expected one of ',' or '\]', got EOT`)
	require.Panic(t, func() { newtype.Parse(`apple`) }, `reference to unresolved type 'apple'`)
	require.Panic(t, func() { newtype.Parse(`-two`) }, `unexpected character 't'`)
	require.Panic(t, func() { newtype.Parse(`[3e,0]`) }, `unexpected character ','`)
	require.Panic(t, func() { newtype.Parse(`[3e1.0]`) }, `unexpected character '.'`)
	require.Panic(t, func() { newtype.Parse(`[3e23r,4]`) }, `unexpected character 'r'`)
	require.Panic(t, func() { newtype.Parse(`[3e1`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { newtype.Parse(`[3e`) }, `unexpected end`)
	require.Panic(t, func() { newtype.Parse(`[0x`) }, `unexpected end`)
	require.Panic(t, func() { newtype.Parse(`[0x4`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { newtype.Parse(`[1. 3]`) }, `unexpected character ' '`)
	require.Panic(t, func() { newtype.Parse(`[1 . 3]`) }, `expected one of ',' or ']', got '.'`)
	require.Panic(t, func() { newtype.Parse(`[/\`) }, `unterminated regexp`)
	require.Panic(t, func() { newtype.Parse(`[/\/`) }, `unterminated regexp`)
	require.Panic(t, func() { newtype.Parse(`[/\//`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { newtype.Parse(`[/\t/`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { newtype.Parse(`{"a":string,...`) }, `expected '}', got EOT`)
	require.Panic(t, func() { newtype.Parse(`{a:32, 4}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { newtype.Parse(`{"a":32, 4}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { newtype.Parse(`{4, a:32}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { newtype.Parse(`{4, "a":32}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { newtype.Parse(`{4, a}`) }, `reference to unresolved type 'a'`)
}

func TestParseFile_errors(t *testing.T) {
	require.Panic(t, func() { newtype.ParseFile(`foo.dgo`, `[1 2]`) }, `expected one of ',' or '\]', got 2: \(file: foo\.dgo, line: 1, column: 4\)`)
}
