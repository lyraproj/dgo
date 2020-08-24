package parser_test

import (
	"math"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/stringer"

	"github.com/lyraproj/dgo/internal"

	"github.com/lyraproj/dgo/dgo"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestParse_default(t *testing.T) {
	require.Equal(t, typ.Any, tf.ParseType(`any`))
	require.Equal(t, typ.Boolean, tf.ParseType(`bool`))
	require.Equal(t, typ.Integer, tf.ParseType(`int`))
	require.Equal(t, typ.True, tf.ParseType(`true`))
	require.Equal(t, typ.False, tf.ParseType(`false`))
	require.Equal(t, typ.Float, tf.ParseType(`float`))
	require.Equal(t, typ.String, tf.ParseType(`string`))
	require.Equal(t, typ.Binary, tf.ParseType(`binary`))
	require.Equal(t, typ.Nil, tf.ParseType(`nil`))
	require.Equal(t, typ.Array, tf.ParseType(`[]any`))
	require.Equal(t, typ.Map, tf.ParseType(`map[any]any`))
}

func TestParse_exact(t *testing.T) {
	require.Equal(t, vf.Value(32).Type(), tf.ParseType(`32`))
	require.Equal(t, vf.Value(3.14).Type(), tf.ParseType(`3.14`))
	require.Equal(t, vf.Value("pelle").Type(), tf.ParseType(`"pelle"`))
	require.Equal(t, tf.Tuple(vf.Value("pelle").Type(), vf.Value(3.14).Type(), typ.Boolean), tf.ParseType(`{"pelle", 3.14, bool}`))

	st := tf.ParseType(`{"a":1..5,"b"?:2}`)
	require.Equal(t, tf.StructMap(false,
		tf.StructMapEntry(`a`, tf.Integer(1, 5, true), true),
		tf.StructMapEntry(`b`, vf.Value(2).Type(), false)), st)

	st = tf.ParseType(`{"a":1..5,"b"?:2,...}`)
	require.Equal(t, tf.StructMap(true,
		tf.StructMapEntry(`a`, tf.Integer(1, 5, true), true),
		tf.StructMapEntry(`b`, vf.Value(2).Type(), false)), st)

	require.Equal(t, `{"a":1..5,"b"?:2,...}`, st.String())

	st = tf.ParseType(`{}`)
	require.Equal(t, tf.StructMap(false), st)
	require.Equal(t, `{}`, st.String())

	st = tf.ParseType(`{...}`)
	require.Equal(t, tf.StructMap(true), st)
	require.Equal(t, `{...}`, st.String())
}

func TestParse_func(t *testing.T) {
	tt := tf.ParseType(`func(string,...any) (string, bool)`)
	require.Equal(t, tf.Function(
		tf.VariadicTuple(typ.String, typ.Any),
		tf.Tuple(typ.String, typ.Boolean)), tt)
}

func TestParse_tuple(t *testing.T) {
	require.Equal(t, tf.VariadicTuple(typ.String, tf.String(10)), tf.ParseType(`{string,...string[10]}`))
	require.Equal(t, tf.Function(typ.EmptyTuple, typ.EmptyTuple), tf.ParseType(`func()`))
}

func TestParse_ciEnum(t *testing.T) {
	st := tf.ParseType(`~"foo"|~"fee"`)
	require.Equal(t, tf.CiEnum(`foo`, `fee`), st)
	require.Panic(t, func() { tf.ParseType(`~32`) }, `expected a literal string, got 32`)
}

func TestParse_sized(t *testing.T) {
	require.Equal(t, tf.String(1), tf.ParseType(`string[1]`))
	require.Equal(t, tf.String(1, 10), tf.ParseType(`string[1,10]`))
	require.Equal(t, tf.Array(1), tf.ParseType(`[1]any`))
	require.Equal(t, tf.Map(1), tf.ParseType(`map[any,1]any`))
	require.Equal(t, tf.Pattern(regexp.MustCompile(`a.*`)), tf.ParseType(`/a.*/`))
}

func TestParse_nestedSized(t *testing.T) {
	require.Equal(t, tf.Array(tf.String(1), 1), tf.ParseType(`[1]string[1]`))
	require.Equal(t, tf.Array(tf.String(1, 10), 2, 5), tf.ParseType(`[2,5]string[1,10]`))
	require.Equal(t, tf.Map(typ.String, tf.String(1), 1), tf.ParseType(`map[string,1]string[1]`))
	require.Equal(t, tf.Map(typ.String, tf.String(1, 10), 2, 5), tf.ParseType(`map[string,2,5]string[1,10]`))
	require.Equal(t, tf.Map(tf.Map(typ.String, typ.Integer), tf.String(1, 10), 2, 5),
		tf.ParseType(`map[map[string]int,2,5]string[1,10]`))
}

func TestParse_aliasBad(t *testing.T) {
	internal.ResetDefaultAliases()
	require.Panic(t, func() { tf.ParseType(`"f"=map[string]int`) }, `expected end of expression, got '='`)
	internal.ResetDefaultAliases()
	require.Panic(t, func() { tf.ParseType(`m=map[string](int|<"m">)`) }, `expected an identifier, got "m"`)
	internal.ResetDefaultAliases()
	require.Panic(t, func() { tf.ParseType(`m=map[string](int|<m)`) }, `expected '>', got '\)'`)
	internal.ResetDefaultAliases()
	require.Panic(t, func() { tf.ParseType(`map[string](int|n)`) }, `reference to unresolved type 'n'`)
	internal.ResetDefaultAliases()
	require.Panic(t, func() { tf.ParseType(`m=map[string](int|n)`) }, `reference to unresolved type 'n'`)
	internal.ResetDefaultAliases()
	require.Panic(t, func() { tf.ParseType(`int=map[string](int|int)`) }, `attempt to redeclare identifier 'int'`)
}

func TestParse_aliasInUnary(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`type[m=map[string](string|m)]`).(dgo.UnaryType)
	require.Equal(t, `type[m]`, tp.String())

	internal.ResetDefaultAliases()
	tp = tf.ParseType(`!m=map[string](string|m)`).(dgo.UnaryType)
	require.Equal(t, `!m`, tp.String())
}

