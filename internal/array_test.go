package internal_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/internal"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestArray_max_min(t *testing.T) {
	tp := tf.Array(2, 1)
	require.Equal(t, tp.Min(), 1)
	require.Equal(t, tp.Max(), 2)
}

func TestArray_negative_min(t *testing.T) {
	tp := tf.Array(-2, 3)
	require.Equal(t, tp.Min(), 0)
	require.Equal(t, tp.Max(), 3)
}

func TestArray_negative_min_max(t *testing.T) {
	tp := tf.Array(-2, -3)
	require.Equal(t, tp.Min(), 0)
	require.Equal(t, tp.Max(), 0)
}

func TestArray_explicit_unbounded(t *testing.T) {
	tp := tf.Array(0, math.MaxInt64)
	require.Equal(t, tp, typ.Array)
	require.True(t, tp.Unbounded())
}

func TestArray_badOneArg(t *testing.T) {
	require.Panic(t, func() { tf.Array(`bad`) }, `illegal argument`)
}

func TestArray_badTwoArg(t *testing.T) {
	require.Panic(t, func() { tf.Array(`bad`, 2) }, `illegal argument 1`)
	require.Panic(t, func() { tf.Array(typ.String, `bad`) }, `illegal argument 2`)
}

func TestArray_badThreeArg(t *testing.T) {
	require.Panic(t, func() { tf.Array(`bad`, 2, 2) }, `illegal argument 1`)
	require.Panic(t, func() { tf.Array(typ.String, `bad`, 2) }, `illegal argument 2`)
	require.Panic(t, func() { tf.Array(typ.String, 2, `bad`) }, `illegal argument 3`)
}

func TestArray_badArgCount(t *testing.T) {
	require.Panic(t, func() { tf.Array(typ.String, 2, 2, true) }, `illegal number of arguments`)
}

func TestArrayType(t *testing.T) {
	tp := tf.Array()
	v := vf.Strings(`a`, `b`)
	require.Instance(t, tp, v)
	require.Assignable(t, tp, tf.AnyOf(tf.Array(tf.String(5, 5)), tf.Array(tf.String(8, 8))))
	require.Equal(t, typ.Any, tp.ElementType())
	require.Equal(t, 0, tp.Min())
	require.Equal(t, math.MaxInt64, tp.Max())
	require.True(t, tp.Unbounded())

	require.Instance(t, tp.Type(), tp)
	require.Equal(t, `[]any`, tp.String())

	require.NotEqual(t, 0, tp.HashCode())
	require.Equal(t, tp.HashCode(), tp.HashCode())
}

func TestArrayType_New(t *testing.T) {
	at := tf.Array(typ.Integer)
	a := vf.Values(1, 2)
	require.Same(t, a, vf.New(typ.Array, a))
	require.Same(t, a, vf.New(at, a))
	require.Same(t, a, vf.New(a.Type(), a))
	require.Same(t, a, vf.New(tf.Tuple(typ.Integer, typ.Integer), a))
	require.Same(t, a, vf.New(at, vf.Arguments(a)))
	require.Equal(t, a, vf.New(at, vf.Values(1, 2)))

	require.Panic(t, func() { vf.New(at, vf.Values(`first`, `second`)) }, `cannot be assigned`)
}

func TestSizedArrayType(t *testing.T) {
	tp := tf.Array(typ.String)
	v := vf.Strings(`a`, `b`)
	require.Instance(t, tp, v)
	require.NotInstance(t, tp, `a`)
	require.NotAssignable(t, tp, typ.Array)
	require.Assignable(t, tp, tf.AnyOf(tf.Array(tf.String(5, 5)), tf.Array(tf.String(8, 8))))
	require.Equal(t, tp, tp)
	require.NotEqual(t, tp, typ.Array)
	require.NotEqual(t, tp, tf.Map(typ.String, typ.String))

	require.Instance(t, tp.Type(), tp)
	require.Equal(t, `[]string`, tp.String())

	tp = tf.Array(typ.Integer)
	v = vf.Strings(`a`, `b`)
	require.NotInstance(t, tp, v)
	require.NotAssignable(t, tp, v.Type())

	tp = tf.Array(tf.AnyOf(typ.String, typ.Integer))
	v = vf.Values(`a`, 3)
	require.Instance(t, tp, v)
	require.Assignable(t, tp, v.Type())
	require.Instance(t, v.Type(), v)
	v = v.With(vf.True)
	require.NotInstance(t, tp, v)
	require.NotAssignable(t, tp, v.Type())

	tp = tf.Array(0, 2)
	v = vf.Strings(`a`, `b`)
	require.Instance(t, tp, v)

	tp = tf.Array(0, 1)
	require.NotInstance(t, tp, v)

	tp = tf.Array(2, 3)
	require.Instance(t, tp, v)

	tp = tf.Array(3, 3)
	require.NotInstance(t, tp, v)

	tp = tf.Array(typ.String, 2, 3)
	require.Instance(t, tp, v)

	tp = tf.Array(typ.Integer, 2, 3)
	require.NotInstance(t, tp, v)

	require.NotEqual(t, 0, tp.HashCode())
	require.NotEqual(t, tp.HashCode(), tf.Array(typ.Integer).HashCode())
	require.Equal(t, `[2,3]int`, tp.String())

	tp = tf.Array(tf.IntegerRange(0, 15, true), 2, 3)
	require.Equal(t, `[2,3]0..15`, tp.String())

	tp = tf.Array(tf.Array(2, 2), 0, 10)
	require.Equal(t, `[0,10][2,2]any`, tp.String())

	require.Equal(t, tf.Array(typ.Any).ReflectType(), typ.Array.ReflectType())
}

