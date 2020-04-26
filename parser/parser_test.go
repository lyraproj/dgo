package parser_test

import (
	"math"
	"math/big"
	"regexp"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
	"github.com/tada/dgo/stringer"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/test/require"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestParse_default(t *testing.T) {
	assert.Equal(t, typ.Any, tf.ParseType(`any`))
	assert.Equal(t, typ.Boolean, tf.ParseType(`bool`))
	assert.Equal(t, typ.Integer, tf.ParseType(`int`))
	assert.Equal(t, typ.True, tf.ParseType(`true`))
	assert.Equal(t, typ.False, tf.ParseType(`false`))
	assert.Equal(t, typ.Float, tf.ParseType(`float`))
	assert.Equal(t, typ.String, tf.ParseType(`string`))
	assert.Equal(t, typ.Binary, tf.ParseType(`binary`))
	assert.Equal(t, typ.Nil, tf.ParseType(`nil`))
	assert.Equal(t, typ.Array, tf.ParseType(`[]any`))
	assert.Equal(t, typ.Map, tf.ParseType(`map[any]any`))
}

func TestParse_exact(t *testing.T) {
	assert.Equal(t, vf.Value(32).Type(), tf.ParseType(`32`))
	assert.Equal(t, vf.Value(3.14).Type(), tf.ParseType(`3.14`))
	assert.Equal(t, vf.Value("pelle").Type(), tf.ParseType(`"pelle"`))
	assert.Equal(t, tf.Tuple(vf.Value("pelle").Type(), vf.Value(3.14).Type(), typ.Boolean), tf.ParseType(`{"pelle", 3.14, bool}`))

	st := tf.ParseType(`{"a":1..5,"b"?:2}`)
	assert.Equal(t, tf.StructMap(false,
		tf.StructMapEntry(`a`, tf.Integer64(1, 5, true), true),
		tf.StructMapEntry(`b`, vf.Value(2).Type(), false)), st)

	st = tf.ParseType(`{"a":1..5,"b"?:2,...}`)
	assert.Equal(t, tf.StructMap(true,
		tf.StructMapEntry(`a`, tf.Integer64(1, 5, true), true),
		tf.StructMapEntry(`b`, vf.Value(2).Type(), false)), st)

	assert.Equal(t, `{"a":1..5,"b"?:2,...}`, st.String())

	st = tf.ParseType(`{}`)
	assert.Equal(t, tf.StructMap(false), st)
	assert.Equal(t, `{}`, st.String())

	st = tf.ParseType(`{...}`)
	assert.Equal(t, tf.StructMap(true), st)
	assert.Equal(t, `{...}`, st.String())
}

func TestParse_func(t *testing.T) {
	tt := tf.ParseType(`func(string,...any) (string, bool)`)
	assert.Equal(t, tf.Function(
		tf.VariadicTuple(typ.String, typ.Any),
		tf.Tuple(typ.String, typ.Boolean)), tt)
}

func TestParse_tuple(t *testing.T) {
	assert.Equal(t, tf.VariadicTuple(typ.String, tf.String(10)), tf.ParseType(`{string,...string[10]}`))
	assert.Equal(t, tf.Function(typ.EmptyTuple, typ.EmptyTuple), tf.ParseType(`func()`))
}

func TestParse_ciEnum(t *testing.T) {
	st := tf.ParseType(`~"foo"|~"fee"`)
	assert.Equal(t, tf.CiEnum(`foo`, `fee`), st)
	assert.Panic(t, func() { tf.ParseType(`~32`) }, `expected a literal string, got 32`)
}

func TestParse_sized(t *testing.T) {
	assert.Equal(t, tf.String(1), tf.ParseType(`string[1]`))
	assert.Equal(t, tf.String(1, 10), tf.ParseType(`string[1,10]`))
	assert.Equal(t, tf.Array(1), tf.ParseType(`[1]any`))
	assert.Equal(t, tf.Map(1), tf.ParseType(`map[any,1]any`))
	assert.Equal(t, tf.Pattern(regexp.MustCompile(`a.*`)), tf.ParseType(`/a.*/`))
}

