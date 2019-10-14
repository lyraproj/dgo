package internal_test

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/util"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func ExampleMap_SetType_structType() {
	m := vf.Map(`host`, `example.com`, `port`, 22).Copy(false)
	m.SetType(`{host: string, port: int}`)
	fmt.Println(m)
	// Output: {"host":"example.com","port":22}
}

func ExampleMap_Put_structType() {
	m := vf.Map(`host`, `example.com`, `port`, 22).Copy(false)
	m.SetType(`{host: string, port: int}`)
	m.Put(`port`, 465)
	fmt.Println(m)
	// Output: {"host":"example.com","port":465}
}

func ExampleMap_Put_structTypeAdditionalKey() {
	m := vf.Map(`host`, `example.com`, `port`, 22).Copy(false)
	m.SetType(`{host: string, port: int, ...}`)
	m.Put(`login`, `bob`)
	fmt.Println(m)
	// Output: {"host":"example.com","port":22,"login":"bob"}
}

func ExampleMap_Put_structTypeIllegalKey() {
	m := vf.Map(`host`, `example.com`, `port`, 22).Copy(false)
	m.SetType(`{host: string, port: int}`)
	if err := util.Catch(func() { m.Put(`login`, `bob`) }); err != nil {
		fmt.Println(err)
	}
	// Output: key "login" cannot added to type {"host":string,"port":int}
}

func ExampleMap_Put_structTypeIllegalValue() {
	m := vf.Map(`host`, `example.com`, `port`, 22).Copy(false)
	m.SetType(`{host: string, port: int}`)
	if err := util.Catch(func() { m.Put(`port`, `22`) }); err != nil {
		fmt.Println(err)
	}
	// Output: the string "22" cannot be assigned to a variable of type int
}

func TestStructType_Get(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	require.Equal(t, tp.Get(`a`).Value(), typ.Integer)
	require.Nil(t, tp.Get(`c`))
}

func TestStructType_Validate(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	es := tp.Validate(nil, vf.Map(`a`, 1, `b`, `yes`))
	require.Equal(t, 0, len(es))
}

func TestStructType_Validate_valueType(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	es := tp.Validate(nil, vf.Map(`a`, `no`, `b`, `yes`))
	require.Equal(t, 1, len(es))
	require.Equal(t, `parameter 'a' is not an instance of type int`, es[0].Error())
}

func TestStructType_Validate_missingKey(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	es := tp.Validate(nil, vf.Map(`a`, 1))
	require.Equal(t, 1, len(es))
	require.Equal(t, `missing required parameter 'b'`, es[0].Error())
}

func TestStructType_Validate_unknownKey(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	es := tp.Validate(nil, vf.Map(`a`, 1, `b`, `yes`, `c`, `no`))
	require.Equal(t, 1, len(es))
	require.Equal(t, `unknown parameter 'c'`, es[0].Error())
}

func TestStructType_Validate_notMap(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	es := tp.Validate(nil, vf.Values(1, 2))
	require.Equal(t, 1, len(es))
	require.NotOk(t, `value is not a Map`, es[0])
}

func TestStructType_ValidateVerbose_valueType(t *testing.T) {
	tp := newtype.Parse(`{a:int}`).(dgo.StructMapType)
	out := util.NewIndenter(`  `)
	ok := tp.ValidateVerbose(vf.Map(`a`, `no`), out)
	es := out.String()
	require.False(t, ok)
	require.Equal(t, `Validating 'a' against definition int
  'a' FAILED!
  Reason: expected a value of type int, got "no"
`, es)
}

func TestStructType_ValidateVerbose_missingKey(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	out := util.NewIndenter(`  `)
	ok := tp.ValidateVerbose(vf.Map(`a`, 1), out)
	es := out.String()
	require.False(t, ok)
	require.Equal(t, `Validating 'a' against definition int
  'a' OK!
Validating 'b' against definition string
  'b' FAILED!
  Reason: required key not found in input
`, es)
}