func TestExactArrayType(t *testing.T) {
	v := vf.Strings()
	tp := v.Type().(dgo.TupleType)
	require.Equal(t, typ.Any, tp.ElementType())

	v = vf.Strings(`a`)
	tp = v.Type().(dgo.TupleType)
	et := vf.String(`a`).Type()
	require.Equal(t, et, tp.ElementType())
	require.Assignable(t, tp, tf.Array(et, 1, 1))
	require.NotAssignable(t, tp, tf.Array(typ.String, 1, 1))

	v = vf.Strings(`a`, `b`)
	tp = v.Type().(dgo.TupleType)
	require.Instance(t, tp, v)
	require.Equal(t, tp, vf.Strings(`a`, `b`).Type())
	require.NotInstance(t, tp, `a`)
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.Array)
	require.NotAssignable(t, tp, tf.Array(et, 1, 1))

	require.Assignable(t, tp, tf.Tuple(vf.String(`a`).Type(), vf.String(`b`).Type()))
	require.NotAssignable(t, tp, tf.Tuple(vf.String(`a`).Type(), vf.String(`b`).Type(), vf.String(`c`).Type()))
	require.NotAssignable(t, tp, tf.Tuple(vf.String(`a`).Type(), typ.String))

	require.Equal(t, 2, tp.Min())
	require.Equal(t, 2, tp.Max())
	require.False(t, tp.Unbounded())

	require.Equal(t, tf.Array(typ.String), typ.Generic(tp))

	require.NotAssignable(t, tp, tf.AnyOf(tf.Array(tf.String(5, 5)), tf.Array(tf.String(8, 8))))
	require.Equal(t, tp, tp)
	require.Equal(t, vf.Values(vf.String(`a`).Type(), vf.String(`b`).Type()), tp.ElementTypes())
	require.NotEqual(t, tp, typ.Array)

	require.Instance(t, tp.Type(), tp)
	require.Equal(t, `{"a","b"}`, tp.String())

	require.NotEqual(t, 0, tp.HashCode())
	require.NotEqual(t, tp.HashCode(), tf.Array(typ.Integer).HashCode())
}

func TestArrayElementType_singleElement(t *testing.T) {
	v := vf.Strings(`hello`)
	at := v.Type()
	et := at.(dgo.ArrayType).ElementType()
	require.Assignable(t, at, at)
	tp := tf.Array(typ.String)
	require.Assignable(t, tp, at)
	require.NotAssignable(t, at, tp)
	require.Assignable(t, et, vf.String(`hello`).Type())
	require.NotAssignable(t, et, vf.String(`hey`).Type())
	require.Instance(t, et, `hello`)
	require.NotInstance(t, et, `world`)
	require.Equal(t, et, vf.Strings(`hello`).Type().(dgo.ArrayType).ElementType())
	//	require.Equal(t, et.(dgo.ExactType).Value(), vf.Strings(`hello`))
	require.NotEqual(t, et, vf.Strings(`hello`).Type().(dgo.ArrayType))

	require.NotEqual(t, 0, et.HashCode())
	require.Equal(t, et.HashCode(), et.HashCode())
}

func TestArrayElementType_multipleElements(t *testing.T) {
	v := vf.Strings(`hello`, `world`)
	at := v.Type()
	et := at.(dgo.ArrayType).ElementType()
	require.Assignable(t, at, at)
	tp := tf.Array(typ.String)
	require.Assignable(t, tp, at)
	require.NotAssignable(t, at, tp)
	require.NotAssignable(t, et, vf.String(`hello`).Type())
	require.NotAssignable(t, et, vf.String(`world`).Type())
	require.NotAssignable(t, vf.String(`hello`).Type(), et)
	require.NotAssignable(t, vf.String(`world`).Type(), et)
	require.Assignable(t, vf.Strings(`hello`, `world`).Type(), at)
	require.Assignable(t, vf.Strings(`hello`, `world`).Type().(dgo.ArrayType).ElementType(), et)
	require.Assignable(t, vf.Strings(`world`, `hello`).Type().(dgo.ArrayType).ElementType(), et)

	require.Assignable(t, tf.Array(2, 2), at)
	require.Assignable(t, tf.Array(et, 2, 2), at)

	et.(dgo.TernaryType).Operands().Each(func(v dgo.Value) {
		t.Helper()
		require.Instance(t, typ.String, v)
	})

	require.Instance(t, et.Type(), et)
	require.Equal(t, `"hello"&"world"`, et.String())
}

