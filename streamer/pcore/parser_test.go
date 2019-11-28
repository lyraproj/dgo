package pcore_test

import (
	"fmt"
	"math"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/streamer/pcore"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func ExampleParse_qName() {
	defer tf.RemoveNamed(`foo.bar`)
	tf.NewNamed(`foo.bar`, nil, nil, nil, nil, nil)
	fmt.Println(pcore.Parse(`Foo::Bar`))
	// Output: foo.bar
}

func ExampleParse_int() {
	t := pcore.Parse(`23`)
	fmt.Println(t)
	// Output: 23
}

type booBar string

func (b booBar) String() string {
	return b.Type().(dgo.NamedType).ValueString(b)
}

func (b booBar) Type() dgo.Type {
	return tf.Named(`boo.bar`)
}

func (b booBar) Equals(other interface{}) bool {
	return b == other
}

func (b booBar) HashCode() int {
	return util.StringHash(string(b))
}

func ExampleParse_entry() {
	defer tf.RemoveNamed(`foo.bar`)
	defer tf.RemoveNamed(`boo.bar`)
	defer tf.RemoveNamed(`baz`)
	tf.NewNamed(`foo.bar`, nil, nil, nil, nil, nil)
	tf.NewNamed(`boo.bar`,
		func(value dgo.Value) dgo.Value {
			a := value.(dgo.Arguments)
			return booBar(a.Get(0).(dgo.String).GoString())
		},
		func(value dgo.Value) dgo.Value {
			return vf.String(string(value.(booBar)))
		}, nil, nil, nil)
	tf.NewNamed(`baz`, nil, nil, nil, nil, nil)
	const src = `# This is scanned code.
    constants => {
      first => 0,
      second => 0x32,
      third => 2e4,
      fourth => 2.3e-2,
      fifth => 'hello',
      sixth => "world",
      type => Foo::Bar[1,'23',Baz[1,2]],
      value => "String\nWith \\Escape",
      array => [a, b, c],
      call => Boo::Bar("with args")
    }
  `
	v := typ.ExactValue(pcore.Parse(src))
	sb := util.NewIndenter(`  `)
	sb.AppendValue(v)
	fmt.Println(sb.String())
	// Output:
	// {
	//   "constants": {
	//     "first": 0,
	//     "second": 50,
	//     "third": 20000.0,
	//     "fourth": 0.023,
	//     "fifth": "hello",
	//     "sixth": "world",
	//     "type": foo.bar[1,"23",type[baz[1,2]]],
	//     "value": "String\nWith \\Escape",
	//     "array": {
	//       "a",
	//       "b",
	//       "c"
	//     },
	//     "call": boo.bar"with args"
	//   }
	// }
}

func ExampleParse_hash() {
	v := pcore.Parse(`{value=>-1}`)
	fmt.Println(v)
	// Output: {"value":-1}
}

func TestParse_emptyTypeArgs(t *testing.T) {
	require.Panic(t,
		func() { pcore.Parse(`String[]`) },
		`illegal number of arguments for String\. Expected 1 or 2, got 0: \(column: 8\)`)
}

func TestParse_default(t *testing.T) {
	require.Equal(t, typ.Any, pcore.Parse(`Any`))
	require.Equal(t, typ.Array, pcore.Parse(`Array`))
	require.Equal(t, typ.Binary, pcore.Parse(`Binary`))
	require.Equal(t, typ.Boolean, pcore.Parse(`Boolean`))
	require.Equal(t, typ.Number, pcore.Parse(`Number`))
	require.Equal(t, typ.Integer, pcore.Parse(`Integer`))
	require.Equal(t, typ.True, pcore.Parse(`True`))
	require.Equal(t, typ.False, pcore.Parse(`False`))
	require.Equal(t, typ.Float, pcore.Parse(`Float`))
	require.Equal(t, typ.Map, pcore.Parse(`Hash`))
	require.Equal(t, typ.Nil, pcore.Parse(`Undef`))
	require.Equal(t, typ.Sensitive, pcore.Parse(`Sensitive`))
	require.Equal(t, typ.String, pcore.Parse(`Enum`))
	require.Equal(t, typ.String, pcore.Parse(`Pattern`))
	require.Equal(t, typ.String, pcore.Parse(`String`))
	require.Equal(t, typ.Type, pcore.Parse(`Type`))
}

