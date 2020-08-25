package internal_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestArray_max_min(t *testing.T) {
	tp := tf.Array(2, 1)
	assert.Equal(t, tp.Min(), 1)
	assert.Equal(t, tp.Max(), 2)
}

func TestArray_negative_min(t *testing.T) {
	tp := tf.Array(-2, 3)
	assert.Equal(t, tp.Min(), 0)
	assert.Equal(t, tp.Max(), 3)
}

func TestArray_negative_min_max(t *testing.T) {
	tp := tf.Array(-2, -3)
	assert.Equal(t, tp.Min(), 0)
	assert.Equal(t, tp.Max(), 0)
}

func TestArray_explicit_unbounded(t *testing.T) {
	tp := tf.Array(0, dgo.UnboundedSize)
	assert.Equal(t, tp, typ.Array)
	assert.True(t, tp.Unbounded())
}

func TestArray_badOneArg(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	n := testNamed(0)
	tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(n), nil, nil)
	assert.Panic(t, func() { tf.Array(n) }, `illegal argument`)
}

func TestArray_badTwoArg(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	n := testNamed(0)
	tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(n), nil, nil)
	assert.Panic(t, func() { tf.Array(n, 2) }, `illegal argument 1`)
	assert.Panic(t, func() { tf.Array(typ.String, `bad`) }, `illegal argument 2`)
}

func TestArray_badThreeArg(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	n := testNamed(0)
	tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(n), nil, nil)
	assert.Panic(t, func() { tf.Array(n, 2, 2) }, `illegal argument 1`)
	assert.Panic(t, func() { tf.Array(typ.String, `bad`, 2) }, `illegal argument 2`)
	assert.Panic(t, func() { tf.Array(typ.String, 2, `bad`) }, `illegal argument 3`)
}

func TestArray_badArgCount(t *testing.T) {
	assert.Panic(t, func() { tf.Array(typ.String, 2, 2, true) }, `illegal number of arguments`)
}

func TestArrayType(t *testing.T) {
	tp := tf.Array()
	v := vf.Strings(`a`, `b`)
	assert.Instance(t, tp, v)
	assert.Assignable(t, tp, tf.AnyOf(tf.Array(tf.String(5, 5)), tf.Array(tf.String(8, 8))))
	assert.Equal(t, typ.Any, tp.ElementType())
	assert.Equal(t, 0, tp.Min())
	assert.Equal(t, dgo.UnboundedSize, tp.Max())
	assert.True(t, tp.Unbounded())
	assert.Instance(t, tp.Type(), tp)
	assert.Equal(t, `[]any`, tp.String())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Equal(t, tp.HashCode(), tp.HashCode())
}

func TestArrayType_New(t *testing.T) {
	at := tf.Array(typ.Integer)
	a := vf.Values(1, 2)
	assert.Same(t, a, vf.New(typ.Array, a))
	assert.Same(t, a, vf.New(at, a))
	assert.Same(t, a, vf.New(a.Type(), a))
	assert.Same(t, a, vf.New(tf.Tuple(typ.Integer, typ.Integer), a))
	assert.Same(t, a, vf.New(at, vf.Arguments(a)))
	assert.Equal(t, a, vf.New(at, vf.Values(1, 2)))
	assert.Panic(t, func() { vf.New(at, vf.Values(`first`, `second`)) }, `cannot be assigned`)
}

func TestSizedArrayType(t *testing.T) {
	tp := tf.Array(typ.String)
	v := vf.Strings(`a`, `b`)
	assert.Instance(t, tp, v)
	assert.NotInstance(t, tp, `a`)
	assert.NotAssignable(t, tp, typ.Array)
	assert.Assignable(t, tp, tf.AnyOf(tf.Array(tf.String(5, 5)), tf.Array(tf.String(8, 8))))
	assert.Equal(t, tp, tp)
	assert.NotEqual(t, tp, typ.Array)
	assert.NotEqual(t, tp, tf.Map(typ.String, typ.String))
	assert.Instance(t, tp.Type(), tp)
	assert.Equal(t, `[]string`, tp.String())

	tp = tf.Array(typ.Integer)
	v = vf.Strings(`a`, `b`)
	assert.NotInstance(t, tp, v)
	assert.NotAssignable(t, tp, v.Type())

	tp = tf.Array(tf.AnyOf(typ.String, typ.Integer))
	v = vf.Values(`a`, 3)
	assert.Instance(t, tp, v)
	assert.Assignable(t, tp, v.Type())
	assert.Instance(t, v.Type(), v)
	v = v.With(vf.True)
	assert.NotInstance(t, tp, v)
	assert.NotAssignable(t, tp, v.Type())

	tp = tf.Array(0, 2)
	v = vf.Strings(`a`, `b`)
	assert.Instance(t, tp, v)

	tp = tf.Array(0, 1)
	assert.NotInstance(t, tp, v)

	tp = tf.Array(2, 3)
	assert.Instance(t, tp, v)

	tp = tf.Array(3, 3)
	assert.NotInstance(t, tp, v)

	tp = tf.Array(typ.String, 2, 3)
	assert.Instance(t, tp, v)

	tp = tf.Array(typ.Integer, 2, 3)
	assert.NotInstance(t, tp, v)
	assert.NotEqual(t, 0, tp.HashCode())
	assert.NotEqual(t, tp.HashCode(), tf.Array(typ.Integer).HashCode())
	assert.Equal(t, `[2,3]int`, tp.String())

	tp = tf.Array(tf.Integer64(0, 15, true), 2, 3)
	assert.Equal(t, `[2,3]0..15`, tp.String())

	tp = tf.Array(tf.Array(2, 2), 0, 10)
	assert.Equal(t, `[0,10][2,2]any`, tp.String())
	assert.Equal(t, tf.Array(typ.Any).ReflectType(), typ.Array.ReflectType())
}