func TestTupleType(t *testing.T) {
	tt := typ.Tuple
	require.Same(t, tt, tf.VariadicTuple(typ.Any))

	tt = typ.EmptyTuple
	require.Same(t, tt, tf.Tuple())
	require.Assignable(t, tt, tf.Array(0, 0))

	tt = tf.Tuple(typ.String, typ.Any, typ.Float)
	require.Assignable(t, tt, tf.Tuple(typ.String, typ.Integer, typ.Float))
	require.Assignable(t, tt, vf.Values(`one`, 2, 3.0).Type())
	require.Equal(t, 3, tt.Min())
	require.Equal(t, 3, tt.Max())
	require.False(t, tt.Unbounded())

	et := tf.String(1, 100)
	tt = tf.Tuple(et)
	require.Same(t, tt.ElementType(), et)

	tt = tf.Tuple(typ.String, typ.Integer, typ.Float)
	require.Assignable(t, tt, tt)
	require.Assignable(t, tt, tf.Tuple(typ.String, typ.Integer, typ.Float))
	require.NotAssignable(t, tt, tf.Tuple(typ.String, typ.Integer, typ.Boolean))
	require.Assignable(t, typ.Array, tt)
	require.Assignable(t, tf.Array(0, 3), tt)

	require.Assignable(t, tt, vf.Values(`one`, 2, 3.0).Type())
	require.NotAssignable(t, tt, vf.Values(`one`, 2, 3.0, `four`).Type())
	require.NotAssignable(t, tt, vf.Values(`one`, 2, 3).Type())
	require.NotAssignable(t, tt, typ.Array)
	require.NotAssignable(t, tt, typ.Tuple)
	require.NotEqual(t, tt, typ.String)

	require.Equal(t, 3, tt.Min())
	require.Equal(t, 3, tt.Max())
	require.False(t, tt.Unbounded())

	require.Assignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer, typ.Float)), tt)
	require.NotAssignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer)), tt)
	require.Assignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer, typ.Float, typ.Boolean)), tt)

	require.Assignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer, typ.Float), 0, 3), tt)
	require.NotAssignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer, typ.Float), 0, 2), tt)

	okv := vf.Values(`hello`, 1, 2.0)
	require.Instance(t, tt, okv)
	require.NotInstance(t, tt, okv.Get(0))
	require.Assignable(t, tt, okv.Type())
	require.Assignable(t, typ.Array, tt)

	okv = vf.Values(`hello`, 1, 2)
	require.NotInstance(t, tt, okv)

	okv = vf.Values(`hello`, 1, 2.0, true)
	require.NotInstance(t, tt, okv)

	okm := vf.MutableValues(`world`, 2, 3.0)
	okm.SetType(tt)
	require.Panic(t, func() { okm.Add(3) }, tf.IllegalSize(tt, 4))
	require.Panic(t, func() { okm.Set(2, 3) }, tf.IllegalAssignment(typ.Float, vf.Value(3)))
	require.Equal(t, `world`, okm.Set(0, `earth`))

	tt = tf.Tuple(typ.String, typ.String)
	require.Assignable(t, tt, tf.Array(typ.String, 2, 2))
	require.Equal(t, tt.ReflectType(), tf.Array(typ.String).ReflectType())

	require.Assignable(t, tt, tf.Array(typ.String, 2, 2))
	require.NotAssignable(t, tt, tf.AnyOf(typ.Nil, tf.Array(typ.String, 2, 2)))
	tt = tf.Tuple(typ.String, typ.Integer)
	require.NotAssignable(t, tt, tf.Array(typ.String, 2, 2))
	require.Equal(t, tt.ReflectType(), typ.Array.ReflectType())

	require.Equal(t, typ.Any, typ.Tuple.ElementType())
	require.Equal(t, tf.AllOf(typ.String, typ.Integer), tt.ElementType())
	require.Equal(t, vf.Values(typ.String, typ.Integer), tt.ElementTypes())
	require.Equal(t, typ.Array, typ.Generic(tt))

	require.Instance(t, tt.Type(), tt)
	require.Equal(t, `{string,int}`, tt.String())

	require.NotEqual(t, 0, tt.HashCode())
	require.Equal(t, tt.HashCode(), tt.HashCode())

	tt = tf.Tuple(typ.String, typ.String)
	require.Equal(t, tf.Array(typ.String), typ.Generic(tt))
	require.Equal(t, vf.Values(typ.String, typ.String), typ.ExactValue(tt))

	te := tf.Tuple(vf.String(`a`).Type(), vf.String(`b`).Type())
	require.Assignable(t, tt, te)
	require.NotAssignable(t, te, tt)
	require.Equal(t, tf.Array(typ.String), typ.Generic(te))
	require.Equal(t, vf.Strings(`a`, `b`), typ.ExactValue(te))
}

func TestTupleType_selfReference(t *testing.T) {
	tp := tf.ParseType(`x={string,x}`).(dgo.ArrayType)
	d := vf.MutableValues()
	d.Add(`hello`)
	d.Add(d)
	require.Instance(t, tp, d)

	t2 := tf.ParseType(`x={string,{string,x}}`)
	require.Assignable(t, tp, t2)

	require.Equal(t, `{string,<recursive self reference to tuple type>}`, tp.String())
}

func TestVariadicTupleType(t *testing.T) {
	tt := tf.ParseType(`{string,...string}`).(dgo.TupleType)
	require.Instance(t, tt, vf.Values(`one`))
	require.Instance(t, tt, vf.Values(`one`, `two`))
	require.NotInstance(t, tt, vf.Values(`one`, 2))
	require.NotInstance(t, tt, vf.Values(1))

	tt = tf.VariadicTuple(typ.String, typ.String)
	require.NotInstance(t, tt, vf.Values())
	require.Instance(t, tt, vf.Values(`one`))
	require.Instance(t, tt, vf.Values(`one`, `two`))
	require.Instance(t, tt, vf.Values(`one`, `two`, `three`))
	require.True(t, tt.Unbounded())
	require.Equal(t, 1, tt.Min())
	require.Equal(t, math.MaxInt64, tt.Max())

	a := vf.MutableValues(`one`, `two`)
	a.SetType(tt)
	require.Panic(t, func() { a.Set(0, 1) }, `cannot be assigned`)
	require.Panic(t, func() { a.Set(1, 2) }, `cannot be assigned`)

	require.Panic(t, func() { tf.VariadicTuple() }, `must have at least one element`)
}