func TestParse_exact(t *testing.T) {
	require.Equal(t, vf.Value(32), pcore.Parse(`32`))
	require.Equal(t, vf.Value(3.14), pcore.Parse(`3.14`))
	require.Equal(t, vf.Value("pelle"), pcore.Parse(`"pelle"`))
	require.Equal(t, vf.Values("pelle", 3.14, typ.Boolean), pcore.Parse(`["pelle", 3.14, Boolean]`))

	st := pcore.Parse(`Struct[a => Integer[1, 5], Optional[b] => Integer[2,2]]`)
	require.Equal(t, tf.StructMap(false,
		tf.StructMapEntry(`a`, tf.Integer(1, 5, true), true),
		tf.StructMapEntry(`b`, vf.Value(2).Type(), false)), st)
}

func TestParse_func(t *testing.T) {
	tt := pcore.Parse(`Callable[[String,Any,1],String]`)
	require.Equal(t, tf.Function(
		tf.VariadicTuple(typ.String, typ.Any),
		tf.Tuple(typ.String)), tt)
}

func testA(m string, e error) (string, error) {
	return m, e
}

func testB(m string, vs ...int) string {
	return fmt.Sprintf(`%s with %v`, m, vs)
}

func TestFunctionType(t *testing.T) {
	ft := tf.Function(typ.Tuple, typ.Tuple)
	require.Equal(t, pcore.ParseType(`Callable`), ft)

	ft = tf.Function(typ.EmptyTuple, tf.Tuple(typ.Any))
	require.Equal(t, pcore.ParseType(`Callable[0,0]`), ft)

	ft = pcore.Parse(`Callable[String,String,1]`).(dgo.FunctionType)
	require.Assignable(t, typ.Any, ft)
	require.NotAssignable(t, ft, typ.Any)
	require.Assignable(t, ft, pcore.ParseType(`Callable[String,String,1]`))
	require.Instance(t, ft.Type(), ft)

	require.Panic(t, func() { ft.ReflectType() }, `unable to build`)
	require.Equal(t, ft, ft)
	require.Equal(t, ft, pcore.ParseType(`Callable[String,String,1]`))
	require.Equal(t, ft.String(), `func(string,...string) any`)
	require.NotEqual(t, ft, pcore.ParseType(`Callable[Integer,String,1]`))
	require.NotEqual(t, ft, pcore.ParseType(`Callable[String,Array[String],1]`))
	require.NotEqual(t, ft, pcore.ParseType(`Callable[[String,String,1], String]`))
	require.NotEqual(t, ft, pcore.ParseType(`Tuple[String,String,1]`))

	in := ft.In()
	require.Assignable(t, in, pcore.ParseType(`Tuple[String,String,1]`))
	require.Assignable(t, in, pcore.ParseType(`Array[String,1]`))
	require.NotAssignable(t, in, pcore.ParseType(`Array[String]`))
	require.Same(t, typ.String, in.ElementType())
	require.Same(t, typ.String, pcore.ParseType(`Tuple[String,0]`).(dgo.TupleType).ElementType())

	require.NotEqual(t, 0, ft.HashCode())
	require.Equal(t, ft.HashCode(), pcore.ParseType(`Callable[String,String,1]`).HashCode())

	ft = pcore.ParseType(`Callable[[String,Integer,1], String]`).(dgo.FunctionType)
	require.Instance(t, ft, testB)
	require.NotInstance(t, ft, testA)
	require.NotInstance(t, ft, ft)
}