func TestStructType_ValidateVerbose_unknownKey(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	out := util.NewIndenter(`  `)
	ok := tp.ValidateVerbose(vf.Map(`a`, 1, `b`, `yes`, `c`, `no`), out)
	es := out.String()
	require.False(t, ok)
	require.Equal(t, `Validating 'a' against definition int
  'a' OK!
Validating 'b' against definition string
  'b' OK!
Validating 'c'
  'c' FAILED!
  Reason: key is not found in definition
`, es)
}

func TestStructType_ValidateVerbose(t *testing.T) {
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructMapType)
	out := util.NewIndenter(``)
	require.False(t, tp.ValidateVerbose(vf.Values(1, 2), out))
	require.Equal(t, `value is not a Map`, out.String())
}

func TestStructType_alias(t *testing.T) {
	tp := newtype.Parse(`person={name:string,mom:person,dad:person}`).(dgo.StructMapType)
	require.Same(t, tp, tp.Get(`mom`).Value())
}

func TestStructType(t *testing.T) {
	tp := newtype.StructMap(false)
	require.Equal(t, tp, tp)
	require.Equal(t, tp, newtype.Parse(`{}`))
	require.Equal(t, tp.KeyType(), typ.Any)
	require.Equal(t, tp.ValueType(), typ.Any)
	require.True(t, tp.Unbounded())

	tp = newtype.StructMap(false,
		newtype.StructMapEntry(`a`, typ.Integer, true))
	require.Equal(t, tp, newtype.Parse(`{a:int}`))
	require.Assignable(t, tp, newtype.Parse(`{a:0..5}`))
	require.Assignable(t, tp, newtype.AnyOf(newtype.Parse(`{a:0..5}`), newtype.Parse(`{a:10..15}`)))
	require.NotAssignable(t, tp, newtype.Parse(`{a:float}`))
	require.Equal(t, tp.KeyType(), vf.String(`a`).Type())
	require.Equal(t, tp.ValueType(), typ.Integer)
	require.Equal(t, tp.Min(), 1)
	require.Equal(t, tp.Max(), 1)

	tp = newtype.StructMap(false,
		newtype.StructMapEntry(`a`, typ.Integer, true),
		newtype.StructMapEntry(`b`, typ.String, false))

	require.Equal(t, tp, tp)
	require.Equal(t, tp, newtype.Parse(`{a:int,b?:string}`))
	require.Equal(t, tp.KeyType(), newtype.Parse(`"a"&"b"`))
	require.Equal(t, tp.ValueType(), newtype.Parse(`int&string`))
	require.NotEqual(t, tp, newtype.Parse(`map[string](int|string)`))
	require.Equal(t, tp.Min(), 1)
	require.Equal(t, tp.Max(), 2)
	require.False(t, tp.Additional())

	m := vf.Map(`a`, 3, `b`, `yes`)
	require.Assignable(t, tp, tp)
	require.Assignable(t, tp, m.Type())
	require.Instance(t, tp, m)

	m = vf.Map(`a`, 3)
	require.Assignable(t, tp, m.Type())
	require.Instance(t, tp, m)

	m = vf.Map(`b`, `yes`)
	require.NotAssignable(t, tp, m.Type())
	require.NotInstance(t, tp, m)

	m = vf.Map(`a`, 3, `b`, 4)
	require.NotAssignable(t, tp, m.Type())
	require.NotInstance(t, tp, m)

	require.NotInstance(t, tp, vf.Values(`a`, `b`))

	require.Instance(t, tp.Type(), tp)

	tps := newtype.StructMap(false,
		newtype.StructMapEntry(`a`, newtype.IntegerRange(0, 10, true), true),
		newtype.StructMapEntry(`b`, newtype.String(20), false))
	require.Equal(t, tps, newtype.Parse(`{a:0..10,b?:string[20]}`))
	require.Assignable(t, tp, tps)

	tps = newtype.StructMap(false,
		newtype.StructMapEntry(`a`, typ.Integer, true),
		newtype.StructMapEntry(`b`, typ.String, true))
	require.Equal(t, tps, newtype.Parse(`{a:int,b:string}`))
	require.Assignable(t, tp, tps)

	tps = newtype.StructMap(false,
		newtype.StructMapEntry(`a`, typ.Integer, false),
		newtype.StructMapEntry(`b`, typ.String, false))
	require.Equal(t, tps, newtype.Parse(`{a?:int,b?:string}`))
	require.NotAssignable(t, tp, tps)

	tps = newtype.StructMap(false,
		newtype.StructMapEntry(`a`, typ.Integer, true))
	require.Equal(t, tps, newtype.Parse(`{a:int}`))
	require.Assignable(t, tp, tps)

	tps = newtype.StructMap(false,
		newtype.StructMapEntry(`a`, 3, true))
	require.Equal(t, tps, newtype.Parse(`{a:3}`))
	require.Assignable(t, tp, tps)

	tps = newtype.StructMap(false,
		newtype.StructMapEntry(`b`, typ.String, false))
	require.Equal(t, tps, newtype.Parse(`{b?:string}`))
	require.NotAssignable(t, tp, tps)

	tps = newtype.StructMap(true,
		newtype.StructMapEntry(`a`, typ.Integer, true))
	require.Equal(t, tps, newtype.Parse(`{a:int,...}`))
	require.NotAssignable(t, tp, tps)
	require.Equal(t, tps.Min(), 1)
	require.Equal(t, tps.Max(), math.MaxInt64)
	require.True(t, tps.Additional())

	require.NotEqual(t, 0, tp.HashCode())
	require.NotEqual(t, tp.HashCode(), newtype.Parse(`{a:int,b?:string,...}`).HashCode())

	require.Panic(t, func() {
		newtype.StructMap(false,
			newtype.StructMapEntry(newtype.Pattern(regexp.MustCompile(`a*`)), typ.Integer, true))
	}, `non exact key types`)

	tps = newtype.Parse(`{a:0..10,b?:int}`).(dgo.StructMapType)
	require.True(t, reflect.ValueOf(map[string]int64{}).Type().AssignableTo(tps.ReflectType()))
}

