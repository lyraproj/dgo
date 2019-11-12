package internal_test

import (
	"math"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestParse_default(t *testing.T) {
	require.Equal(t, typ.Any, tf.Parse(`any`))
	require.Equal(t, typ.Boolean, tf.Parse(`bool`))
	require.Equal(t, typ.Integer, tf.Parse(`int`))
	require.Equal(t, typ.True, tf.Parse(`true`))
	require.Equal(t, typ.False, tf.Parse(`false`))
	require.Equal(t, typ.Float, tf.Parse(`float`))
	require.Equal(t, typ.String, tf.Parse(`string`))
	require.Equal(t, typ.Binary, tf.Parse(`binary`))
	require.Equal(t, typ.Nil, tf.Parse(`nil`))
	require.Equal(t, typ.Array, tf.Parse(`[]any`))
	require.Equal(t, typ.Map, tf.Parse(`map[any]any`))
}

func TestParse_exact(t *testing.T) {
	require.Equal(t, vf.Value(32).Type(), tf.Parse(`32`))
	require.Equal(t, vf.Value(3.14).Type(), tf.Parse(`3.14`))
	require.Equal(t, vf.Value("pelle").Type(), tf.Parse(`"pelle"`))
	require.Equal(t, tf.Tuple(vf.Value("pelle").Type(), vf.Value(3.14).Type(), typ.Boolean), tf.Parse(`{"pelle", 3.14, bool}`))

	st := tf.Parse(`{"a":1..5,"b"?:2}`)
	require.Equal(t, tf.StructMap(false,
		tf.StructMapEntry(`a`, tf.IntegerRange(1, 5, true), true),
		tf.StructMapEntry(`b`, vf.Value(2).Type(), false)), st)

	st = tf.Parse(`{"a":1..5,"b"?:2,...}`)
	require.Equal(t, tf.StructMap(true,
		tf.StructMapEntry(`a`, tf.IntegerRange(1, 5, true), true),
		tf.StructMapEntry(`b`, vf.Value(2).Type(), false)), st)

	require.Equal(t, `{"a":1..5,"b"?:2,...}`, st.String())

	st = tf.Parse(`{}`)
	require.Equal(t, tf.StructMap(false), st)
	require.Equal(t, `{}`, st.String())

	st = tf.Parse(`{...}`)
	require.Equal(t, tf.StructMap(true), st)
	require.Equal(t, `{...}`, st.String())
}

func TestParse_func(t *testing.T) {
	tt := tf.Parse(`func(string,...any) (string, bool)`)
	require.Equal(t, tf.Function(
		tf.VariadicTuple(typ.String, tf.Array(typ.Any)),
		tf.Tuple(typ.String, typ.Boolean)), tt)
}

func TestParse_tuple(t *testing.T) {
	tt := tf.Parse(`{string,...string[10]}`)
	require.Equal(t, tf.VariadicTuple(typ.String, tf.Array(tf.String(10))), tt)
}

func TestParse_ciEnum(t *testing.T) {
	st := tf.Parse(`~"foo"|~"fee"`)
	require.Equal(t, tf.CiEnum(`foo`, `fee`), st)
	require.Panic(t, func() { tf.Parse(`~32`) }, `expected a literal string, got 32`)
}

func TestParse_sized(t *testing.T) {
	require.Equal(t, tf.String(1), tf.Parse(`string[1]`))
	require.Equal(t, tf.String(1, 10), tf.Parse(`string[1,10]`))
	require.Equal(t, tf.Array(1), tf.Parse(`[1]any`))
	require.Equal(t, tf.Map(1), tf.Parse(`map[any,1]any`))
	require.Equal(t, tf.Pattern(regexp.MustCompile(`a.*`)), tf.Parse(`/a.*/`))
}

func TestParse_nestedSized(t *testing.T) {
	require.Equal(t, tf.Array(tf.String(1), 1), tf.Parse(`[1]string[1]`))
	require.Equal(t, tf.Array(tf.String(1, 10), 2, 5), tf.Parse(`[2,5]string[1,10]`))
	require.Equal(t, tf.Map(typ.String, tf.String(1), 1), tf.Parse(`map[string,1]string[1]`))
	require.Equal(t, tf.Map(typ.String, tf.String(1, 10), 2, 5), tf.Parse(`map[string,2,5]string[1,10]`))
	require.Equal(t, tf.Map(tf.Map(typ.String, typ.Integer), tf.String(1, 10), 2, 5),
		tf.Parse(`map[map[string]int,2,5]string[1,10]`))
}

func TestParse_aliasBad(t *testing.T) {
	require.Panic(t, func() { tf.Parse(`"f"=map[string]int`) }, `expected end of expression, got '='`)
	require.Panic(t, func() { tf.Parse(`m=map[string](int|<"m">)`) }, `expected an identifier, got "m"`)
	require.Panic(t, func() { tf.Parse(`m=map[string](int|<m)`) }, `expected '>', got '\)'`)
	require.Panic(t, func() { tf.Parse(`map[string](int|n)`) }, `reference to unresolved type 'n'`)
	require.Panic(t, func() { tf.Parse(`m=map[string](int|n)`) }, `reference to unresolved type 'n'`)
	require.Panic(t, func() { tf.Parse(`int=map[string](int|int)`) }, `attempt redeclare identifier 'int'`)
}

func TestParse_aliasInUnary(t *testing.T) {
	tp := tf.Parse(`type[m=map[string](string|m)]`).(dgo.UnaryType)
	require.Equal(t, `type[map[string](string|<recursive self reference to map type>)]`, tp.String())

	tp = tf.Parse(`!m=map[string](string|m)`).(dgo.UnaryType)
	require.Equal(t, `!map[string](string|<recursive self reference to map type>)`, tp.String())
}

func TestParse_aliasInAllOf(t *testing.T) {
	tp := tf.Parse(`m=[](0..5&3..8&m)`)
	require.Equal(t, `[](0..5&3..8&<recursive self reference to slice type>)`, tp.String())
}

func TestParse_aliasInOneOf(t *testing.T) {
	tp := tf.Parse(`m=[](0..5^3..8^m)`)
	require.Equal(t, `[](0..5^3..8^<recursive self reference to slice type>)`, tp.String())
}

func TestParse_range(t *testing.T) {
	require.Equal(t, tf.IntegerRange(1, 10, true), tf.Parse(`1..10`))
	require.Equal(t, tf.IntegerRange(1, 10, false), tf.Parse(`1...10`))
	require.Equal(t, tf.IntegerRange(1, math.MaxInt64, true), tf.Parse(`1..`))
	require.Equal(t, tf.IntegerRange(1, math.MaxInt64, false), tf.Parse(`1...`))
	require.Equal(t, tf.IntegerRange(math.MinInt64, 0, true), tf.Parse(`..0`))
	require.Equal(t, tf.IntegerRange(math.MinInt64, 0, false), tf.Parse(`...0`))
	require.Equal(t, tf.FloatRange(1, 10, true), tf.Parse(`1.0..10`))
	require.Equal(t, tf.FloatRange(1, 10, true), tf.Parse(`1..10.0`))
	require.Equal(t, tf.FloatRange(1, 10, true), tf.Parse(`1.0..10.0`))
	require.Equal(t, tf.FloatRange(1, 10, false), tf.Parse(`1.0...10`))
	require.Equal(t, tf.FloatRange(1, 10, false), tf.Parse(`1...10.0`))
	require.Equal(t, tf.FloatRange(1, 10, false), tf.Parse(`1.0...10.0`))
	require.Equal(t, tf.FloatRange(1, math.MaxFloat64, true), tf.Parse(`1.0..`))
	require.Equal(t, tf.FloatRange(1, math.MaxFloat64, false), tf.Parse(`1.0...`))
	require.Equal(t, tf.FloatRange(-math.MaxFloat64, 0, true), tf.Parse(`..0.0`))
	require.Equal(t, tf.FloatRange(-math.MaxFloat64, 0, false), tf.Parse(`...0.0`))

	require.Panic(t, func() { tf.Parse(`.."b"`) }, `expected an integer or a float, got "b"`)
	require.Panic(t, func() { tf.Parse(`../a*/`) }, `expected an integer or a float, got /a\*/`)
	require.Panic(t, func() { tf.Parse(`..."b"`) }, `expected an integer or a float, got "b"`)
	require.Panic(t, func() { tf.Parse(`.../a*/`) }, `expected an integer or a float, got /a\*/`)
}

func TestParse_unary(t *testing.T) {
	require.Equal(t, tf.Not(typ.String), tf.Parse(`!string`))
	require.Equal(t, typ.String.Type(), tf.Parse(`type[string]`))
	require.Equal(t, typ.Any.Type(), tf.Parse(`type`))
	require.Panic(t, func() { tf.Parse(`type[string`) }, `expected ']', got EOT`)
}

func TestParse_ternary(t *testing.T) {
	require.Equal(t, tf.AnyOf(typ.String, typ.Integer, typ.Float), tf.Parse(`string|int|float`))
	require.Equal(t, tf.OneOf(typ.String, typ.Integer, typ.Float), tf.Parse(`string^int^float`))
	require.Equal(t, tf.AllOf(typ.String, typ.Integer, typ.Float), tf.Parse(`string&int&float`))
	require.Equal(t, tf.AnyOf(typ.String, tf.OneOf(
		tf.AllOf(typ.Integer, typ.Float), typ.Boolean)), tf.Parse(`string|int&float^bool`))
	require.Equal(t, tf.AnyOf(typ.String,
		tf.AllOf(typ.Integer, tf.OneOf(typ.Float, typ.Boolean))), tf.Parse(`string|int&(float^bool)`))
	require.Equal(t, tf.Enum(`a`, `b`, `c`), tf.Parse(`"a"|"b"|"c"`))
}

func TestParse_string(t *testing.T) {
	require.Equal(t, vf.String("\r").Type(), tf.Parse(`"\r"`))
	require.Equal(t, vf.String("\n").Type(), tf.Parse(`"\n"`))
	require.Equal(t, vf.String("\t").Type(), tf.Parse(`"\t"`))
	require.Equal(t, vf.String("\"").Type(), tf.Parse(`"\""`))

	require.Equal(t, vf.String("\"").Type(), tf.Parse("`\"`"))

	require.Panic(t, func() { tf.Parse(`"\`) }, `unterminated string`)
	require.Panic(t, func() { tf.Parse(`"\"`) }, `unterminated string`)
	require.Panic(t, func() { tf.Parse(`"\y"`) }, `illegal escape`)
	require.Panic(t, func() {
		tf.Parse(`["
"`)
	}, `unterminated string`)

	require.Panic(t, func() { tf.Parse("`x") }, `unterminated string`)
}

func TestParse_multiAliases(t *testing.T) {
	tp := tf.Parse(`
{
  types: {
    ascii=1..127,
    slug=/^[a-z0-9-]+$/
  },
  x: map[slug]{token:ascii,value:string}
}`).(dgo.StructMapType)
	require.Equal(t, `map[/^[a-z0-9-]+$/]{"token":1..127,"value":string}`, tp.Get(`x`).Value().String())
}

func TestParse_errors(t *testing.T) {
	require.Panic(t, func() { tf.Parse(`[1 23]`) }, `expected one of ',' or '\]', got 23: \(column: 4\)`)
	require.Panic(t, func() { tf.Parse(`map{}`) }, `expected '\[', got '\{': \(column: 4\)`)
	require.Panic(t, func() { tf.Parse(`map[string}string`) }, `expected '\]', got '\}': \(column: 11\)`)
	require.Panic(t, func() { tf.Parse(`map[string](string|int]`) }, `expected '\)', got '\]': \(column: 23\)`)
	require.Panic(t, func() { tf.Parse(`{1 23}`) }, `expected one of ',' or '\}', got 23`)
	require.Panic(t, func() { tf.Parse(`[}int`) }, `expected a type expression, got '\}'`)
	require.Panic(t, func() { tf.Parse(`[]string[1][2]`) }, `expected end of expression, got '\['`)
	require.Panic(t, func() { tf.Parse(`[]string[1`) }, `expected one of ',' or '\]', got EOT`)
	require.Panic(t, func() { tf.Parse(`apple`) }, `reference to unresolved type 'apple'`)
	require.Panic(t, func() { tf.Parse(`-two`) }, `unexpected character 't'`)
	require.Panic(t, func() { tf.Parse(`[3e,0]`) }, `unexpected character ','`)
	require.Panic(t, func() { tf.Parse(`[3e1.0]`) }, `unexpected character '.'`)
	require.Panic(t, func() { tf.Parse(`[3e23r,4]`) }, `unexpected character 'r'`)
	require.Panic(t, func() { tf.Parse(`[3e1`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { tf.Parse(`[3e`) }, `unexpected end`)
	require.Panic(t, func() { tf.Parse(`[0x`) }, `unexpected end`)
	require.Panic(t, func() { tf.Parse(`[0x4`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { tf.Parse(`[1. 3]`) }, `unexpected character ' '`)
	require.Panic(t, func() { tf.Parse(`[1 . 3]`) }, `expected one of ',' or ']', got '.'`)
	require.Panic(t, func() { tf.Parse(`[/\`) }, `unterminated regexp`)
	require.Panic(t, func() { tf.Parse(`[/\/`) }, `unterminated regexp`)
	require.Panic(t, func() { tf.Parse(`[/\//`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { tf.Parse(`[/\t/`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { tf.Parse(`{"a":string,...`) }, `expected '}', got EOT`)
	require.Panic(t, func() { tf.Parse(`{a:32, 4}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { tf.Parse(`{"a":32, 4}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { tf.Parse(`{4, a:32}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { tf.Parse(`{4, "a":32}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { tf.Parse(`{4, a}`) }, `reference to unresolved type 'a'`)
	require.Panic(t, func() { tf.Parse(`{func, 3}`) }, `expected '\(', got ','`)
}

func TestParseFile_errors(t *testing.T) {
	require.Panic(t, func() { tf.ParseFile(nil, `foo.dgo`, `[1 2]`) }, `expected one of ',' or '\]', got 2: \(file: foo\.dgo, line: 1, column: 4\)`)
}