func TestExactArrayType(t *testing.T) {
	v := vf.Strings()
	tp := v.Type().(dgo.TupleType)
	assert.Equal(t, typ.Any, tp.ElementType())

	v = vf.Strings(`a`)
	tp = v.Type().(dgo.TupleType)
	et := vf.String(`a`).Type()
	assert.Equal(t, et, tp.ElementType())
	assert.Assignable(t, tp, tf.Array(et, 1, 1))
	assert.NotAssignable(t, tp, tf.Array(typ.String, 1, 1))

	v = vf.Strings(`a`, `b`)
	tp = v.Type().(dgo.TupleType)
	assert.Instance(t, tp, v)
	assert.Equal(t, tp, vf.Strings(`a`, `b`).Type())
	assert.NotInstance(t, tp, `a`)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.Array)
	assert.NotAssignable(t, tp, tf.Array(et, 1, 1))
	assert.Equal(t, vf.String(`a`).Type(), tp.ElementTypeAt(0))
	assert.Equal(t, vf.String(`b`).Type(), tp.ElementTypeAt(1))
	assert.Assignable(t, tp, tf.Tuple(vf.String(`a`).Type(), vf.String(`b`).Type()))
	assert.NotAssignable(t, tp, tf.Tuple(vf.String(`a`).Type(), vf.String(`b`).Type(), vf.String(`c`).Type()))
	assert.NotAssignable(t, tp, tf.Tuple(vf.String(`a`).Type(), typ.String))
	assert.Equal(t, 2, tp.Len())
	assert.Equal(t, 2, tp.Min())
	assert.Equal(t, 2, tp.Max())
	assert.False(t, tp.Unbounded())
	assert.False(t, tp.Variadic())
	assert.Equal(t, tf.Array(typ.String), typ.Generic(tp))
	assert.NotAssignable(t, tp, tf.AnyOf(tf.Array(tf.String(5, 5)), tf.Array(tf.String(8, 8))))
	assert.Equal(t, tp, tp)
	assert.Equal(t, vf.Values(vf.String(`a`).Type(), vf.String(`b`).Type()), tp.ElementTypes())
	assert.NotEqual(t, tp, typ.Array)
	assert.Instance(t, tp.Type(), tp)
	assert.Equal(t, `{"a","b"}`, tp.String())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.NotEqual(t, tp.HashCode(), tf.Array(typ.Integer).HashCode())
}

func TestArrayElementType_singleElement(t *testing.T) {
	v := vf.Strings(`hello`)
	at := v.Type()
	et := at.(dgo.ArrayType).ElementType()
	assert.Assignable(t, at, at)
	tp := tf.Array(typ.String)
	assert.Assignable(t, tp, at)
	assert.NotAssignable(t, at, tp)
	assert.Assignable(t, et, vf.String(`hello`).Type())
	assert.NotAssignable(t, et, vf.String(`hey`).Type())
	assert.Instance(t, et, `hello`)
	assert.NotInstance(t, et, `world`)
	assert.Equal(t, et, vf.Strings(`hello`).Type().(dgo.ArrayType).ElementType())
	assert.NotEqual(t, et, vf.Strings(`hello`).Type().(dgo.ArrayType))
	assert.NotEqual(t, 0, et.HashCode())
	assert.Equal(t, et.HashCode(), et.HashCode())
}