func TestParse_nestedSized(t *testing.T) {
	assert.Equal(t, tf.Array(tf.String(1), 1), tf.ParseType(`[1]string[1]`))
	assert.Equal(t, tf.Array(tf.String(1, 10), 2, 5), tf.ParseType(`[2,5]string[1,10]`))
	assert.Equal(t, tf.Map(typ.String, tf.String(1), 1), tf.ParseType(`map[string,1]string[1]`))
	assert.Equal(t, tf.Map(typ.String, tf.String(1, 10), 2, 5), tf.ParseType(`map[string,2,5]string[1,10]`))
	assert.Equal(t, tf.Map(tf.Map(typ.String, typ.Integer), tf.String(1, 10), 2, 5),
		tf.ParseType(`map[map[string]int,2,5]string[1,10]`))
}

func TestParse_aliasBad(t *testing.T) {
	internal.ResetDefaultAliases()
	assert.Panic(t, func() { tf.ParseType(`"f"=map[string]int`) }, `expected end of expression, got '='`)
	internal.ResetDefaultAliases()
	assert.Panic(t, func() { tf.ParseType(`m=map[string](int|<"m">)`) }, `expected an identifier, got "m"`)
	internal.ResetDefaultAliases()
	assert.Panic(t, func() { tf.ParseType(`m=map[string](int|<m)`) }, `expected '>', got '\)'`)
	internal.ResetDefaultAliases()
	assert.Panic(t, func() { tf.ParseType(`map[string](int|n)`) }, `reference to unresolved type 'n'`)
	internal.ResetDefaultAliases()
	assert.Panic(t, func() { tf.ParseType(`m=map[string](int|n)`) }, `reference to unresolved type 'n'`)
	internal.ResetDefaultAliases()
	assert.Panic(t, func() { tf.ParseType(`int=map[string](int|int)`) }, `attempt to redeclare identifier 'int'`)
}

func TestParse_aliasInUnary(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`type[m=map[string](string|m)]`).(dgo.UnaryType)
	assert.Equal(t, `type[m]`, tp.String())

	internal.ResetDefaultAliases()
	tp = tf.ParseType(`!m=map[string](string|m)`).(dgo.UnaryType)
	assert.Equal(t, `!m`, tp.String())
}

func TestParse_aliasInAllOf(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`m=[](0..5&3..8&m)`)
	assert.Equal(t, `m`, tp.String())
}

func TestParse_aliasInAllOf_separateAm(t *testing.T) {
	internal.ResetDefaultAliases()
	tf.BuiltInAliases().Collect(func(a dgo.AliasAdder) {
		tp := tf.ParseFile(a, `internal`, `m=[](0..5&3..8&m)`)
		assert.Equal(t, `[](0..5&3..8&<recursive self reference to slice type>)`, tp.String())
	})
}

func TestParse_aliasInOneOf(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`m=[](0..5^3..8^m)`)
	assert.Equal(t, `m`, tp.String())
}

func TestParse_data(t *testing.T) {
	tf.BuiltInAliases().Collect(func(aa dgo.AliasAdder) {
		tp := tf.ParseFile(aa, `internal`, `data`)
		assert.Equal(t, `data`, tp.String())
	})
}

func TestParse_range(t *testing.T) {
	assert.Equal(t, tf.Integer64(1, 10, true), tf.ParseType(`1..10`))
	assert.Equal(t, tf.Integer64(1, 10, false), tf.ParseType(`1...10`))
	assert.Equal(t, tf.Integer64(1, math.MaxInt64, true), tf.ParseType(`1..`))
	assert.Equal(t, tf.Integer64(1, math.MaxInt64, false), tf.ParseType(`1...`))
	assert.Equal(t, tf.Integer64(math.MinInt64, 0, true), tf.ParseType(`..0`))
	assert.Equal(t, tf.Integer64(math.MinInt64, 0, false), tf.ParseType(`...0`))
	assert.Equal(t, tf.Float64(1, 10, true), tf.ParseType(`1.0..10`))
	assert.Equal(t, tf.Float64(1, 10, true), tf.ParseType(`1..10.0`))
	assert.Equal(t, tf.Float64(1, 10, true), tf.ParseType(`1.0..10.0`))
	assert.Equal(t, tf.Float64(1, 10, false), tf.ParseType(`1.0...10`))
	assert.Equal(t, tf.Float64(1, 10, false), tf.ParseType(`1...10.0`))
	assert.Equal(t, tf.Float64(1, 10, false), tf.ParseType(`1.0...10.0`))
	assert.Equal(t, tf.Float64(1, math.MaxFloat64, true), tf.ParseType(`1.0..`))
	assert.Equal(t, tf.Float64(1, math.MaxFloat64, false), tf.ParseType(`1.0...`))
	assert.Equal(t, tf.Float64(-math.MaxFloat64, 0, true), tf.ParseType(`..0.0`))
	assert.Equal(t, tf.Float64(-math.MaxFloat64, 0, false), tf.ParseType(`...0.0`))

	assert.Panic(t, func() { tf.ParseType(`.."b"`) }, `expected an integer or a float, got "b"`)
	assert.Panic(t, func() { tf.ParseType(`../a*/`) }, `expected an integer or a float, got /a\*/`)
	assert.Panic(t, func() { tf.ParseType(`..."b"`) }, `expected an integer or a float, got "b"`)
	assert.Panic(t, func() { tf.ParseType(`.../a*/`) }, `expected an integer or a float, got /a\*/`)
}