func TestMutableValues_withoutType(t *testing.T) {
	a := vf.MutableValues(nil)
	require.True(t, vf.Nil == a.Get(0))
}

func TestMutableValues_maxSizeMismatch(t *testing.T) {
	a := vf.MutableValues(true)
	a.SetType(tf.Array(0, 1))
	require.Panic(t, func() { a.Add(vf.False) }, `size constraint`)
	require.Panic(t, func() { a.With(vf.False) }, `size constraint`)
	require.Panic(t, func() { a.WithAll(vf.Values(false)) }, `size constraint`)
	require.Panic(t, func() { a.WithValues(false) }, `size constraint`)

	a.WithAll(vf.Values()) // No panic
	a.WithValues()         // No panic
}

func TestMutableValues_minSizeMismatch(t *testing.T) {
	require.Panic(t, func() { vf.MutableValues().SetType(tf.Array(1, 1)) }, `cannot be assigned`)

	a := vf.MutableValues(true)
	a.SetType(`[1,1]bool`)
	require.Panic(t, func() { a.Remove(0) }, `size constraint`)
	require.Panic(t, func() { a.RemoveValue(vf.True) }, `size constraint`)
	require.Panic(t, func() { a.Add(vf.True) }, `size constraint`)
	require.Panic(t, func() { a.AddAll(vf.Values(vf.True)) }, `size constraint`)
}

func TestMutableValues_elementTypeMismatch(t *testing.T) {
	a := vf.MutableValues()
	a.SetType(`[]string`)
	a.Add(`hello`)
	a.AddAll(vf.Values(`hello`))
	a.AddAll(vf.Values())
	require.Panic(t, func() { a.Add(vf.True) }, `cannot be assigned`)
	require.Panic(t, func() { a.AddAll(vf.Values(vf.True)) }, `cannot be assigned`)

	a = vf.MutableValues(`a`, 2)
	a.SetType(tf.Tuple(typ.String, typ.Integer))
	a.Set(0, `hello`)
	a.Set(1, 3)
	require.Panic(t, func() { a.Set(0, 3) }, `cannot be assigned`)
}

func TestMutableValues_tupleTypeMismatch(t *testing.T) {
	require.Panic(t, func() { vf.MutableValues(true).SetType(tf.Tuple(typ.String)) }, `cannot be assigned`)
}

func TestMutableArray(t *testing.T) {
	s := []dgo.Value{nil}
	a := vf.WrapSlice(s)
	require.True(t, vf.Nil == s[0])
	a.Set(0, `hello`)
	require.Equal(t, `hello`, s[0])
}

func TestMutableArray_stringType(t *testing.T) {
	a := vf.WrapSlice(nil)
	a.SetType(`[]int`)
	a.Add(3)
	require.Equal(t, 3, a.Get(0))
}

func TestMutableArray_dgoStringType(t *testing.T) {
	a := vf.WrapSlice(nil)
	a.SetType(vf.String(`[]int`))
	a.Add(3)
	require.Equal(t, 3, a.Get(0))
}

func TestMutableArray_badType(t *testing.T) {
	require.Panic(t, func() { vf.WrapSlice(nil).SetType(`map[string]int`) }, `does not evaluate to an array type`)
}

func TestMutableArray_zeroType(t *testing.T) {
	a := vf.WrapSlice(nil)
	a.SetType([]int{})
	a.Add(3)
	require.Equal(t, 3, a.Get(0))
}

func TestArray(t *testing.T) {
	a := vf.Array([]int{1, 2})
	require.Equal(t, 2, a.Len())
	require.Equal(t, 1, a.Get(0))
	require.Equal(t, 2, a.Get(1))

	require.Same(t, a, vf.Array(a))
	require.Same(t, a, vf.Array(reflect.ValueOf(a)))
}

func TestArray_Set(t *testing.T) {
	a := vf.MutableValues()
	a.SetType(tf.Array(typ.Integer))
	a.Add(1)
	a.Set(0, 2)
	require.Equal(t, 2, a.Get(0))

	require.Panic(t, func() { a.Set(0, 1.0) }, `cannot be assigned`)

	f := a.Copy(true)
	require.Panic(t, func() { f.Set(0, 1) }, `Set .* frozen`)
}

func TestArray_SetType(t *testing.T) {
	a := vf.MutableValues(1, 2.0, `three`)
	adt := a.Type()

	at := tf.Array(tf.AnyOf(typ.Integer, typ.Float, typ.String))
	a.SetType(at)
	require.Same(t, at, a.Type())

	require.Panic(t, func() { a.SetType(`[](float|string)`) },
		`cannot be assigned`)

	require.Panic(t, func() { a.SetType(vf.String(`[](float|string)`)) },
		`cannot be assigned`)

	require.Panic(t, func() { a.SetType(`(float|string)`) },
		`does not evaluate to an array type`)

	a.SetType(nil)
	require.Equal(t, adt, a.Type())

	a.Freeze()
	require.Panic(t, func() { a.SetType(at) }, `SetType .* frozen`)
}

func TestArray_selfReference(t *testing.T) {
	tp := tf.ParseType(`x=[](string|x)`).(dgo.ArrayType)
	d := vf.MutableValues(tp, `hello`)
	d.Add(d)
	require.Instance(t, tp, d)

	t2 := tf.ParseType(`x=[](string|[](string|x))`)
	require.Assignable(t, tp, t2)
}