func TestArrayElementType_multipleElements(t *testing.T) {
	v := vf.Strings(`hello`, `world`)
	at := v.Type()
	et := at.(dgo.ArrayType).ElementType()
	assert.Assignable(t, at, at)
	tp := tf.Array(typ.String)
	assert.Assignable(t, tp, at)
	assert.NotAssignable(t, at, tp)
	assert.NotAssignable(t, et, vf.String(`hello`).Type())
	assert.NotAssignable(t, et, vf.String(`world`).Type())
	assert.NotAssignable(t, vf.String(`hello`).Type(), et)
	assert.NotAssignable(t, vf.String(`world`).Type(), et)
	assert.Assignable(t, vf.Strings(`hello`, `world`).Type(), at)
	assert.Assignable(t, vf.Strings(`hello`, `world`).Type().(dgo.ArrayType).ElementType(), et)
	assert.Assignable(t, vf.Strings(`world`, `hello`).Type().(dgo.ArrayType).ElementType(), et)
	assert.Assignable(t, tf.Array(2, 2), at)

	et.(dgo.TernaryType).Operands().Each(func(v dgo.Value) {
		t.Helper()
		assert.Instance(t, typ.String, v)
	})
	assert.Instance(t, et.Type(), et)
	assert.Equal(t, `"hello"&"world"`, et.String())
}

func TestMutableValues_withoutNil(t *testing.T) {
	a := vf.MutableValues(nil)
	assert.True(t, vf.Nil == a.Get(0))
}

func TestMutableValues_maxSizeMismatch(t *testing.T) {
	tp := tf.Array(0, 1)
	a := vf.Values(true)
	assert.Instance(t, tp, a)
	assert.NotInstance(t, tp, a.With(false))
	assert.NotInstance(t, tp, a.WithAll(vf.Values(false)))
	assert.NotInstance(t, tp, a.WithValues(false))
	assert.Instance(t, tp, a.WithAll(vf.Values()))
	assert.Instance(t, tp, a.WithValues())
}

func TestMutableValues_minSizeMismatch(t *testing.T) {
	a := vf.MutableValues(true)
	tp := tf.ParseType(`[1,1]bool`)
	assert.Instance(t, tp, a)
	a.Remove(0)
	assert.NotInstance(t, tp, a)
}

func TestMutableArray(t *testing.T) {
	s := []dgo.Value{nil}
	a := vf.WrapSlice(s)
	assert.True(t, vf.Nil == s[0])
	a.Set(0, `hello`)
	assert.Equal(t, `hello`, s[0])
}

func TestMutableArray_nilSlice(t *testing.T) {
	a := vf.WrapSlice(nil)
	a.Add(3)
	assert.Equal(t, 3, a.Get(0))
}

func TestArray(t *testing.T) {
	a := vf.Array([]int{1, 2})
	assert.Equal(t, 2, a.Len())
	assert.Equal(t, 1, a.Get(0))
	assert.Equal(t, 2, a.Get(1))
	assert.Same(t, a, vf.Array(a))
	assert.Same(t, a, vf.Array(reflect.ValueOf(a)))
}

func TestArray_FrozenCopy(t *testing.T) {
	a := vf.MutableValues()
	a.Add(1)
	a.Set(0, 2)
	assert.Equal(t, 2, a.Get(0))

	f := a.Copy(true)
	assert.Panic(t, func() { f.Set(0, 1) }, `Set .* frozen`)
}

func TestArray_selfReference(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`x=[](string|x)`).(dgo.ArrayType)
	d := vf.MutableValues(tp, `hello`)
	d.Add(d)
	assert.Instance(t, tp, d)

	internal.ResetDefaultAliases()
	t2 := tf.ParseType(`x=[](string|[](string|x))`)
	assert.Assignable(t, tp, t2)
}

func TestArray_recursiveFreeze(t *testing.T) {
	a := vf.Array([]dgo.Value{vf.MutableValues(`b`)})
	assert.True(t, a.Get(0).(dgo.Array).Frozen())
}

func TestArray_recursiveReflectiveFreeze(t *testing.T) {
	a := vf.Value(
		reflect.ValueOf([]interface{}{
			reflect.ValueOf(vf.MutableValues(`b`))})).(dgo.Array)
	assert.False(t, a.Get(0).(dgo.Array).Frozen())
}

func TestArray_replaceNil(t *testing.T) {
	a := vf.Array([]dgo.Value{nil})
	assert.True(t, vf.Nil == a.Get(0))
}

func TestArray_fromReflected(t *testing.T) {
	a := vf.Value([]interface{}{`a`, 1, nil}).(dgo.Array)
	assert.False(t, a.Frozen())
	assert.True(t, vf.Nil == a.Get(2))
}

func TestArray_Add(t *testing.T) {
	a := vf.Values(`a`)
	assert.Panic(t, func() { a.Add(vf.Value(`b`)) }, `Add .* frozen`)
	m := a.Copy(false)
	m.Add(vf.Value(`b`))
	assert.Equal(t, vf.Values(`a`), a)
	assert.Equal(t, vf.Values(`a`, `b`), m)
}