func TestParse_unary(t *testing.T) {
	assert.Equal(t, tf.Not(typ.String), tf.ParseType(`!string`))
	assert.Equal(t, typ.String.Type(), tf.ParseType(`type[string]`))
	assert.Equal(t, typ.Any.Type(), tf.ParseType(`type`))
	assert.Panic(t, func() { tf.ParseType(`type[string`) }, `expected ']', got EOT`)
}

func TestParse_ternary(t *testing.T) {
	assert.Equal(t, typ.Number, tf.ParseType(`int|float`))
	assert.Equal(t, tf.AnyOf(typ.String, typ.Integer, typ.Float), tf.ParseType(`string|int|float`))
	assert.Equal(t, tf.OneOf(typ.String, typ.Integer, typ.Float), tf.ParseType(`string^int^float`))
	assert.Equal(t, tf.AllOf(typ.String, typ.Integer, typ.Float), tf.ParseType(`string&int&float`))
	assert.Equal(t, tf.AnyOf(typ.String, tf.OneOf(
		tf.AllOf(typ.Integer, typ.Float), typ.Boolean)), tf.ParseType(`string|int&float^bool`))
	assert.Equal(t, tf.AnyOf(typ.String,
		tf.AllOf(typ.Integer, tf.OneOf(typ.Float, typ.Boolean))), tf.ParseType(`string|int&(float^bool)`))
	assert.Equal(t, tf.Enum(`a`, `b`, `c`), tf.ParseType(`"a"|"b"|"c"`))
}