func TestParse_aliasInAllOf(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`m=[](0..5&3..8&m)`)
	require.Equal(t, `m`, tp.String())
}

func TestParse_aliasInAllOf_separateAm(t *testing.T) {
	internal.ResetDefaultAliases()
	tf.BuiltInAliases().Collect(func(a dgo.AliasAdder) {
		tp := tf.ParseFile(a, `internal`, `m=[](0..5&3..8&m)`)
		require.Equal(t, `[](0..5&3..8&<recursive self reference to slice type>)`, tp.String())
	})
}

func TestParse_aliasInOneOf(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`m=[](0..5^3..8^m)`)
	require.Equal(t, `m`, tp.String())
}

func TestParse_data(t *testing.T) {
	tf.BuiltInAliases().Collect(func(aa dgo.AliasAdder) {
		tp := tf.ParseFile(aa, `internal`, `data`)
		require.Equal(t, `data`, tp.String())
	})
}

func TestParse_range(t *testing.T) {
	require.Equal(t, tf.Integer(1, 10, true), tf.ParseType(`1..10`))
	require.Equal(t, tf.Integer(1, 10, false), tf.ParseType(`1...10`))
	require.Equal(t, tf.Integer(1, math.MaxInt64, true), tf.ParseType(`1..`))
	require.Equal(t, tf.Integer(1, math.MaxInt64, false), tf.ParseType(`1...`))
	require.Equal(t, tf.Integer(math.MinInt64, 0, true), tf.ParseType(`..0`))
	require.Equal(t, tf.Integer(math.MinInt64, 0, false), tf.ParseType(`...0`))
	require.Equal(t, tf.Float(1, 10, true), tf.ParseType(`1.0..10`))
	require.Equal(t, tf.Float(1, 10, true), tf.ParseType(`1..10.0`))
	require.Equal(t, tf.Float(1, 10, true), tf.ParseType(`1.0..10.0`))
	require.Equal(t, tf.Float(1, 10, false), tf.ParseType(`1.0...10`))
	require.Equal(t, tf.Float(1, 10, false), tf.ParseType(`1...10.0`))
	require.Equal(t, tf.Float(1, 10, false), tf.ParseType(`1.0...10.0`))
	require.Equal(t, tf.Float(1, math.MaxFloat64, true), tf.ParseType(`1.0..`))
	require.Equal(t, tf.Float(1, math.MaxFloat64, false), tf.ParseType(`1.0...`))
	require.Equal(t, tf.Float(-math.MaxFloat64, 0, true), tf.ParseType(`..0.0`))
	require.Equal(t, tf.Float(-math.MaxFloat64, 0, false), tf.ParseType(`...0.0`))

	require.Panic(t, func() { tf.ParseType(`.."b"`) }, `expected an integer or a float, got "b"`)
	require.Panic(t, func() { tf.ParseType(`../a*/`) }, `expected an integer or a float, got /a\*/`)
	require.Panic(t, func() { tf.ParseType(`..."b"`) }, `expected an integer or a float, got "b"`)
	require.Panic(t, func() { tf.ParseType(`.../a*/`) }, `expected an integer or a float, got /a\*/`)
}