func TestArray_AddAll(t *testing.T) {
	a := vf.Values(`a`)
	assert.Panic(t, func() { a.AddAll(vf.Values(`b`)) }, `AddAll .* frozen`)
	c := a.Copy(false)
	c.AddAll(vf.Values(`b`))
	assert.Equal(t, vf.Values(`a`), a)
	assert.Equal(t, vf.Values(`a`, `b`), c)

	m := vf.Map(`c`, `C`, `d`, `D`)
	c.AddAll(m)
	assert.Equal(t, vf.Values(`a`, `b`, internal.NewMapEntry(`c`, `C`), internal.NewMapEntry(`d`, `D`)), c)
}

func TestArray_AddValues(t *testing.T) {
	a := vf.Values(`a`)
	assert.Panic(t, func() { a.AddValues(`b`) }, `AddValues .* frozen`)
	m := a.Copy(false)
	m.AddValues(`b`)
	assert.Equal(t, vf.Values(`a`), a)
	assert.Equal(t, vf.Values(`a`, `b`), m)
}

func TestArray_All(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)
	i := 0
	assert.True(t, a.All(func(e dgo.Value) bool {
		i++
		return e.Equals(`a`) || e.Equals(`b`) || e.Equals(`c`)
	}))
	assert.Equal(t, 3, i)

	i = 0
	assert.False(t, a.All(func(e dgo.Value) bool {
		i++
		return e.Equals(`a`)
	}))
	assert.Equal(t, 2, i)
}

func TestArray_Any(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)
	i := 0
	assert.True(t, a.Any(func(e dgo.Value) bool {
		i++
		return e.Equals(`b`)
	}))
	assert.Equal(t, 2, i)

	i = 0
	assert.False(t, a.Any(func(e dgo.Value) bool {
		i++
		return e.Equals(`d`)
	}))
	assert.Equal(t, 3, i)
}

func TestArray_One(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)
	i := 0
	assert.True(t, a.One(func(e dgo.Value) bool {
		i++
		return e.Equals(`b`)
	}))
	assert.Equal(t, 3, i)

	a = vf.Strings(`a`, `b`, `c`, `b`)
	i = 0
	assert.False(t, a.One(func(e dgo.Value) bool {
		i++
		return e.Equals(`b`)
	}))
	assert.Equal(t, 4, i)

	i = 0
	assert.False(t, a.One(func(e dgo.Value) bool {
		i++
		return e.Equals(`d`)
	}))
	assert.Equal(t, 4, i)
}

func TestArray_CompareTo(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)

	c, ok := a.CompareTo(a)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = a.CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	_, ok = a.CompareTo(vf.String(`a`))
	assert.False(t, ok)

	b := vf.Strings(`a`, `b`, `c`)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	b = vf.Strings(`a`, `b`, `c`, `d`)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	b = vf.Values(`a`, `b`, `d`, `d`)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	b = vf.Values(`a`, `b`, 3, `d`)
	_, ok = a.CompareTo(b)
	assert.False(t, ok)

	b = vf.Strings(`a`, `b`)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	b = vf.Strings(`a`, `b`, `d`)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	b = vf.Strings(`a`, `b`, `b`)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	a = vf.MutableValues(`a`, `b`)
	a.Add(a)
	b = vf.MutableValues(`a`, `b`)
	b.Add(b)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	b = vf.MutableValues(`a`, `b`)
	b.Add(a)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	a = vf.Values(`a`, 1, nil)
	b = vf.Values(`a`, 1, nil)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	b = vf.Values(`a`, 1, 2)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	a = vf.Values(`a`, 1, 2)
	b = vf.Values(`a`, 1, nil)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	a = vf.Values(`a`, 1, 2)
	m := vf.Map(`a`, 1)
	_, ok = a.CompareTo(m)
	assert.False(t, ok)

	a = vf.Values(`a`, 1, []int{2})
	b = vf.Values(`a`, 1, 2)
	_, ok = a.CompareTo(b)
	assert.False(t, ok)

	a = vf.Values(m)
	b = vf.Values(`a`)
	_, ok = a.CompareTo(b)
	assert.False(t, ok)
}

func TestArray_Copy(t *testing.T) {
	a := vf.Values(`a`, `b`, vf.MutableValues(`c`))
	assert.Same(t, a, a.Copy(true))
	assert.True(t, a.Get(2).(dgo.Mutability).Frozen())

	c := a.Copy(false)
	assert.False(t, c.Frozen())
	assert.NotSame(t, c, c.Copy(false))

	c = c.Copy(true)
	assert.True(t, c.Frozen())
	assert.Same(t, c, c.Copy(true))
}

func TestArray_Flatten_mutable(t *testing.T) {
	a := vf.MutableValues(`a`, `b`, vf.Values(`c`, `d`), vf.Values(`e`, vf.Values(`f`, `g`)))
	af := a.Flatten()
	b := vf.MutableValues(`a`, `b`, `c`, `d`, `e`, `f`, `g`)
	assert.Equal(t, b, af)
	assert.Same(t, b, b.Flatten())
	assert.False(t, af.Frozen())
}