func TestParse_tuple(t *testing.T) {
	tt := pcore.Parse(`Tuple[String,String[10],1]`)
	require.Equal(t, tf.VariadicTuple(typ.String, tf.String(10)), tt)
}

func TestParse_ciEnum(t *testing.T) {
	st := pcore.Parse(`Enum[foo, fee, true]`)
	require.Equal(t, tf.CiEnum(`foo`, `fee`), st)
	require.Panic(t, func() { pcore.Parse(`Enum[32]`) }, `illegal argument for Enum. Expected string, got 32`)
}

func TestParse_sized(t *testing.T) {
	require.Equal(t, tf.String(1), pcore.Parse(`String[1]`))
	require.Equal(t, tf.String(1, 10), pcore.Parse(`String[1,10]`))
	require.Equal(t, tf.Array(1), pcore.Parse(`Array[Any,1]`))
	require.Equal(t, tf.Map(1), pcore.Parse(`Hash[Any,Any,1]`))
	require.Equal(t, tf.Pattern(regexp.MustCompile(`a.*`)), pcore.Parse(`Pattern[/a.*/]`))
}

func TestParse_nestedSized(t *testing.T) {
	require.Equal(t, tf.Array(tf.String(1), 1), pcore.Parse(`Array[String[1],1]`))
	require.Equal(t, tf.Array(tf.String(1, 10), 2, 5), pcore.Parse(`Array[String[1,10],2,5]`))
	require.Equal(t, tf.Map(typ.String, tf.String(1), 1), pcore.Parse(`Hash[String,String[1],1]`))
	require.Equal(t, tf.Map(typ.String, tf.String(1, 10), 2, 5), pcore.Parse(`Hash[String,String[1,10],2,5]`))
	require.Equal(t, tf.Map(tf.Map(typ.String, typ.Integer), tf.String(1, 10), 2, 5),
		pcore.Parse(`Hash[Hash[String,Integer],String[1,10],2,5]`))
}

func TestParse_aliasBad(t *testing.T) {
	require.Panic(t,
		func() { pcore.Parse(`f=Hash[String,Integer]`) }, `expected end of expression, got '='`)
	require.Panic(t,
		func() { pcore.Parse(`type m=Hash[String,Variant[Integer,M]]`) }, `expected end of expression, got m`)
	require.Panic(t,
		func() { pcore.Parse(`type Integer=Hash[String,Integer]`) }, `attempt to redeclare identifier 'Integer'`)
}

func TestParse_aliasInUnary(t *testing.T) {
	tp := pcore.Parse(`Type[type M=Hash[String,Variant[String,M]]]`).(dgo.UnaryType)
	require.Equal(t, `type[map[string](string|<recursive self reference to map type>)]`, tp.String())
}

func TestParse_aliasInAnyOf(t *testing.T) {
	tp := pcore.Parse(`type M=Array[Variant[Integer[0, 5], Integer[3, 8], M]]`).(dgo.Type)
	require.Equal(t, `[](0..5|3..8|<recursive self reference to slice type>)`, tp.String())
}

func TestParse_aliasInMap(t *testing.T) {
	tp := pcore.Parse(`type M=Array[Variant[Integer[0, 5], Integer[3, 8], M]]`).(dgo.Type)
	require.Equal(t, `[](0..5|3..8|<recursive self reference to slice type>)`, tp.String())
}