func TestArray_recursiveFreeze(t *testing.T) {
	a := vf.Array([]dgo.Value{vf.MutableValues(`b`)})
	require.True(t, a.Get(0).(dgo.Array).Frozen())
}

func TestArray_recursiveReflectiveFreeze(t *testing.T) {
	a := vf.Value(
		reflect.ValueOf([]interface{}{
			reflect.ValueOf(vf.MutableValues(`b`))})).(dgo.Array)
	require.True(t, a.Get(0).(dgo.Array).Frozen())
}

func TestArray_replaceNil(t *testing.T) {
	a := vf.Array([]dgo.Value{nil})
	require.True(t, vf.Nil == a.Get(0))
}

func TestArray_fromReflected(t *testing.T) {
	a := vf.Value([]interface{}{`a`, 1, nil}).(dgo.Array)
	require.True(t, a.Frozen())
	require.True(t, vf.Nil == a.Get(2))
}

func TestArray_Add(t *testing.T) {
	a := vf.Values(`a`)
	require.Panic(t, func() { a.Add(vf.Value(`b`)) }, `Add .* frozen`)
	m := a.Copy(false)
	m.Add(vf.Value(`b`))
	require.Equal(t, vf.Values(`a`), a)
	require.Equal(t, vf.Values(`a`, `b`), m)
}

func TestArray_AddAll(t *testing.T) {
	a := vf.Values(`a`)
	require.Panic(t, func() { a.AddAll(vf.Values(`b`)) }, `AddAll .* frozen`)
	c := a.Copy(false)
	c.AddAll(vf.Values(`b`))
	require.Equal(t, vf.Values(`a`), a)
	require.Equal(t, vf.Values(`a`, `b`), c)

	m := vf.Map(`c`, `C`, `d`, `D`)
	c.AddAll(m)
	require.Equal(t, vf.Values(`a`, `b`, internal.NewMapEntry(`c`, `C`), internal.NewMapEntry(`d`, `D`)), c)
}

func TestArray_AddValues(t *testing.T) {
	a := vf.Values(`a`)
	require.Panic(t, func() { a.AddValues(`b`) }, `AddValues .* frozen`)
	m := a.Copy(false)
	m.AddValues(`b`)
	require.Equal(t, vf.Values(`a`), a)
	require.Equal(t, vf.Values(`a`, `b`), m)
}

func TestArray_All(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)
	i := 0
	require.True(t, a.All(func(e dgo.Value) bool {
		i++
		return e.Equals(`a`) || e.Equals(`b`) || e.Equals(`c`)
	}))
	require.Equal(t, 3, i)

	i = 0
	require.False(t, a.All(func(e dgo.Value) bool {
		i++
		return e.Equals(`a`)
	}))
	require.Equal(t, 2, i)
}

func TestArray_Any(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)
	i := 0
	require.True(t, a.Any(func(e dgo.Value) bool {
		i++
		return e.Equals(`b`)
	}))
	require.Equal(t, 2, i)

	i = 0
	require.False(t, a.Any(func(e dgo.Value) bool {
		i++
		return e.Equals(`d`)
	}))
	require.Equal(t, 3, i)
}

func TestArray_One(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)
	i := 0
	require.True(t, a.One(func(e dgo.Value) bool {
		i++
		return e.Equals(`b`)
	}))
	require.Equal(t, 3, i)

	a = vf.Strings(`a`, `b`, `c`, `b`)
	i = 0
	require.False(t, a.One(func(e dgo.Value) bool {
		i++
		return e.Equals(`b`)
	}))
	require.Equal(t, 4, i)

	i = 0
	require.False(t, a.One(func(e dgo.Value) bool {
		i++
		return e.Equals(`d`)
	}))
	require.Equal(t, 4, i)
}

func TestArray_CompareTo(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)

	c, ok := a.CompareTo(a)
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = a.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	_, ok = a.CompareTo(vf.String(`a`))
	require.False(t, ok)

	b := vf.Strings(`a`, `b`, `c`)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 0, c)

	b = vf.Strings(`a`, `b`, `c`, `d`)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, -1, c)

	b = vf.Values(`a`, `b`, `d`, `d`)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, -1, c)

	b = vf.Values(`a`, `b`, 3, `d`)
	_, ok = a.CompareTo(b)
	require.False(t, ok)

	b = vf.Strings(`a`, `b`)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 1, c)

	b = vf.Strings(`a`, `b`, `d`)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, -1, c)

	b = vf.Strings(`a`, `b`, `b`)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 1, c)

	a = vf.MutableValues(`a`, `b`)
	a.Add(a)
	b = vf.MutableValues(`a`, `b`)
	b.Add(b)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 0, c)

	b = vf.MutableValues(`a`, `b`)
	b.Add(a)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 0, c)

	a = vf.Values(`a`, 1, nil)
	b = vf.Values(`a`, 1, nil)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 0, c)

	b = vf.Values(`a`, 1, 2)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, -1, c)

	a = vf.Values(`a`, 1, 2)
	b = vf.Values(`a`, 1, nil)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 1, c)

	a = vf.Values(`a`, 1, 2)
	m := vf.Map(`a`, 1)
	_, ok = a.CompareTo(m)
	require.False(t, ok)

	a = vf.Values(`a`, 1, []int{2})
	b = vf.Values(`a`, 1, 2)
	_, ok = a.CompareTo(b)
	require.False(t, ok)

	a = vf.Values(m)
	b = vf.Values(`a`)
	_, ok = a.CompareTo(b)
	require.False(t, ok)
}