func TestArray_Flatten_frozen(t *testing.T) {
	a := vf.Values(`a`, `b`, vf.Values(`c`, `d`), vf.Values(`e`, vf.Values(`f`, `g`)))
	af := a.Flatten()
	b := vf.Values(`a`, `b`, `c`, `d`, `e`, `f`, `g`)
	assert.Equal(t, b, af)
	assert.Same(t, b, b.Flatten())
	assert.True(t, af.Frozen())
}

func TestArray_FromReflected(t *testing.T) {
	vs := []dgo.Value{vf.Int64(2), vf.String(`b`)}
	a := internal.ArrayFromReflected(reflect.ValueOf(vs), false).(dgo.Array)
	assert.Equal(t, reflect.ValueOf(a.GoSlice()).Pointer(), reflect.ValueOf(vs).Pointer())

	a = internal.ArrayFromReflected(reflect.ValueOf(vs), true).(dgo.Array)
	assert.NotEqual(t, reflect.ValueOf(a.GoSlice()).Pointer(), reflect.ValueOf(vs).Pointer())

	vs = append(vs, vf.MutableValues(3, 5))
	a = internal.ArrayFromReflected(reflect.ValueOf(vs), true).(dgo.Array)
	assert.True(t, a.Get(2).(dgo.Mutability).Frozen())
}

func TestArray_Equal(t *testing.T) {
	a := vf.Values(1, 2)
	assert.True(t, a.Equals(a))

	b := vf.Values(1, nil)
	assert.False(t, a.Equals(b))

	a = vf.Values(1, nil)
	assert.True(t, a.Equals(b))

	b = vf.Values(1, 2)
	assert.False(t, a.Equals(b))

	a = vf.Values(1, []int{2})
	assert.False(t, a.Equals(b))

	b = vf.Values(1, map[int]int{2: 1})
	assert.False(t, a.Equals(b))

	b = vf.Values(1, `2`)
	assert.False(t, a.Equals(b))

	// Values containing themselves.
	a = vf.MutableValues(`2`)
	a.Add(a)

	b = vf.MutableValues(`2`)
	b.Add(b)

	m := vf.MutableMap()
	m.Put(`me`, m)
	a.Add(m)
	b.Add(m)
	assert.True(t, a.Equals(b))
	assert.Equal(t, a.HashCode(), b.HashCode())
}

func TestArray_EachWithIndex(t *testing.T) {
	ni := 0
	vf.Values(1, 2, 3).EachWithIndex(func(v dgo.Value, i int) {
		assert.Equal(t, ni, i)
		assert.Equal(t, v, i+1)
		ni++
	})
	assert.Equal(t, 3, ni)
}

func TestArray_ElementType(t *testing.T) {
	assert.Equal(t, typ.Any, vf.MutableValues().(dgo.ArrayType).ElementType())
	v := vf.String(`hello`)
	assert.Same(t, v, vf.MutableValues(v).(dgo.ArrayType).ElementType())
}

func TestArray_Find(t *testing.T) {
	v := vf.Values(`a`, `b`, 3, `d`).Find(func(v dgo.Value) interface{} {
		if v.Equals(3) {
			return `three`
		}
		return nil
	})
	assert.Equal(t, v, `three`)
}

func TestArray_Find_notFound(t *testing.T) {
	v := vf.Values(`a`, `b`, `d`).Find(func(v dgo.Value) interface{} { return nil })
	assert.True(t, v == nil)
}

func TestArray_FrozenEqual(t *testing.T) {
	f := vf.Values(1, 2, 3)
	assert.True(t, f.Frozen(), `not frozen`)

	a := f.Copy(false)
	assert.False(t, a.Frozen(), `frozen`)
	assert.Equal(t, f, a)
	assert.Equal(t, a, f)

	a = a.Copy(true)
	assert.True(t, a.Frozen(), `not frozen`)
	assert.Same(t, a, a.Copy(true))

	b := a.Copy(false)
	assert.NotSame(t, a, b)
	assert.NotSame(t, b, b.Copy(true))
	assert.NotSame(t, b, b.Copy(false))
}

func TestArray_Format(t *testing.T) {
	a := vf.Values(1, 3)
	assert.Equal(t, `[]int{1, 3}`, fmt.Sprintf("%#v", a))
}

func TestArray_Format_namedElementType(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tf.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil, nil)
	a := vf.Values(testNamed(1), testNamed(3))
	assert.Equal(t, `[]testNamed{1, 3}`, fmt.Sprintf("%#v", a))
}

func TestArray_Format_nativeElementType(t *testing.T) {
	a := vf.Values(&testStruct{A: `a`}, &testStruct{A: `b`})
	assert.Equal(t, `[]*internal_test.testStruct{&internal_test.testStruct{A:"a"}, &internal_test.testStruct{A:"b"}}`, fmt.Sprintf("%#v", a))
}

