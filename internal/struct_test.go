package internal_test

import (
	"fmt"
	"math"
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
	tp := newtype.Parse(`{a:int,b:string}`).(dgo.StructType)
	a, ok := tp.Get(`a`)
	require.True(t, ok)
	require.Equal(t, a.Value(), typ.Integer)

	_, ok = tp.Get(`c`)
	require.False(t, ok)
}

func TestStructType_alias(t *testing.T) {
	tp := newtype.Parse(`person={name:string,mom:person,dad:person}`).(dgo.StructType)
	mom, ok := tp.Get(`mom`)
	require.True(t, ok)
	require.Same(t, tp, mom.Value())
}

func TestStructType(t *testing.T) {
	tp := newtype.Struct(false)
	require.Equal(t, tp, tp)
	require.Equal(t, tp, newtype.Parse(`{}`))
	require.Equal(t, tp.KeyType(), typ.Any)
	require.Equal(t, tp.ValueType(), typ.Any)
	require.True(t, tp.Unbounded())

	tp = newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, true))
	require.Equal(t, tp, newtype.Parse(`{a:int}`))
	require.Assignable(t, tp, newtype.Parse(`{a:0..5}`))
	require.Assignable(t, tp, newtype.AnyOf(newtype.Parse(`{a:0..5}`), newtype.Parse(`{a:10..15}`)))
	require.NotAssignable(t, tp, newtype.Parse(`{a:float}`))
	require.Equal(t, tp.KeyType(), vf.String(`a`).Type())
	require.Equal(t, tp.ValueType(), typ.Integer)
	require.Equal(t, tp.Min(), 1)
	require.Equal(t, tp.Max(), 1)

	tp = newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, true),
		newtype.StructEntry(`b`, typ.String, false))

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

	tps := newtype.Struct(false,
		newtype.StructEntry(`a`, newtype.IntegerRange(0, 10), true),
		newtype.StructEntry(`b`, newtype.String(20), false))
	require.Equal(t, tps, newtype.Parse(`{a:0..10,b?:string[20]}`))
	require.Assignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, true),
		newtype.StructEntry(`b`, typ.String, true))
	require.Equal(t, tps, newtype.Parse(`{a:int,b:string}`))
	require.Assignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, false),
		newtype.StructEntry(`b`, typ.String, false))
	require.Equal(t, tps, newtype.Parse(`{a?:int,b?:string}`))
	require.NotAssignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, true))
	require.Equal(t, tps, newtype.Parse(`{a:int}`))
	require.Assignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`a`, 3, true))
	require.Equal(t, tps, newtype.Parse(`{a:3}`))
	require.Assignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`b`, typ.String, false))
	require.Equal(t, tps, newtype.Parse(`{b?:string}`))
	require.NotAssignable(t, tp, tps)

	tps = newtype.Struct(true,
		newtype.StructEntry(`a`, typ.Integer, true))
	require.Equal(t, tps, newtype.Parse(`{a:int,...}`))
	require.NotAssignable(t, tp, tps)
	require.Equal(t, tps.Min(), 1)
	require.Equal(t, tps.Max(), math.MaxInt64)
	require.True(t, tps.Additional())

	require.NotEqual(t, 0, tp.HashCode())
	require.NotEqual(t, tp.HashCode(), newtype.Parse(`{a:int,b?:string,...}`).HashCode())

	require.Panic(t, func() {
		newtype.Struct(false,
			newtype.StructEntry(newtype.Pattern(regexp.MustCompile(`a*`)), typ.Integer, true))
	}, `non exact key types`)
}

func TestStructEntry(t *testing.T) {
	tp := newtype.StructEntry(`a`, typ.String, true)
	require.Equal(t, tp, newtype.StructEntry(`a`, typ.String, true))
	require.NotEqual(t, tp, newtype.StructEntry(`a`, typ.String, false))
	require.NotEqual(t, tp, vf.Values(`a`, typ.String))
	require.Equal(t, `"a":string`, tp.String())
}