func TestParse_range(t *testing.T) {
	require.Equal(t, tf.Integer(1, 10, true), pcore.Parse(`Integer[1,10]`))
	require.Equal(t, tf.Integer(1, math.MaxInt64, true), pcore.Parse(`Integer[1]`))
	require.Equal(t, tf.Integer(math.MinInt64, 0, true), pcore.Parse(`Integer[default,0]`))
	require.Equal(t, tf.Float(1, 10, true), pcore.Parse(`Float[1.0,10]`))
	require.Equal(t, tf.Float(1, 10, true), pcore.Parse(`Float[1,10.0]`))
	require.Equal(t, tf.Float(1, 10, true), pcore.Parse(`Float[1.0,10.0]`))
	require.Equal(t, tf.Float(1, math.MaxFloat64, true), pcore.Parse(`Float[1.0]`))
	require.Equal(t, tf.Float(-math.MaxFloat64, 0, true), pcore.Parse(`Float[default,0.0]`))

	require.Panic(t, func() { pcore.Parse(`Float[b]`) }, `illegal argument for Float. Expected int|float, got b`)
}

func TestParse_unary(t *testing.T) {
	require.Equal(t, typ.String.Type(), pcore.Parse(`Type[String]`))
	require.Equal(t, typ.Any.Type(), pcore.Parse(`Type`))
	require.Panic(t, func() { pcore.Parse(`Type[String`) }, `expected one of ',' or ']', got EOT`)
}

func TestParse_ternary(t *testing.T) {
	require.Equal(t, typ.Number, pcore.Parse(`Variant[Integer,Float]`))
	require.Equal(t, tf.AnyOf(typ.String, typ.Integer, typ.Float), pcore.Parse(`Variant[String,Integer,Float]`))
	require.Equal(t, tf.Enum(`a`, `b`, `c`), pcore.Parse(`Enum[a, b, c]`))
}

func TestParse_string(t *testing.T) {
	require.Equal(t, vf.String("\r"), pcore.Parse(`"\r"`))
	require.Equal(t, vf.String("\n"), pcore.Parse(`"\n"`))
	require.Equal(t, vf.String("\t"), pcore.Parse(`"\t"`))
	require.Equal(t, vf.String("\""), pcore.Parse(`"\""`))

	require.Panic(t, func() { pcore.Parse(`"\`) }, `unterminated string`)
	require.Panic(t, func() { pcore.Parse(`"\"`) }, `unterminated string`)
	require.Panic(t, func() { pcore.Parse(`"\y"`) }, `illegal escape`)
	require.Panic(t, func() {
		pcore.Parse(`["
"`)
	}, `unterminated string`)
}

func TestParse_multiAliases(t *testing.T) {
	tp := pcore.Parse(`
Struct[
  types => [
    type Ascii=Integer[1,127],
    type Slug=Pattern[/^[a-z0-9-]+$/]
  ],
  x => Hash[Slug,Struct["Token" => Ascii, value => String]]
]`).(dgo.StructMapType)
	require.Equal(t, `map[/^[a-z0-9-]+$/]{"Token":1..127,"value":string}`, tp.Get(`x`).Value().(dgo.Type).String())
}

func TestParse_errors(t *testing.T) {
	require.Panic(t, func() { pcore.Parse(`Integer[1 23]`) }, `expected one of ',' or '\]', got 23: \(column: 11\)`)
	require.Panic(t, func() { pcore.Parse(`Float[1. 3]`) }, `unexpected character ' '`)
	require.Panic(t, func() { pcore.Parse(`String[1 . 3]`) }, `expected one of ',' or '\]', got '.'`)
	require.Panic(t, func() { pcore.Parse(`/\`) }, `unterminated regexp`)
	require.Panic(t, func() { pcore.Parse(`/\/`) }, `unterminated regexp`)
	require.Panic(t, func() { pcore.Parse(`Pattern[/\//`) }, `expected one of ',' or '\]', got EOT`)
	require.Panic(t, func() { pcore.Parse(`Pattern[/\t/`) }, `expected one of ',' or '\]', got EOT`)
}

func TestParseFile_errors(t *testing.T) {
	require.Panic(t,
		func() { pcore.ParseFile(nil, `foo.dgo`, `[1 2]`) },
		`expected one of ',' or '\]', got 2: \(file: foo\.dgo, line: 1, column: 4\)`)
}