func TestParse_unary(t *testing.T) {
	require.Equal(t, tf.Not(typ.String), tf.ParseType(`!string`))
	require.Equal(t, typ.String.Type(), tf.ParseType(`type[string]`))
	require.Equal(t, typ.Any.Type(), tf.ParseType(`type`))
	require.Panic(t, func() { tf.ParseType(`type[string`) }, `expected ']', got EOT`)
}

func TestParse_ternary(t *testing.T) {
	require.Equal(t, typ.Number, tf.ParseType(`int|float`))
	require.Equal(t, tf.AnyOf(typ.String, typ.Integer, typ.Float), tf.ParseType(`string|int|float`))
	require.Equal(t, tf.OneOf(typ.String, typ.Integer, typ.Float), tf.ParseType(`string^int^float`))
	require.Equal(t, tf.AllOf(typ.String, typ.Integer, typ.Float), tf.ParseType(`string&int&float`))
	require.Equal(t, tf.AnyOf(typ.String, tf.OneOf(
		tf.AllOf(typ.Integer, typ.Float), typ.Boolean)), tf.ParseType(`string|int&float^bool`))
	require.Equal(t, tf.AnyOf(typ.String,
		tf.AllOf(typ.Integer, tf.OneOf(typ.Float, typ.Boolean))), tf.ParseType(`string|int&(float^bool)`))
	require.Equal(t, tf.Enum(`a`, `b`, `c`), tf.ParseType(`"a"|"b"|"c"`))
}