func TestArray_Copy(t *testing.T) {
	a := vf.Values(`a`, `b`, vf.MutableValues(`c`))
	require.Same(t, a, a.Copy(true))
	require.True(t, a.Get(2).(dgo.Freezable).Frozen())

	c := a.Copy(false)
	require.False(t, c.Frozen())
	require.NotSame(t, c, c.Copy(false))

	c = c.Copy(true)
	require.True(t, c.Frozen())
	require.Same(t, c, c.Copy(true))
}

func TestArray_Flatten(t *testing.T) {
	a := vf.Values(`a`, `b`, vf.Values(`c`, `d`), vf.Values(`e`, vf.Values(`f`, `g`)))
	b := vf.Values(`a`, `b`, `c`, `d`, `e`, `f`, `g`)
	require.Equal(t, b, a.Flatten())
	require.Same(t, b, b.Flatten())
}

func TestArray_FromReflected(t *testing.T) {
	vs := []dgo.Value{vf.Integer(2), vf.String(`b`)}
	a := internal.ArrayFromReflected(reflect.ValueOf(vs), false).(dgo.Array)
	require.Equal(t, reflect.ValueOf(a.GoSlice()).Pointer(), reflect.ValueOf(vs).Pointer())

	a = internal.ArrayFromReflected(reflect.ValueOf(vs), true).(dgo.Array)
	require.NotEqual(t, reflect.ValueOf(a.GoSlice()).Pointer(), reflect.ValueOf(vs).Pointer())
}

func TestArray_Equal(t *testing.T) {
	a := vf.Values(1, 2)
	require.True(t, a.Equals(a))

	b := vf.Values(1, nil)
	require.False(t, a.Equals(b))

	a = vf.Values(1, nil)
	require.True(t, a.Equals(b))

	b = vf.Values(1, 2)
	require.False(t, a.Equals(b))

	a = vf.Values(1, []int{2})
	require.False(t, a.Equals(b))

	b = vf.Values(1, map[int]int{2: 1})
	require.False(t, a.Equals(b))

	b = vf.Values(1, `2`)
	require.False(t, a.Equals(b))

	// Values containing themselves.
	a = vf.MutableValues(`2`)
	a.Add(a)

	b = vf.MutableValues(`2`)
	b.Add(b)

	m := vf.MutableMap()
	m.Put(`me`, m)
	a.Add(m)
	b.Add(m)
	require.True(t, a.Equals(b))

	require.Equal(t, a.HashCode(), b.HashCode())
}

func TestArray_EachWithIndex(t *testing.T) {
	ni := 0
	vf.Values(1, 2, 3).EachWithIndex(func(v dgo.Value, i int) {
		require.Equal(t, ni, i)
		require.Equal(t, v, i+1)
		ni++
	})
	require.Equal(t, 3, ni)
}

func TestArray_Find(t *testing.T) {
	v := vf.Values(`a`, `b`, 3, `d`).Find(func(v dgo.Value) interface{} {
		if v.Equals(3) {
			return `three`
		}
		return nil
	})
	require.Equal(t, v, `three`)
}

func TestArray_Find_notFound(t *testing.T) {
	v := vf.Values(`a`, `b`, `d`).Find(func(v dgo.Value) interface{} { return nil })
	require.True(t, v == nil)
}

func TestArray_Freeze(t *testing.T) {
	a := vf.MutableValues(`a`, `b`, vf.MutableValues(`c`))
	require.False(t, a.Frozen())

	sa := a.Get(2).(dgo.Array)
	require.False(t, sa.Frozen())

	// In place recursive freeze
	a.Freeze()

	// Sub Array is frozen in place
	require.Same(t, a.Get(2), sa)
	require.True(t, a.Frozen())
	require.True(t, sa.Frozen())
}

func TestArray_FrozenEqual(t *testing.T) {
	f := vf.Values(1, 2, 3)
	require.True(t, f.Frozen(), `not frozen`)

	a := f.Copy(false)
	require.False(t, a.Frozen(), `frozen`)

	require.Equal(t, f, a)
	require.Equal(t, a, f)

	a.Freeze()
	require.True(t, a.Frozen(), `not frozen`)
	require.Same(t, a, a.Copy(true))

	b := a.Copy(false)
	require.NotSame(t, a, b)
	require.NotSame(t, b, b.Copy(true))
	require.NotSame(t, b, b.Copy(false))
}

func TestArray_IndexOf(t *testing.T) {
	a := vf.Values(1, nil, 3)
	require.Equal(t, 2, a.IndexOf(3))
	require.Equal(t, 1, a.IndexOf(nil))
	require.Equal(t, 1, a.IndexOf(vf.Nil))
}

func TestArray_Insert(t *testing.T) {
	a := vf.Values(`a`)
	require.Panic(t, func() { a.Insert(0, vf.Value(`b`)) }, `Insert .* frozen`)
	m := a.Copy(false)
	m.Insert(0, vf.Value(`b`))
	require.Equal(t, vf.Values(`a`), a)
	require.Equal(t, vf.Values(`b`, `a`), m)
}

func TestArray_Map(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)
	require.Equal(t, vf.Strings(`d`, `e`, `f`), a.Map(func(e dgo.Value) interface{} {
		return string([]byte{e.String()[0] + 3})
	}))
	require.Equal(t, vf.Values(vf.Nil, vf.Nil, vf.Nil), a.Map(func(e dgo.Value) interface{} {
		return nil
	}))
}