func TestStructEntry(t *testing.T) {
	tp := newtype.StructMapEntry(`a`, typ.String, true)
	require.Equal(t, tp, newtype.StructMapEntry(`a`, typ.String, true))
	require.NotEqual(t, tp, newtype.StructMapEntry(`a`, typ.String, false))
	require.NotEqual(t, tp, vf.Values(`a`, typ.String))
	require.Equal(t, `"a":string`, tp.String())
}

func TestStructFromMap(t *testing.T) {
	require.Panic(t, func() { newtype.StructMapFromMap(false, vf.Map(`nope`, `dope`)) }, `cannot be assigned to a variable of type map`)

	tp := newtype.StructMapFromMap(false, vf.Map(`first`, typ.String))
	require.Equal(t, tp, newtype.StructMap(false, newtype.StructMapEntry(`first`, typ.String, false)))

	tp = newtype.StructMapFromMap(false, vf.Map(`first`, vf.Map(`type`, typ.String)))
	require.Equal(t, tp, newtype.StructMap(false, newtype.StructMapEntry(`first`, typ.String, false)))

	tp = newtype.StructMapFromMap(false, vf.Map(`first`, vf.Map(`type`, typ.String, `required`, true)))
	require.Equal(t, tp, newtype.StructMap(false, newtype.StructMapEntry(`first`, typ.String, true)))

	tp = newtype.StructMapFromMap(false, vf.Map(`first`, vf.Map(`type`, typ.String, `required`, false)))
	require.Equal(t, tp, newtype.StructMap(false, newtype.StructMapEntry(`first`, typ.String, false)))

	tp = newtype.StructMapFromMap(false, vf.Map(`first`, vf.Map(`type`, `string`, `required`, false)))
	require.Equal(t, tp, newtype.StructMap(false, newtype.StructMapEntry(`first`, typ.String, false)))
}