func TestParse_string(t *testing.T) {
	require.Equal(t, vf.String("\r").Type(), tf.ParseType(`"\r"`))
	require.Equal(t, vf.String("\n").Type(), tf.ParseType(`"\n"`))
	require.Equal(t, vf.String("\t").Type(), tf.ParseType(`"\t"`))
	require.Equal(t, vf.String("\"").Type(), tf.ParseType(`"\""`))

	require.Equal(t, vf.String("\"").Type(), tf.ParseType("`\"`"))

	require.Panic(t, func() { tf.ParseType(`"\`) }, `unterminated string`)
	require.Panic(t, func() { tf.ParseType(`"\"`) }, `unterminated string`)
	require.Panic(t, func() { tf.ParseType(`"\y"`) }, `illegal escape`)
	require.Panic(t, func() {
		tf.ParseType(`["
"`)
	}, `unterminated string`)

	require.Panic(t, func() { tf.ParseType("`x") }, `unterminated string`)
}

func TestParse_multiAliases(t *testing.T) {
	var tp dgo.StructMapType
	am := tf.BuiltInAliases().Collect(func(aa dgo.AliasAdder) {
		tp = tf.ParseFile(aa, `internal`, `
{
  types: {
    ascii=1..127,
    slug=/^[a-z0-9-]+$/
  },
  x: map[slug]{Token:ascii,value:string}
}`).(dgo.StructMapType)
	})

	require.Equal(t, `map[slug]{"Token":ascii,"value":string}`, stringer.TypeStringWithAliasMap(tp.GetEntryType(`x`).Value().(dgo.Type), am))
}

func TestParse_mapKeyAliases(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`{tp = "key", { key: 2 }}`)
	require.Equal(t, `{tp,{tp:2}}`, stringer.TypeString(tp))
}

func TestParse_errors(t *testing.T) {
	require.Panic(t, func() { tf.ParseType(`[1 23]`) }, `expected one of ',' or '\]', got 23: \(column: 4\)`)
	require.Panic(t, func() { tf.ParseType(`map{}`) }, `expected '\[', got '\{': \(column: 4\)`)
	require.Panic(t, func() { tf.ParseType(`map[string}string`) }, `expected '\]', got '\}': \(column: 11\)`)
	require.Panic(t, func() { tf.ParseType(`map[string](string|int]`) }, `expected '\)', got '\]': \(column: 23\)`)
	require.Panic(t, func() { tf.ParseType(`{1 23}`) }, `expected one of ',' or '\}', got 23`)
	require.Panic(t, func() { tf.ParseType(`[}int`) }, `expected a type expression, got '\}'`)
	require.Panic(t, func() { tf.ParseType(`[]string[1][2]`) }, `expected end of expression, got '\['`)
	require.Panic(t, func() { tf.ParseType(`[]string[1`) }, `expected one of ',' or '\]', got EOT`)
	require.Panic(t, func() { tf.ParseType(`apple`) }, `reference to unresolved type 'apple'`)
	require.Panic(t, func() { tf.ParseType(`-two`) }, `unexpected character 't'`)
	require.Panic(t, func() { tf.ParseType(`[3e,0]`) }, `unexpected character ','`)
	require.Panic(t, func() { tf.ParseType(`[3e1.0]`) }, `unexpected character '.'`)
	require.Panic(t, func() { tf.ParseType(`[3e23r,4]`) }, `unexpected character 'r'`)
	require.Panic(t, func() { tf.ParseType(`[3e1`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { tf.ParseType(`[3e`) }, `unexpected end`)
	require.Panic(t, func() { tf.ParseType(`[0x`) }, `unexpected end`)
	require.Panic(t, func() { tf.ParseType(`[0x4`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { tf.ParseType(`[1. 3]`) }, `unexpected character ' '`)
	require.Panic(t, func() { tf.ParseType(`[1 . 3]`) }, `expected one of ',' or ']', got '.'`)
	require.Panic(t, func() { tf.ParseType(`[/\`) }, `unterminated regexp`)
	require.Panic(t, func() { tf.ParseType(`[/\/`) }, `unterminated regexp`)
	require.Panic(t, func() { tf.ParseType(`[/\//`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { tf.ParseType(`[/\t/`) }, `expected one of ',' or ']', got EOT`)
	require.Panic(t, func() { tf.ParseType(`{"a":string,...`) }, `expected '}', got EOT`)
	require.Panic(t, func() { tf.ParseType(`{a:32, 4}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { tf.ParseType(`{"a":32, 4}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { tf.ParseType(`{4, a:32}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { tf.ParseType(`{4, "a":32}`) }, `mix of elements and map entries`)
	require.Panic(t, func() { tf.ParseType(`{4, a}`) }, `reference to unresolved type 'a'`)
	require.Panic(t, func() { tf.ParseType(`{func, 3}`) }, `expected '\(', got ','`)
}

func TestParseFile_errors(t *testing.T) {
	require.Panic(t,
		func() { tf.ParseFile(nil, `foo.dgo`, `[1 2]`) },
		`expected one of ',' or '\]', got 2: \(file: foo\.dgo, line: 1, column: 4\)`)
}

func TestParse_value(t *testing.T) {
	require.Equal(t, vf.Map(), tf.Parse(`{}`))
}