func TestParse_string(t *testing.T) {
	assert.Equal(t, vf.String("\r").Type(), tf.ParseType(`"\r"`))
	assert.Equal(t, vf.String("\n").Type(), tf.ParseType(`"\n"`))
	assert.Equal(t, vf.String("\t").Type(), tf.ParseType(`"\t"`))
	assert.Equal(t, vf.String("\"").Type(), tf.ParseType(`"\""`))

	assert.Equal(t, vf.String("\"").Type(), tf.ParseType("`\"`"))

	assert.Panic(t, func() { tf.ParseType(`"\`) }, `unterminated string`)
	assert.Panic(t, func() { tf.ParseType(`"\"`) }, `unterminated string`)
	assert.Panic(t, func() { tf.ParseType(`"\y"`) }, `illegal escape`)
	assert.Panic(t, func() {
		tf.ParseType(`["
"`)
	}, `unterminated string`)

	assert.Panic(t, func() { tf.ParseType("`x") }, `unterminated string`)
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

	assert.Equal(t, `map[slug]{"Token":ascii,"value":string}`, stringer.TypeStringWithAliasMap(tp.GetEntryType(`x`).Value().(dgo.Type), am))
}

func TestParse_mapKeyAliases(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`{tp = "key", { key: 2 }}`)
	assert.Equal(t, `{tp,{tp:2}}`, stringer.TypeString(tp))
}

func TestParse_errors(t *testing.T) {
	assert.Panic(t, func() { tf.ParseType(`[1 23]`) }, `expected one of ',' or '\]', got 23: \(column: 4\)`)
	assert.Panic(t, func() { tf.ParseType(`map{}`) }, `expected '\[', got '\{': \(column: 4\)`)
	assert.Panic(t, func() { tf.ParseType(`map[string}string`) }, `expected '\]', got '\}': \(column: 11\)`)
	assert.Panic(t, func() { tf.ParseType(`map[string](string|int]`) }, `expected '\)', got '\]': \(column: 23\)`)
	assert.Panic(t, func() { tf.ParseType(`{1 23}`) }, `expected one of ',' or '\}', got 23`)
	assert.Panic(t, func() { tf.ParseType(`[}int`) }, `expected a type expression, got '\}'`)
	assert.Panic(t, func() { tf.ParseType(`[]string[1][2]`) }, `expected end of expression, got '\['`)
	assert.Panic(t, func() { tf.ParseType(`[]string[1`) }, `expected one of ',' or '\]', got EOT`)
	assert.Panic(t, func() { tf.ParseType(`apple`) }, `reference to unresolved type 'apple'`)
	assert.Panic(t, func() { tf.ParseType(`-two`) }, `unexpected character 't'`)
	assert.Panic(t, func() { tf.ParseType(`[3e,0]`) }, `unexpected character ','`)
	assert.Panic(t, func() { tf.ParseType(`[3e1.0]`) }, `unexpected character '.'`)
	assert.Panic(t, func() { tf.ParseType(`[3e23r,4]`) }, `unexpected character 'r'`)
	assert.Panic(t, func() { tf.ParseType(`[3e1`) }, `expected one of ',' or ']', got EOT`)
	assert.Panic(t, func() { tf.ParseType(`[3e`) }, `unexpected end`)
	assert.Panic(t, func() { tf.ParseType(`[0x`) }, `unexpected end`)
	assert.Panic(t, func() { tf.ParseType(`[0x4`) }, `expected one of ',' or ']', got EOT`)
	assert.Panic(t, func() { tf.ParseType(`[1. 3]`) }, `unexpected character ' '`)
	assert.Panic(t, func() { tf.ParseType(`[1 . 3]`) }, `expected one of ',' or ']', got '.'`)
	assert.Panic(t, func() { tf.ParseType(`[/\`) }, `unterminated regexp`)
	assert.Panic(t, func() { tf.ParseType(`[/\/`) }, `unterminated regexp`)
	assert.Panic(t, func() { tf.ParseType(`[/\//`) }, `expected one of ',' or ']', got EOT`)
	assert.Panic(t, func() { tf.ParseType(`[/\t/`) }, `expected one of ',' or ']', got EOT`)
	assert.Panic(t, func() { tf.ParseType(`{"a":string,...`) }, `expected '}', got EOT`)
	assert.Panic(t, func() { tf.ParseType(`{a:32, 4}`) }, `mix of elements and map entries`)
	assert.Panic(t, func() { tf.ParseType(`{"a":32, 4}`) }, `mix of elements and map entries`)
	assert.Panic(t, func() { tf.ParseType(`{4, a:32}`) }, `mix of elements and map entries`)
	assert.Panic(t, func() { tf.ParseType(`{4, "a":32}`) }, `mix of elements and map entries`)
	assert.Panic(t, func() { tf.ParseType(`{4, a}`) }, `reference to unresolved type 'a'`)
	assert.Panic(t, func() { tf.ParseType(`{func, 3}`) }, `expected '\(', got ','`)
}

func TestParseFile_errors(t *testing.T) {
	assert.Panic(t,
		func() { tf.ParseFile(nil, `foo.dgo`, `[1 2]`) },
		`expected one of ',' or '\]', got 2: \(file: foo\.dgo, line: 1, column: 4\)`)
}

func TestParse_value(t *testing.T) {
	assert.Equal(t, vf.Map(), tf.Parse(`{}`))
}

func TestParse_bigInt_hex(t *testing.T) {
	v, ok := tf.Parse(`big 0x10000000000000000`).(dgo.BigInt)
	require.True(t, ok)
	e, _ := new(big.Int).SetString(`10000000000000000`, 16)
	assert.Equal(t, e, v)
}

func TestParse_maxUint64_hex(t *testing.T) {
	v, ok := tf.Parse(`big 0xffffffffffffffff`).(dgo.BigInt)
	require.True(t, ok)
	e := new(big.Int).SetUint64(uint64(math.MaxUint64))
	assert.Equal(t, e, v)
}

func TestParse_maxInt64_hex(t *testing.T) {
	v := tf.Parse(`0x7fffffffffffffff`)
	_, ok := v.(dgo.BigInt)
	assert.False(t, ok)
	assert.Equal(t, math.MaxInt64, v)
}

func TestParse_bigInt_decimal(t *testing.T) {
	v, ok := tf.Parse(`big -170141183460469231731687303715884105728`).(dgo.BigInt)
	require.True(t, ok)
	e, _ := new(big.Int).SetString(`-170141183460469231731687303715884105728`, 10)
	assert.Equal(t, e, v)
}

func TestParse_bigInt_inclusive_range(t *testing.T) {
	v, ok := tf.Parse(`big -170141183460469231731687303715884105728..170141183460469231731687303715884105727`).(dgo.IntegerType)
	require.True(t, ok)
	min, _ := new(big.Int).SetString(`-170141183460469231731687303715884105728`, 10)
	max, _ := new(big.Int).SetString(`170141183460469231731687303715884105727`, 10)
	assert.Equal(t, tf.Integer(vf.BigInt(min), vf.BigInt(max), true), v)
}

func TestParse_bigInt_exclusive_range(t *testing.T) {
	v, ok := tf.Parse(`big -170141183460469231731687303715884105728...170141183460469231731687303715884105727`).(dgo.IntegerType)
	require.True(t, ok)
	min, _ := new(big.Int).SetString(`-170141183460469231731687303715884105728`, 10)
	max, _ := new(big.Int).SetString(`170141183460469231731687303715884105727`, 10)
	assert.Equal(t, tf.Integer(vf.BigInt(min), vf.BigInt(max), false), v)
}

func TestParse_bigInt_float_range(t *testing.T) {
	v, ok := tf.Parse(`big -170141183460469231731687303715884105728..170141183460469231731687303715884105727.123`).(dgo.FloatType)
	require.True(t, ok)
	min, _ := new(big.Float).SetString(`-170141183460469231731687303715884105728`)
	max, _ := new(big.Float).SetString(`170141183460469231731687303715884105727.123`)
	assert.Equal(t, tf.Float(vf.BigFloat(min), vf.BigFloat(max), true), v)
}

func TestParse_big_float(t *testing.T) {
	v, ok := tf.Parse(`big -170141183460469231731687303715884105728.123`).(dgo.BigFloat)
	require.True(t, ok)
	e, _ := new(big.Float).SetString(`-170141183460469231731687303715884105728.123`)
	assert.Equal(t, e, v)
}

func TestParse_big_string(t *testing.T) {
	require.Panic(t, func() {
		_, _ = tf.Parse(`big "obviously bogus"`).(dgo.BigFloat)
	}, `expected an integer or a float`)
}

func TestParse_bigFloat_inclusive_range(t *testing.T) {
	v, ok := tf.Parse(`big -170141183460469231731687303715884105728.123..170141183460469231731687303715884105727.123`).(dgo.FloatType)
	require.True(t, ok)
	min, _ := new(big.Float).SetString(`-170141183460469231731687303715884105728.123`)
	max, _ := new(big.Float).SetString(`170141183460469231731687303715884105727.123`)
	assert.Equal(t, tf.Float(vf.BigFloat(min), vf.BigFloat(max), true), v)
}

func TestParse_bigFloat_exclusive_range(t *testing.T) {
	v, ok := tf.Parse(`big -170141183460469231731687303715884105728.123...170141183460469231731687303715884105727.123`).(dgo.FloatType)
	require.True(t, ok)
	min, _ := new(big.Float).SetString(`-170141183460469231731687303715884105728.123`)
	max, _ := new(big.Float).SetString(`170141183460469231731687303715884105727.123`)
	assert.Equal(t, tf.Float(vf.BigFloat(min), vf.BigFloat(max), false), v)
}

func TestParse_inclusive_range(t *testing.T) {
	v, ok := tf.Parse(`-3..3`).(dgo.IntegerType)
	require.True(t, ok)
	assert.Equal(t, tf.Integer64(-3, 3, true), v)
}

func TestParse_exclusive_range(t *testing.T) {
	v, ok := tf.Parse(`big -170141183460469231731687303715884105728...170141183460469231731687303715884105727`).(dgo.IntegerType)
	require.True(t, ok)
	min, _ := new(big.Int).SetString(`-170141183460469231731687303715884105728`, 10)
	max, _ := new(big.Int).SetString(`170141183460469231731687303715884105727`, 10)
	assert.Equal(t, tf.Integer(vf.BigInt(min), vf.BigInt(max), false), v)
}