func TestArray_IndexOf(t *testing.T) {
	a := vf.Values(1, nil, 3)
	assert.Equal(t, 2, a.IndexOf(3))
	assert.Equal(t, 1, a.IndexOf(nil))
	assert.Equal(t, 1, a.IndexOf(vf.Nil))
}

func TestArray_Insert(t *testing.T) {
	a := vf.Values(`a`)
	assert.Panic(t, func() { a.Insert(0, vf.Value(`b`)) }, `Insert .* frozen`)
	m := a.Copy(false)
	m.Insert(0, vf.Value(`b`))
	assert.Equal(t, vf.Values(`a`), a)
	assert.Equal(t, vf.Values(`b`, `a`), m)
}

func TestArray_Map(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`)
	assert.Equal(t, vf.Strings(`d`, `e`, `f`), a.Map(func(e dgo.Value) interface{} {
		return string([]byte{e.String()[0] + 3})
	}))
	assert.Equal(t, vf.Values(vf.Nil, vf.Nil, vf.Nil), a.Map(func(e dgo.Value) interface{} {
		return nil
	}))

	ma := a.Map(func(e dgo.Value) interface{} {
		return vf.MutableValues(e)
	})
	assert.True(t, ma.All(func(e dgo.Value) bool {
		return e.(dgo.Mutability).Frozen()
	}))
}

func TestArray_Pop(t *testing.T) {
	a := vf.Strings(`a`, `b`)
	assert.Panic(t, func() { a.Pop() }, `Pop .* frozen`)
	a = a.Copy(false)
	l, ok := a.Pop()
	assert.True(t, ok)
	assert.Equal(t, `b`, l)
	assert.Equal(t, a, vf.Strings(`a`))
	l, ok = a.Pop()
	assert.True(t, ok)
	assert.Equal(t, `a`, l)
	assert.Equal(t, 0, a.Len())
	_, ok = a.Pop()
	assert.False(t, ok)
}

func TestArray_Reduce(t *testing.T) {
	a := vf.Integers(1, 2, 3)
	assert.Equal(t, 6, a.Reduce(nil, func(memo, v dgo.Value) interface{} {
		if memo == vf.Nil {
			return v
		}
		return memo.(dgo.Integer).GoInt() + v.(dgo.Integer).GoInt()
	}))
	assert.Equal(t, vf.Nil, a.Reduce(nil, func(memo, v dgo.Value) interface{} {
		return nil
	}))
}

func TestArray_ReflectTo(t *testing.T) {
	var s []string
	a := vf.Strings(`a`, `b`)
	a.ReflectTo(reflect.ValueOf(&s).Elem())
	assert.Equal(t, a, s)

	var sp *[]string
	a.ReflectTo(reflect.ValueOf(&sp).Elem())
	s = *sp
	assert.Equal(t, a, s)

	var mi interface{}
	mip := &mi
	a.ReflectTo(reflect.ValueOf(mip).Elem())

	ac, ok := mi.([]string)
	require.True(t, ok)
	assert.Equal(t, a, ac)

	os := []dgo.Value{vf.String(`a`), vf.Int64(23)}
	a = vf.WrapSlice(os)
	var as []dgo.Value
	a.ReflectTo(reflect.ValueOf(&as).Elem())
	assert.Equal(t, os, as)

	// test that os and as is the same slice
	as[0] = vf.String(`b`)
	assert.Equal(t, os, as)

	a = vf.Array(os)
	a.ReflectTo(reflect.ValueOf(&as).Elem())
	assert.Equal(t, os, as)

	// test that os and as are different slices
	as[0] = vf.String(`a`)
	assert.NotEqual(t, os, as)
}

func TestArray_ReflectTo_nestedStructPointer(t *testing.T) {
	m := vf.Values(vf.Map(
		`First`, 1,
		`Second`, 2.0))

	type structA struct {
		First  int
		Second float64
	}

	var sb []*structA
	m.ReflectTo(reflect.ValueOf(&sb))
	require.Equal(t, 1, len(sb))
	assert.Equal(t, 1, sb[0].First)
	assert.Equal(t, 2.0, sb[0].Second)
}

func TestArray_Remove(t *testing.T) {
	s := vf.Integers(1, 2, 3, 4, 5)
	a := s.Copy(false)
	a.Remove(0)
	assert.Equal(t, vf.Integers(2, 3, 4, 5), a)

	a = s.Copy(false)
	a.Remove(4)
	assert.Equal(t, vf.Integers(1, 2, 3, 4), a)

	a = s.Copy(false)
	a.Remove(2)
	assert.Equal(t, vf.Integers(1, 2, 4, 5), a)
	assert.Panic(t, func() { s.Remove(3) }, `Remove .* frozen`)
}

func TestArray_RemoveValue(t *testing.T) {
	s := vf.Integers(1, 2, 3, 4, 5)
	a := s.Copy(false)
	assert.True(t, a.RemoveValue(vf.Int64(1)))
	assert.Equal(t, vf.Integers(2, 3, 4, 5), a)

	a = s.Copy(false)
	assert.True(t, a.RemoveValue(vf.Int64(5)))
	assert.Equal(t, vf.Integers(1, 2, 3, 4), a)

	a = s.Copy(false)
	assert.True(t, a.RemoveValue(vf.Int64(3)))
	assert.Equal(t, vf.Integers(1, 2, 4, 5), a)

	a = s.Copy(false)
	assert.False(t, a.RemoveValue(vf.Int64(0)))
	assert.Equal(t, vf.Integers(1, 2, 3, 4, 5), a)
	assert.Panic(t, func() { s.RemoveValue(vf.Int64(3)) }, `RemoveValue .* frozen`)
}

func TestArray_Reject(t *testing.T) {
	vr := vf.Values(1, 2, vf.Nil, 4, 5).Reject(func(e dgo.Value) bool {
		return e == vf.Nil
	})
	assert.Equal(t, vf.Values(1, 2, 4, 5), vr)
	assert.True(t, vr.Frozen())

	vr = vf.MutableValues(1, 2, vf.Nil, 4, 5).Reject(func(e dgo.Value) bool {
		return e == vf.Nil
	})
	assert.Equal(t, vf.Values(1, 2, 4, 5), vr)
	assert.False(t, vr.Frozen())
}

func TestArray_ContainsAll(t *testing.T) {
	assert.True(t, vf.Values(1, 2, 3).ContainsAll(vf.Values(2, 1)))
	assert.False(t, vf.Values(1, 2).ContainsAll(vf.Values(3, 2, 1)))
	assert.False(t, vf.Values(1, 2).ContainsAll(vf.Values(1, 4)))

	m := vf.Map(`c`, `C`, `d`, `D`)
	assert.True(t, vf.Values(internal.NewMapEntry(`c`, `C`), internal.NewMapEntry(`d`, `D`)).ContainsAll(m))
}

func TestArray_SameValues(t *testing.T) {
	assert.True(t, vf.Values().SameValues(vf.Values()))
	assert.True(t, vf.Values(1, 2, 3).SameValues(vf.Values(3, 2, 1)))
	assert.False(t, vf.Values(1, 2, 4).SameValues(vf.Values(3, 2, 1)))
	assert.False(t, vf.Values(1, 2).SameValues(vf.Values(3, 2, 1)))
}

func TestArray_Select(t *testing.T) {
	vr := vf.Values(1, 2, vf.Nil, 4, 5).Select(func(e dgo.Value) bool {
		return e != vf.Nil
	})
	assert.Equal(t, vf.Values(1, 2, 4, 5), vr)
	assert.True(t, vr.Frozen())

	vr = vf.MutableValues(1, 2, vf.Nil, 4, 5).Select(func(e dgo.Value) bool {
		return e != vf.Nil
	})
	assert.Equal(t, vf.Values(1, 2, 4, 5), vr)
	assert.False(t, vr.Frozen())
}

func TestArray_Slice(t *testing.T) {
	a := vf.Values(1, 2, 3, 4)
	assert.Equal(t, a.Slice(0, 3), vf.Values(1, 2, 3))
	assert.Equal(t, a.Slice(1, 3), vf.Values(2, 3))
	assert.Same(t, a.Slice(0, 4), a)

	a = vf.MutableValues(1, 2, 3, 4)
	b := a.Slice(0, 4)
	assert.Equal(t, a, b)
	assert.NotSame(t, a, b)

	// Setting a value in b should affect a
	b.Set(2, 8)
	assert.Equal(t, a.Get(2), 8)
}

func TestArray_Sort(t *testing.T) {
	a := vf.Strings(`some`, `arbitrary`, `unsorted`, `words`)
	b := a.Sort()
	assert.NotEqual(t, a, b)
	c := vf.Strings(`arbitrary`, `some`, `unsorted`, `words`)
	assert.Equal(t, b, c)
	assert.Equal(t, b, c.Sort())

	a = vf.Values(3.14, -4.2, 2.0)
	b = a.Sort()
	assert.NotEqual(t, a, b)
	c = vf.Values(-4.2, 2.0, 3.14)
	assert.Equal(t, b, c)
	assert.Equal(t, b, c.Sort())

	a = vf.Strings(`the one and only`)
	assert.Same(t, a, a.Sort())

	a = vf.Values(4.2, `hello`, -3.14)
	b = a.Sort()
	assert.Equal(t, b, vf.Values(-3.14, 4.2, `hello`))
	assert.True(t, b.Frozen())

	a = vf.MutableValues(4.2, `hello`, -3.14)
	b = a.Sort()
	assert.Equal(t, b, vf.Values(-3.14, 4.2, `hello`))
	assert.False(t, b.Frozen())

	a = vf.Values(`hello`)
	assert.Same(t, a, a.Sort())

	a = vf.MutableValues(`hello`)
	assert.Same(t, a, a.Sort())
}

func TestArray_ToMap(t *testing.T) {
	a := vf.Strings(`a`, `b`, `c`, `d`)
	b := a.ToMap()
	assert.Equal(t, b, vf.Map(map[string]string{`a`: `b`, `c`: `d`}))

	a = vf.Strings(`a`, `b`, `c`)
	b = a.ToMap()
	assert.Equal(t, b, vf.Map(map[string]interface{}{`a`: `b`, `c`: nil}))
}

func TestArray_ToMapFromEntries(t *testing.T) {
	a := vf.Values(vf.Strings(`a`, `b`), vf.Strings(`c`, `d`))
	b, ok := a.ToMapFromEntries()
	assert.True(t, ok)
	assert.Equal(t, b, vf.Map(`a`, `b`, `c`, `d`))

	a = vf.Array(b)
	b, ok = a.ToMapFromEntries()
	assert.True(t, ok)
	assert.Equal(t, b, vf.Map(`a`, `b`, `c`, `d`))
	assert.True(t, b.Frozen())

	a = a.Copy(false)
	b, ok = a.ToMapFromEntries()
	assert.True(t, ok)
	assert.Equal(t, b, vf.Map(`a`, `b`, `c`, `d`))
	assert.False(t, b.Frozen())

	a = vf.Values(vf.Strings(`a`, `b`), `c`)
	_, ok = a.ToMapFromEntries()
	assert.False(t, ok)

	a = vf.Values(vf.Strings(`a`, `b`), vf.Strings(`c`))
	_, ok = a.ToMapFromEntries()
	assert.False(t, ok)
}

func TestArray_String(t *testing.T) {
	assert.Equal(t, `{1,"two",3.1,true,nil}`, vf.Values(1, "two", 3.1, true, nil).String())
}

func TestArray_Unique(t *testing.T) {
	a := vf.MutableValues(`and`, `some`, `more`, `arbitrary`, `unsorted`, `yes`, `unsorted`, `and`, `yes`, `arbitrary`, `words`)
	b := a.Unique()
	assert.NotEqual(t, a, b)
	assert.False(t, b.Frozen())
	c := vf.Strings(`and`, `some`, `more`, `arbitrary`, `unsorted`, `yes`, `words`)
	assert.Equal(t, b, c)
	assert.NotSame(t, b, c)
	assert.Equal(t, b, c.Unique())
	assert.Same(t, c, c.Unique())

	a = vf.Strings(`the one and only`)
	assert.Same(t, a, a.Unique())
}

func TestArray_With(t *testing.T) {
	a := vf.Values(`a`).With(vf.MutableValues(`b`))
	assert.True(t, a.Get(1).(dgo.Mutability).Frozen())

	a = vf.MutableValues(`a`).With(vf.MutableValues(`b`))
	assert.False(t, a.Get(1).(dgo.Mutability).Frozen())
}

func TestArray_WithAll_mutable(t *testing.T) {
	a := vf.MutableValues(`a`)
	assert.Same(t, a, a.WithAll(vf.Values()))
	c := a.WithAll(vf.MutableValues(vf.MutableValues(`b`)))
	assert.False(t, c.Get(1).(dgo.Mutability).Frozen())
	assert.Equal(t, vf.Values(`a`), a)
	assert.Equal(t, vf.Values(`a`, vf.Values(`b`)), c)

	m := vf.MutableMap(`c`, `C`, `d`, `D`)
	c = a.WithAll(m)
	assert.Equal(t, vf.Values(`a`, internal.NewMapEntry(`c`, `C`), internal.NewMapEntry(`d`, `D`)), c)
	assert.False(t, c.Frozen())
}

func TestArray_WithAll_frozen(t *testing.T) {
	a := vf.Values(`a`)
	assert.Same(t, a, a.WithAll(vf.Values()))
	c := a.WithAll(vf.MutableValues(vf.MutableValues(`b`)))
	assert.True(t, c.Get(1).(dgo.Mutability).Frozen())
	assert.Equal(t, vf.Values(`a`), a)
	assert.Equal(t, vf.Values(`a`, vf.Values(`b`)), c)

	m := vf.Map(`c`, `C`, `d`, `D`)
	c = a.WithAll(m)
	assert.Equal(t, vf.Values(`a`, internal.NewMapEntry(`c`, `C`), internal.NewMapEntry(`d`, `D`)), c)
	assert.True(t, c.Frozen())
}

func TestArray_WithValues_mutable(t *testing.T) {
	a := vf.MutableValues(`a`)
	assert.Same(t, a, a.WithValues())
	c := a.WithValues(vf.MutableValues(`b`))
	assert.False(t, c.Get(1).(dgo.Mutability).Frozen())
	assert.Equal(t, vf.Values(`a`), a)
	assert.Equal(t, vf.Values(`a`, vf.Values(`b`)), c)
}