func TestArray_MapTo(t *testing.T) {
	require.Equal(t, vf.Integers(97), vf.Strings(`a`).MapTo(nil, func(e dgo.Value) interface{} {
		return int64(e.String()[0])
	}))

	at := tf.Array(typ.Integer, 2, 3)
	b := vf.Strings(`a`, `b`, `c`).MapTo(at, func(e dgo.Value) interface{} {
		return int64(e.String()[0])
	})

	require.Equal(t, at, b.Type())
	require.Equal(t, vf.Integers(97, 98, 99), b)

	require.Panic(t, func() {
		vf.Strings(`a`, `b`, `c`, `d`).MapTo(at, func(e dgo.Value) interface{} {
			return int64(e.String()[0])
		})
	}, `size constraint`)

	require.Panic(t, func() {
		vf.Strings(`a`).MapTo(at, func(e dgo.Value) interface{} {
			return int64(e.String()[0])
		})
	}, `size constraint`)

	oat := tf.Array(tf.AnyOf(typ.Nil, typ.Integer), 0, 3)
	require.Equal(t, vf.Values(97, 98, vf.Nil), vf.Strings(`a`, `b`, `c`).MapTo(oat, func(e dgo.Value) interface{} {
		if e.Equals(`c`) {
			return nil
		}
		return vf.Integer(int64(e.String()[0]))
	}))
	require.Panic(t, func() {
		vf.Strings(`a`, `b`, `c`).MapTo(at, func(e dgo.Value) interface{} {
			if e.Equals(`c`) {
				return nil
			}
			return int64(e.String()[0])
		})
	}, `cannot be assigned`)
}

func TestArray_Pop(t *testing.T) {
	a := vf.Strings(`a`, `b`)
	require.Panic(t, func() { a.Pop() }, `Pop .* frozen`)
	a = a.Copy(false)
	l, ok := a.Pop()
	require.True(t, ok)
	require.Equal(t, `b`, l)
	require.Equal(t, a, vf.Strings(`a`))
	l, ok = a.Pop()
	require.True(t, ok)
	require.Equal(t, `a`, l)
	require.Equal(t, 0, a.Len())
	_, ok = a.Pop()
	require.False(t, ok)
}

func TestArray_Reduce(t *testing.T) {
	a := vf.Integers(1, 2, 3)
	require.Equal(t, 6, a.Reduce(nil, func(memo, v dgo.Value) interface{} {
		if memo == vf.Nil {
			return v
		}
		return memo.(dgo.Integer).GoInt() + v.(dgo.Integer).GoInt()
	}))

	require.Equal(t, vf.Nil, a.Reduce(nil, func(memo, v dgo.Value) interface{} {
		return nil
	}))
}

func TestArray_ReflectTo(t *testing.T) {
	var s []string
	a := vf.Strings(`a`, `b`)
	a.ReflectTo(reflect.ValueOf(&s).Elem())
	require.Equal(t, a, s)

	var sp *[]string
	a.ReflectTo(reflect.ValueOf(&sp).Elem())
	s = *sp
	require.Equal(t, a, s)

	var mi interface{}
	mip := &mi
	a.ReflectTo(reflect.ValueOf(mip).Elem())

	ac, ok := mi.([]string)
	require.True(t, ok)
	require.Equal(t, a, ac)

	os := []dgo.Value{vf.String(`a`), vf.Integer(23)}
	a = vf.WrapSlice(os)
	var as []dgo.Value
	a.ReflectTo(reflect.ValueOf(&as).Elem())
	require.Equal(t, os, as)

	// test that os and as is the same slice
	as[0] = vf.String(`b`)
	require.Equal(t, os, as)

	a = vf.Array(os)
	a.ReflectTo(reflect.ValueOf(&as).Elem())
	require.Equal(t, os, as)

	// test that os and as are different slices
	as[0] = vf.String(`a`)
	require.NotEqual(t, os, as)
}

func TestArray_Remove(t *testing.T) {
	s := vf.Integers(1, 2, 3, 4, 5)
	a := s.Copy(false)
	a.Remove(0)
	require.Equal(t, vf.Integers(2, 3, 4, 5), a)

	a = s.Copy(false)
	a.Remove(4)
	require.Equal(t, vf.Integers(1, 2, 3, 4), a)

	a = s.Copy(false)
	a.Remove(2)
	require.Equal(t, vf.Integers(1, 2, 4, 5), a)

	require.Panic(t, func() { s.Remove(3) }, `Remove .* frozen`)
}

func TestArray_RemoveValue(t *testing.T) {
	s := vf.Integers(1, 2, 3, 4, 5)
	a := s.Copy(false)
	require.True(t, a.RemoveValue(vf.Integer(1)))
	require.Equal(t, vf.Integers(2, 3, 4, 5), a)

	a = s.Copy(false)
	require.True(t, a.RemoveValue(vf.Integer(5)))
	require.Equal(t, vf.Integers(1, 2, 3, 4), a)

	a = s.Copy(false)
	require.True(t, a.RemoveValue(vf.Integer(3)))
	require.Equal(t, vf.Integers(1, 2, 4, 5), a)

	a = s.Copy(false)
	require.False(t, a.RemoveValue(vf.Integer(0)))
	require.Equal(t, vf.Integers(1, 2, 3, 4, 5), a)

	require.Panic(t, func() { s.RemoveValue(vf.Integer(3)) }, `RemoveValue .* frozen`)
}

