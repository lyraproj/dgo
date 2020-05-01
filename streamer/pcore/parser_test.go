package pcore_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/streamer/pcore"
	"github.com/lyraproj/dgo/test/require"
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

func (b booBar) HashCode() int32 {
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
      literals => [true, false, undef],
      type => Foo::Bar[1,'23',Baz[1,2]],
      value => "String\nWith \\Escape",
      array => [a, b, c],
      call => Boo::Bar("with args")
    }
  `
	v := pcore.Parse(src)
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
	//     "literals": {
	//       true,
	//       false,
	//       nil
	//     },
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
	fmt.Println(v.String())
	// Output: {"value":-1}
}

func TestParse_call(t *testing.T) {
	v := pcore.Parse(`String(23)`)
	require.Equal(t, vf.New(typ.String, vf.Integer(23)), v)
}

func TestParse_emptyHash(t *testing.T) {
	v := pcore.Parse(`{}`)
	require.Equal(t, vf.Map(), v)
}

func TestParse_emptyTypeArgs(t *testing.T) {
	require.Panic(t,
		func() { pcore.Parse(`String[]`) },
		`illegal number of arguments for String\. Expected 1 or 2, got 0: \(column: 8\)`)
}

func TestParse_identifier(t *testing.T) {
	v := pcore.Parse(`m`)
	require.Equal(t, `m`, v)
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
	require.Equal(t, pcore.Parse(`Callable`), ft)
	require.Equal(t, ft.String(), `func(...any) (...any)`)

	ft = tf.Function(typ.EmptyTuple, tf.Tuple(typ.Any))
	require.Equal(t, ft.String(), `func() any`)
	require.Equal(t, pcore.Parse(`Callable[[0,0],Any]`), ft)

	ft = pcore.Parse(`Callable[String,String,1]`).(dgo.FunctionType)
	require.Assignable(t, typ.Any, ft)
	require.NotAssignable(t, ft, typ.Any)
	require.Assignable(t, ft, pcore.Parse(`Callable[String,String,1]`).(dgo.Type))
	require.Instance(t, ft.Type(), ft)

	require.Panic(t, func() { ft.ReflectType() }, `unable to build`)
	require.Equal(t, ft, ft)
	require.Equal(t, ft, pcore.Parse(`Callable[String,String,1]`))
	require.Equal(t, ft.String(), `func(string,...string)`)
	require.NotEqual(t, ft, pcore.Parse(`Callable[Integer,String,1]`))
	require.NotEqual(t, ft, pcore.Parse(`Callable[String,Array[String],1]`))
	require.NotEqual(t, ft, pcore.Parse(`Callable[[String,String,1], String]`))
	require.NotEqual(t, ft, pcore.Parse(`Tuple[String,String,1]`))

	in := ft.In()
	require.Assignable(t, in, pcore.Parse(`Tuple[String,String,1]`).(dgo.Type))
	require.Assignable(t, in, pcore.Parse(`Array[String,1]`).(dgo.Type))
	require.NotAssignable(t, in, pcore.Parse(`Array[String]`).(dgo.Type))
	require.Same(t, typ.String, in.ElementType())
	require.Same(t, typ.String, pcore.Parse(`Tuple[String,0]`).(dgo.TupleType).ElementType())

	require.NotEqual(t, 0, ft.HashCode())
	require.Equal(t, ft.HashCode(), pcore.Parse(`Callable[String,String,1]`).HashCode())

	ft = pcore.Parse(`Callable[[String,Integer,1], String]`).(dgo.FunctionType)
	require.Instance(t, ft, testB)
	require.NotInstance(t, ft, testA)
	require.NotInstance(t, ft, ft)

	ft = pcore.Parse(`Callable[[String,1,1,Callable], String]`).(dgo.FunctionType)
	require.Equal(t, `func(func(...any) (...any),string) string`, ft.String())

	ft = pcore.Parse(`Callable[[String,default,1,Callable], String]`).(dgo.FunctionType)
	require.Equal(t, `func(func(...any) (...any),...string) string`, ft.String())

	ft = pcore.Parse(`Callable[[String,1,Callable], String]`).(dgo.FunctionType)
	require.Equal(t, `func(func(...any) (...any),...string) string`, ft.String())
}

func TestParse_ciEnum(t *testing.T) {
	st := pcore.Parse(`Enum[foo, fee, true]`)
	require.Equal(t, tf.CiEnum(`foo`, `fee`), st)
	require.Panic(t, func() { pcore.Parse(`Enum[32]`) }, `illegal argument for Enum. Expected string, got 32`)
}

func TestParse_sized(t *testing.T) {
	require.Equal(t, tf.String(1), pcore.Parse(`String[1]`))
	require.Equal(t, tf.String(1, 10), pcore.Parse(`String[1,10]`))
	require.Equal(t, tf.Array(1), pcore.Parse(`Array[1]`))
	require.Equal(t, tf.Array(1, 10), pcore.Parse(`Array[1,10]`))
	require.Equal(t, tf.Array(1), pcore.Parse(`Array[Any,1]`))
	require.Equal(t, tf.Binary(1), pcore.Parse(`Binary[1]`))
	require.Equal(t, tf.Binary(1, 10), pcore.Parse(`Binary[1,10]`))
	require.Equal(t, tf.Map(1), pcore.Parse(`Hash[1]`))
	require.Equal(t, tf.Map(1), pcore.Parse(`Hash[Any,Any,1]`))
	require.Equal(t, tf.Map(typ.String, typ.Integer, 1, 10), pcore.Parse(`Hash[String,Integer,1,10]`))
	require.Equal(t, tf.Map(1, 10), pcore.Parse(`Hash[1,10]`))
	require.Equal(t, tf.Map(typ.Any, typ.Any), pcore.Parse(`Hash[Any,Any]`))
	require.Equal(t, tf.Pattern(regexp.MustCompile(`a.*`)), pcore.Parse(`Pattern[/a.*/]`))
	require.Panic(t,
		func() { pcore.Parse(`Integer[2, x]`) }, `Expected int, got x`)
}

func TestParse_nestedSized(t *testing.T) {
	require.Equal(t, tf.Array(tf.String(1), 1), pcore.Parse(`Array[String[1],1]`))
	require.Equal(t, tf.Array(tf.String(1, 10), 2, 5), pcore.Parse(`Array[String[1,10],2,5]`))
	require.Equal(t, tf.Map(typ.String, tf.String(1), 1), pcore.Parse(`Hash[String,String[1],1]`))
	require.Equal(t, tf.Map(typ.String, tf.String(1, 10), 2, 5), pcore.Parse(`Hash[String,String[1,10],2,5]`))
	require.Equal(t, tf.Map(tf.Map(typ.String, typ.Integer), tf.String(1, 10), 2, 5),
		pcore.Parse(`Hash[Hash[String,Integer],String[1,10],2,5]`))
}

func TestParse_notUndef(t *testing.T) {
	require.Equal(t, tf.Not(typ.Nil), pcore.Parse(`NotUndef`))
	require.Equal(t, tf.Not(typ.Nil), pcore.Parse(`NotUndef[Any]`))
	require.Equal(t, tf.AllOf(tf.Not(typ.Nil), tf.AnyOf(typ.Nil, typ.String)), pcore.Parse(`NotUndef[Variant[Undef,String]]`))
}

func TestParse_number(t *testing.T) {
	require.Equal(t, tf.AnyOf(typ.Integer, typ.Float), pcore.Parse(`Number`))
	require.Equal(t, tf.AnyOf(tf.Integer(vf.Integer(3), nil, true), tf.Float(vf.Float(3), nil, true)), pcore.Parse(`Number[3]`))
	require.Equal(t, tf.AnyOf(tf.Integer(vf.Integer(3), vf.Integer(3), true), tf.Float(vf.Float(3), vf.Float(3), true)), pcore.Parse(`Number[3,3]`))
}

func TestParse_optional(t *testing.T) {
	require.Equal(t, typ.Any, pcore.Parse(`Optional`))
	require.Equal(t, tf.AnyOf(typ.Nil, typ.String), pcore.Parse(`Optional[String]`))
	require.Equal(t, tf.AnyOf(typ.Nil, vf.String(`hello`).Type()), pcore.Parse(`Optional[hello]`))
}

func TestParse_pattern(t *testing.T) {
	require.Equal(t, typ.String, pcore.Parse(`Pattern`))
	require.Equal(t, tf.Pattern(regexp.MustCompile(`a.*`)), pcore.Parse(`Pattern[/a.*/]`))
	require.Equal(t, tf.Pattern(regexp.MustCompile(`a.*`)), pcore.Parse(`Pattern["a.*"]`))
	require.Panic(t,
		func() { pcore.Parse(`Pattern[String]`) }, `a value of type type\[string\] cannot be assigned to a variable of type string\|regexp`)
}

func TestParse_regexp(t *testing.T) {
	require.Equal(t, typ.Regexp, pcore.Parse(`Regexp`))
	require.Equal(t, vf.Regexp(regexp.MustCompile(`a.*`)).Type(), pcore.Parse(`Regexp[/a.*/]`))
	require.Equal(t, vf.Regexp(regexp.MustCompile(`a.*`)).Type(), pcore.Parse(`Regexp["a.*"]`))
}

func TestParse_sensitive(t *testing.T) {
	require.Equal(t, typ.Sensitive, pcore.Parse(`Sensitive`))
	require.Equal(t, tf.Sensitive(typ.String), pcore.Parse(`Sensitive[String]`))
}

func TestParse_struct(t *testing.T) {
	require.Equal(t, tf.StructMap(true), pcore.Parse(`Struct`))
	require.Equal(t, tf.StructMap(false), pcore.Parse(`Struct[{}]`))

	st := pcore.Parse(`Struct[a => Integer[1, 5], Optional[b] => Integer[2,2]]`)
	require.Equal(t, tf.StructMap(false,
		tf.StructMapEntry(`a`, tf.Integer(vf.Integer(1), vf.Integer(5), true), true),
		tf.StructMapEntry(`b`, vf.Value(2).Type(), false)), st)
}

func TestParse_tuple(t *testing.T) {
	require.Equal(t, typ.Tuple, pcore.Parse(`Tuple`))
	require.Equal(t, tf.Tuple(typ.String), pcore.Parse(`Tuple[String]`))
	require.Equal(t, tf.Tuple(typ.String), pcore.Parse(`Tuple[String,1,1]`))
	require.Equal(t, tf.Tuple(typ.String), pcore.Parse(`Tuple[String,default,default]`))
	require.Equal(t, tf.VariadicTuple(typ.String), pcore.Parse(`Tuple[String,default,10]`))
	require.Equal(t, tf.VariadicTuple(typ.String), pcore.Parse(`Tuple[String,default]`))
	require.Equal(t, tf.VariadicTuple(typ.String), pcore.Parse(`Tuple[String,1]`))
	require.Equal(t, tf.VariadicTuple(typ.String, tf.String(10)), pcore.Parse(`Tuple[String,String[10],1]`))
}

func TestParse_variant(t *testing.T) {
	require.Equal(t, typ.AnyOf, pcore.Parse(`Variant`))
	require.Equal(t, tf.AnyOf(2, 3), pcore.Parse(`Variant[2, 3]`))
}

func TestParse_aliasBad(t *testing.T) {
	internal.ResetDefaultAliases()
	require.Panic(t,
		func() { pcore.Parse(`f=Hash[String,Integer]`) }, `expected end of expression, got '='`)
	require.Panic(t,
		func() { pcore.Parse(`type m=Hash[String,Variant[Integer,M]]`) }, `expected end of expression, got m`)
	require.Panic(t,
		func() { pcore.Parse(`type Integer=Hash[String,Integer]`) }, `attempt to redeclare identifier 'Integer'`)
	require.Panic(t, func() { pcore.Parse(`MyType`) }, `reference to unresolved type 'myType'`)
}

func TestParse_aliasInUnary(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := pcore.Parse(`Type[type M=Hash[String,Variant[String,M]]]`).(dgo.UnaryType)
	require.Equal(t, `type[m]`, tp.String())
}

func TestParse_aliasInAnyOf(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := pcore.Parse(`type M=Array[Variant[Integer[0, 5], Integer[3, 8], M]]`).(dgo.Type)
	require.Equal(t, `m`, tp.String())
}

func TestParse_aliasInMap(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := pcore.Parse(`type M=Hash[String, Variant[Integer[0, 5], Integer[3, 8], M]]`).(dgo.Type)
	require.Equal(t, `m`, tp.String())
}

func TestParse_range(t *testing.T) {
	require.Equal(t, tf.Integer(vf.Integer(1), vf.Integer(10), true), pcore.Parse(`Integer[1,10]`))
	require.Equal(t, tf.Integer(vf.Integer(1), nil, true), pcore.Parse(`Integer[1]`))
	require.Equal(t, tf.Integer(nil, vf.Integer(0), true), pcore.Parse(`Integer[default,0]`))
	require.Equal(t, tf.Float(vf.Float(1), vf.Float(10), true), pcore.Parse(`Float[1.0,10]`))
	require.Equal(t, tf.Float(vf.Float(1), vf.Float(10), true), pcore.Parse(`Float[1,10.0]`))
	require.Equal(t, tf.Float(vf.Float(1), vf.Float(10), true), pcore.Parse(`Float[1.0,10.0]`))
	require.Equal(t, tf.Float(vf.Float(1), nil, true), pcore.Parse(`Float[1.0]`))
	require.Equal(t, tf.Float(nil, vf.Float(0), true), pcore.Parse(`Float[default,0.0]`))

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

func TestParse_map(t *testing.T) {
	require.Equal(t, vf.Map(`a`, 10), pcore.Parse(`{a => 10}`))
	require.Equal(t, vf.Values(1, vf.Map(`a`, 10), 23), pcore.Parse(`[1, a => 10, 23]`))
	require.Panic(t, func() { pcore.Parse(`{]`) }, `expected a literal, got '\]'`)
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
	internal.ResetDefaultAliases()
	tp := pcore.Parse(`