func TestArray_Reject(t *testing.T) {
	require.Equal(t, vf.Values(1, 2, 4, 5), vf.Values(1, 2, vf.Nil, 4, 5).Reject(func(e dgo.Value) bool {
		return e == vf.Nil
	}))
}

func TestContainsAll(t *testing.T) {
	require.True(t, vf.Values(1, 2, 3).ContainsAll(vf.Values(2, 1)))
	require.False(t, vf.Values(1, 2).ContainsAll(vf.Values(3, 2, 1)))
	require.False(t, vf.Values(1, 2).ContainsAll(vf.Values(1, 4)))

	m := vf.Map(`c`, `C`, `d`, `D`)
	require.True(t, vf.Values(internal.NewMapEntry(`c`, `C`), internal.NewMapEntry(`d`, `D`)).ContainsAll(m))
}

func TestArray_SameValues(t *testing.T) {
	require.True(t, vf.Values().SameValues(vf.Values()))
	require.True(t, vf.Values(1, 2, 3).SameValues(vf.Values(3, 2, 1)))
	require.False(t, vf.Values(1, 2, 4).SameValues(vf.Values(3, 2, 1)))
	require.False(t, vf.Values(1, 2).SameValues(vf.Values(3, 2, 1)))
}

func TestArray_Select(t *testing.T) {
	require.Equal(t, vf.Values(1, 2, 4, 5), vf.Values(1, 2, vf.Nil, 4, 5).Select(func(e dgo.Value) bool {
		return e != vf.Nil
	}))
}

func TestArray_Slice(t *testing.T) {
	a := vf.Values(1, 2, 3, 4)
	require.Equal(t, a.Slice(0, 3), vf.Values(1, 2, 3))
	require.Equal(t, a.Slice(1, 3), vf.Values(2, 3))
	require.Same(t, a.Slice(0, 4), a)

	a = vf.MutableValues(1, 2, 3, 4)
	b := a.Slice(0, 4)
	require.Equal(t, a, b)
	require.NotSame(t, a, b)

	// Setting a value in b should not affect a
	b.Set(2, 8)
	require.Equal(t, a.Get(2), 3)
}

func TestArray_Sort(t *testing.T) {
	a := vf.Strings(`some`, `arbitrary`, `unsorted`, `words`)
	b := a.Sort()
	require.NotEqual(t, a, b)
	c := vf.Strings(`arbitrary`, `some`, `unsorted`, `words`)
	require.Equal(t, b, c)
	require.Equal(t, b, c.Sort())

	a = vf.Values(3.14, -4.2, 2)
	b = a.Sort()
	require.NotEqual(t, a, b)
	c = vf.Values(-4.2, 2, 3.14)
	require.Equal(t, b, c)
	require.Equal(t, b, c.Sort())

	a = vf.Strings(`the one and only`)
	require.Same(t, a, a.Sort())

	a = vf.Values(4.2, `hello`, -3.14)
	b = a.Sort()

	require.Equal(t, b, vf.Values(-3.14, 4.2, `hello`))
}

func TestArray_ToMap(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`, `d`)
	b := a.ToMap()
	require.Equal(t, b, vf.Map(map[string]string{`a`: `b`, `c`: `d`}))

	a = vf.Strings(`a`, `b`, `c`)
	b = a.ToMap()
	require.Equal(t, b, vf.Map(map[string]interface{}{`a`: `b`, `c`: nil}))
}

func TestArray_ToMapFromEntries(t *testing.T) {
	a := vf.Values(vf.Strings(`a`, `b`), vf.Strings(`c`, `d`))
	b, ok := a.ToMapFromEntries()
	require.True(t, ok)
	require.Equal(t, b, vf.Map(`a`, `b`, `c`, `d`))

	a = vf.Array(b)
	b, ok = a.ToMapFromEntries()
	require.True(t, ok)
	require.Equal(t, b, vf.Map(`a`, `b`, `c`, `d`))

	a = vf.Values(vf.Strings(`a`, `b`), `c`)
	_, ok = a.ToMapFromEntries()
	require.False(t, ok)
}

func TestArray_String(t *testing.T) {
	require.Equal(t, `{1,"two",3.1,true,nil}`, vf.Values(1, "two", 3.1, true, nil).String())
}

func TestArray_Unique(t *testing.T) {
	a := vf.Strings(`and`, `some`, `more`, `arbitrary`, `unsorted`, `yes`, `unsorted`, `and`, `yes`, `arbitrary`, `words`)
	b := a.Unique()
	require.NotEqual(t, a, b)
	c := vf.Strings(`and`, `some`, `more`, `arbitrary`, `unsorted`, `yes`, `words`)
	require.Equal(t, b, c)
	require.NotSame(t, b, c)
	require.Equal(t, b, c.Unique())
	require.Same(t, c, c.Unique())

	a = vf.Strings(`the one and only`)
	require.Same(t, a, a.Unique())
}

func TestArray_WithAll(t *testing.T) {
	a := vf.Values(`a`)
	c := a.WithAll(vf.Values(`b`))
	require.Equal(t, vf.Values(`a`), a)
	require.Equal(t, vf.Values(`a`, `b`), c)

	m := vf.Map(`c`, `C`, `d`, `D`)
	c = a.WithAll(m)
	require.Equal(t, vf.Values(`a`, internal.NewMapEntry(`c`, `C`), internal.NewMapEntry(`d`, `D`)), c)
	require.True(t, c.Frozen())
}