Struct[
  types => [
    type Ascii=Integer[1,127],
    type Slug=Pattern[/^[a-z0-9-]+$/]
  ],
  x => Hash[Slug,Struct["Token" => Ascii, value => String]]
]`).(dgo.StructMapType)
	require.Equal(t, `map[slug]{"Token":ascii,"value":string}`, tp.GetEntryType(`x`).Value().(dgo.Type).String())
}

func TestParse_errors(t *testing.T) {
	require.Panic(t, func() { pcore.Parse(`Integer[1 23]`) }, `expected one of ',' or '\]', got 23: \(column: 11\)`)
	require.Panic(t, func() { pcore.Parse(`Float[1. 3]`) }, `unexpected character ' '`)
	require.Panic(t, func() { pcore.Parse(`String[1 . 3]`) }, `expected one of ',' or '\]', got '.'`)
	require.Panic(t, func() { pcore.Parse(`/\`) }, `unterminated regexp`)
	require.Panic(t, func() { pcore.Parse(`/\/`) }, `unterminated regexp`)
	require.Panic(t, func() { pcore.Parse(`Pattern[/\//`) }, `expected one of ',' or '\]', got EOT`)
	require.Panic(t, func() { pcore.Parse(`Pattern[/\t/`) }, `expected one of ',' or '\]', got EOT`)
	require.Panic(t, func() { pcore.Parse(`{1, 2}`) }, `expected '=>', got ','`)
	require.Panic(t, func() { pcore.Parse(`{1 => 2 3}`) }, `expected one of ',' or '\}', got 3`)
	require.Panic(t, func() { pcore.Parse(`?`) }, `expected a literal, got '\?'`)
	require.Panic(t, func() { pcore.Parse(`type Foo {}`) }, `expected '=', got Foo`)
	require.Panic(t, func() { pcore.Parse(`Foo(a b)`) }, `expected one of ',' or '\)', got b`)
	require.Panic(t, func() { pcore.Parse(`+x`) }, `unexpected character '\+'`)
	require.Panic(t, func() { pcore.Parse(`A:B`) }, `unexpected character 'B'`)
	require.Panic(t, func() { pcore.Parse(`a:b`) }, `unexpected character 'b'`)
}

func TestParseFile_errors(t *testing.T) {
	require.Panic(t,
		func() { pcore.ParseFile(nil, `foo.dgo`, `[1 2]`) },
		`expected one of ',' or '\]', got 2: \(file: foo\.dgo, line: 1, column: 4\)`)
}
