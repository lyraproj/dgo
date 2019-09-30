package internal_test

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestTyped(t *testing.T) {
	// value type for the map
	mt := newtype.Map(typ.String, newtype.AnyOf(typ.String, newtype.IntegerRange(0, 15, true)), 0, 2)

	m := vf.MutableMap(mt)
	m.PutAll(vf.Map(map[string]interface{}{
		`first`:  1,
		`second`: `two`,
	}))
	require.Panic(t, func() { m.Put(`third`, 3.2) }, internal.IllegalAssignment(mt.ValueType(), vf.Value(3.2)))
	require.Panic(t, func() { m.Put(1, 2) }, internal.IllegalAssignment(mt.KeyType(), vf.Value(1)))
	require.Panic(t, func() { m.Put(`third`, 2) }, internal.IllegalSize(mt, 3))
}

func TestMap_ValueType(t *testing.T) {
	m1 := vf.Map(
		`first`, 1,
		`second`, 2).Type().(dgo.MapType).ValueType()
	m2 := vf.Map(
		`one`, 1,
		`two`, 2).Type().(dgo.MapType).ValueType()
	m3 := vf.Map(
		`one`, 1,
		`two`, 3).Type().(dgo.MapType).ValueType()
	m4 := vf.Map(
		`two`, 3,
		`three`, 3,
		`four`, 3).Type().(dgo.MapType).ValueType()
	require.Assignable(t, m1, m2)
	require.NotAssignable(t, m1, m3)
	require.Assignable(t, newtype.IntegerRange(1, 2, true), m1)
	require.NotAssignable(t, newtype.IntegerRange(2, 3, true), m1)

	require.NotAssignable(t, m2, vf.Integer(2).Type())
	require.Assignable(t, m4, vf.Integer(3).Type())

	require.Equal(t, m1, m2)
	require.NotEqual(t, m1, m3)
	require.NotEqual(t, m4, vf.Integer(3).Type())
	require.NotEqual(t, m1, m4)

	require.True(t, m1.HashCode() > 0)
	require.Equal(t, m1.HashCode(), m1.HashCode())
	vm := m1.Type()
	require.Instance(t, vm, m1)

	require.True(t, `1&2` == m1.String())
}

func TestNewMapType_DefaultType(t *testing.T) {
	mt := newtype.Map()
	require.Same(t, typ.Map, mt)

	mt = newtype.Map(typ.Any, typ.Any)
	require.Same(t, typ.Map, mt)

	mt = newtype.Map(typ.Any, typ.Any, 0, math.MaxInt64)
	require.Same(t, typ.Map, mt)

	m1 := vf.Map(
		`a`, 1,
		`b`, 2,
	)
	require.Assignable(t, mt, mt)
	require.NotAssignable(t, mt, typ.String)
	require.Instance(t, mt, m1)
	require.NotInstance(t, mt, `a`)

	require.Equal(t, mt, mt)
	require.NotEqual(t, mt, typ.String)

	require.Equal(t, mt.KeyType(), typ.Any)
	require.Equal(t, mt.ValueType(), typ.Any)

	require.Equal(t, mt.Min(), 0)
	require.Equal(t, mt.Max(), math.MaxInt64)
	require.True(t, mt.Unbounded())

	require.True(t, mt.HashCode() > 0)
	require.Equal(t, mt.HashCode(), mt.HashCode())

	vm := mt.Type()
	require.Instance(t, vm, mt)

	require.Equal(t, `map[any]any`, mt.String())
}

func TestMap_ExactType(t *testing.T) {
	m1 := vf.Map(
		`a`, 3,
		`b`, 4)
	t1 := m1.Type().(dgo.MapType)
	m2 := vf.Map(
		`a`, 1,
		`b`, 2)
	t2 := m2.Type().(dgo.MapType)
	t3 := vf.Map(`b`, 2).Type().(dgo.MapType)
	require.Equal(t, 2, t1.Min())
	require.Equal(t, 2, t1.Max())
	require.False(t, t1.Unbounded())
	require.Assignable(t, t1, t1)
	require.NotAssignable(t, t1, t2)
	require.Instance(t, t1, m1)
	require.NotInstance(t, t1, m2)
	require.NotInstance(t, t1, `a`)
	require.Assignable(t, typ.Map, t1)
	require.NotAssignable(t, t1, typ.Map)
	require.NotAssignable(t, t1, t3)
	require.Assignable(t, newtype.Map(typ.String, typ.Integer), t1)
	require.NotAssignable(t, newtype.Map(typ.String, typ.String), t1)
	require.NotAssignable(t, newtype.Map(typ.String, typ.Integer, 3, 3), t1)
	require.NotAssignable(t, t1, newtype.Map(typ.String, typ.Integer))
	require.NotAssignable(t, typ.String, t1)

	require.NotEqual(t, t1, t2)
	require.NotEqual(t, t1, t3)
	require.NotEqual(t, t1, typ.String)
	require.NotEqual(t, newtype.Map(typ.String, typ.Integer), t1)

	require.True(t, t1.HashCode() > 0)
	require.Equal(t, t1.HashCode(), t1.HashCode())
	vm := t1.Type()
	require.Instance(t, vm, t1)

	require.Equal(t, `{"a":3,"b":4}`, t1.String())
}

func TestMap_SizedType(t *testing.T) {
	mt := newtype.Map(typ.String, typ.Integer)

	require.Equal(t, mt, mt)
	require.Equal(t, mt, newtype.Map(typ.String, typ.Integer))
	require.NotEqual(t, mt, newtype.Map(typ.String, typ.Float))
	require.NotEqual(t, mt, newtype.Array(typ.String))

	m1 := vf.Map(
		`a`, 1,
		`b`, 2)
	m2 := vf.Map(
		`a`, 1,
		`b`, 2.0)
	m3 := vf.Map(
		`a`, 1,
		`b`, 2,
		`c`, 3)
	require.Assignable(t, mt, mt)
	require.NotAssignable(t, mt, typ.Any)
	require.Instance(t, mt, m1)
	require.NotInstance(t, mt, m2)
	require.NotInstance(t, mt, `a`)

	mtz := newtype.Map(typ.String, typ.Integer, 3, 3)
	require.Instance(t, mtz, m3)
	require.NotInstance(t, mtz, m1)

	require.Assignable(t, mt, mtz)
	require.NotAssignable(t, mtz, mt)
	require.NotEqual(t, mt, mtz)
	require.NotEqual(t, mt, typ.Any)

	mta := newtype.Map(typ.Any, typ.Any, 3, 3)
	require.NotInstance(t, mta, m1)
	require.Instance(t, mta, m3)

	mtva := newtype.Map(typ.String, typ.Any)
	require.Instance(t, mtva, m1)
	require.Instance(t, mtva, m2)
	require.Instance(t, mtva, m3)

	mtka := newtype.Map(typ.Any, typ.Integer)
	require.Instance(t, mtka, m1)
	require.NotInstance(t, mtka, m2)
	require.Instance(t, mtka, m3)

	require.True(t, mt.HashCode() > 0)
	require.Equal(t, mt.HashCode(), mt.HashCode())
	require.NotEqual(t, mt.HashCode(), mtz.HashCode())
	vm := mt.Type()
	require.Instance(t, vm, mt)

	require.Equal(t, `map[string]int`, mt.String())
	require.Equal(t, `map[string,3,3]int`, mtz.String())

	require.Equal(t, mta.ReflectType(), typ.Map.ReflectType())
}

func TestMap_KeyType(t *testing.T) {
	m1 := vf.Map(`a`, 3, `b`, 4).Type().(dgo.MapType).KeyType()
	m2 := vf.Map(`a`, 1, `b`, 2).Type().(dgo.MapType).KeyType()
	m3 := vf.Map(`b`, 2).Type().(dgo.MapType).KeyType()
	require.Assignable(t, m1, m2)
	require.NotAssignable(t, m1, m3)
	require.Assignable(t, newtype.Enum(`a`, `b`), m1)
	require.NotAssignable(t, newtype.Enum(`b`, `c`), m1)

	require.NotAssignable(t, m2, vf.String(`b`).Type())
	require.Assignable(t, m3, vf.String(`b`).Type())

	require.Equal(t, m1, m2)
	require.NotEqual(t, m1, m3)
	require.NotEqual(t, m3, vf.String(`b`).Type())

	require.True(t, m1.HashCode() > 0)
	require.Equal(t, m1.HashCode(), m1.HashCode())
	vm := m1.Type()
	require.Instance(t, vm, m1)

	require.Equal(t, `"a"&"b"`, m1.String())
}

func TestMap_EntryType(t *testing.T) {
	vf.Map(`a`, 3).EachEntry(func(v dgo.MapEntry) {
		require.True(t, v.Frozen())
		require.Same(t, v, v.FrozenCopy())
		require.NotEqual(t, v, `a`)
		require.True(t, v.HashCode() > 0)
		require.Equal(t, v.HashCode(), v.HashCode())

		vt := v.Type()
		require.Assignable(t, vt, vt)
		require.NotAssignable(t, vt, typ.String)
		require.Instance(t, vt, v)
		require.NotInstance(t, vt, vt)
		require.Equal(t, vt, vt)
		require.NotEqual(t, vt, `a`)
		require.True(t, vt.HashCode() > 0)
		require.Equal(t, vt.HashCode(), vt.HashCode())
		require.Equal(t, `"a":3`, vt.String())

		vm := vt.Type()
		require.Instance(t, vm, vt)

		require.True(t, reflect.ValueOf(v).Type().AssignableTo(vt.ReflectType()))
	})

	m := vf.MutableMap(nil)
	m.Put(`a`, vf.MutableValues(nil, 1, 2))
	m.EachEntry(func(v dgo.MapEntry) {
		require.False(t, v.Frozen())
		require.NotSame(t, v, v.FrozenCopy())

		vt := v.Type()
		require.Equal(t, `"a":{1,2}`, vt.String())
	})
}

func TestNewMapType_max_min(t *testing.T) {
	tp := newtype.Map(2, 1)
	require.Equal(t, tp.Min(), 1)
	require.Equal(t, tp.Max(), 2)
}

func TestNewMapType_negative_min(t *testing.T) {
	tp := newtype.Map(-2, 3)
	require.Equal(t, tp.Min(), 0)
	require.Equal(t, tp.Max(), 3)
}

func TestNewMapType_negative_min_max(t *testing.T) {
	tp := newtype.Map(-2, -3)
	require.Equal(t, tp.Min(), 0)
	require.Equal(t, tp.Max(), 0)
}

func TestNewMapType_explicit_unbounded(t *testing.T) {
	tp := newtype.Map(0, -3)
	require.Equal(t, tp.Min(), 0)
	require.Equal(t, tp.Max(), 0)
}

func TestNewMapType_badOneArg(t *testing.T) {
	require.Panic(t, func() { newtype.Map(`bad`) }, `illegal argument 1`)
}

func TestNewMapType_badTwoArg(t *testing.T) {
	require.Panic(t, func() { newtype.Map(`bad`, 2) }, `illegal argument 1`)
	require.Panic(t, func() { newtype.Map(2, `bad`) }, `illegal argument 2`)
	require.Panic(t, func() { newtype.Map(typ.String, 2) }, `illegal argument 2`)
}

func TestNewMapType_badThreeArg(t *testing.T) {
	require.Panic(t, func() { newtype.Map(`bad`, typ.Integer, 2) }, `illegal argument 1`)
	require.Panic(t, func() { newtype.Map(typ.String, `bad`, 2) }, `illegal argument 2`)
	require.Panic(t, func() { newtype.Map(typ.String, 1, 2) }, `illegal argument 2`)
	require.Panic(t, func() { newtype.Map(typ.String, typ.Integer, `bad`) }, `illegal argument 3`)
}

func TestNewMapType_badFourArg(t *testing.T) {
	require.Panic(t, func() { newtype.Map(`bad`, typ.Integer, 2, 2) }, `illegal argument 1`)
	require.Panic(t, func() { newtype.Map(typ.String, `bad`, 2, 2) }, `illegal argument 2`)
	require.Panic(t, func() { newtype.Map(typ.String, typ.Integer, `bad`, 2) }, `illegal argument 3`)
	require.Panic(t, func() { newtype.Map(typ.String, typ.Integer, 2, `bad`) }, `illegal argument 4`)
}

func TestNewMapType_badArgCount(t *testing.T) {
	require.Panic(t, func() { newtype.Map(typ.String, typ.Integer, 2, 2, true) }, `illegal number of arguments`)
}

func TestMap_illegalArgument(t *testing.T) {
	require.Panic(t, func() { vf.Map('a', 23, 'b') }, `the number of arguments to Map must be 1 or an even number, got: 3`)
	require.Panic(t, func() { vf.Map(23) }, `illegal argument`)
}

func TestMap_empty(t *testing.T) {
	m := vf.Map()
	require.Equal(t, 0, m.Len())
	require.True(t, m.Frozen())
}

func TestMap_fromArray(t *testing.T) {
	m := vf.Map(vf.Values(`a`, 1, `b`, `2`))
	require.Equal(t, 2, m.Len())
	require.True(t, m.Frozen())
}

func TestMap_fromStruct(t *testing.T) {
	type TestMapFromStruct struct {
		A string
		B int
		C *string
		D *int
		E []string
	}

	c := `Charlie`
	d := 42
	s := TestMapFromStruct{A: `Alpha`, B: 32, C: &c, D: &d, E: []string{`Echo`, `Foxtrot`}}

	// Pass pointer to struct
	m := vf.Map(&s)
	require.Equal(t, 5, m.Len())
	require.False(t, m.Frozen())
	m.Freeze()
	require.True(t, m.Frozen())
	require.Equal(t, `Alpha`, m.Get(`A`))
	require.Equal(t, 32, m.Get(`B`))
	require.Equal(t, `Charlie`, m.Get(`C`))
	require.Equal(t, 42, m.Get(`D`))
	e, ok := m.Get(`E`).(dgo.Array)
	require.True(t, ok)
	require.True(t, e.Frozen())

	// Pass by value. Should give us the same result
	m2 := vf.Map(s)
	require.Equal(t, m, m2)
}

func TestMap_immutable(t *testing.T) {
	gm := map[string]int{
		`first`:  1,
		`second`: 2,
	}

	m := vf.Map(gm)
	require.Equal(t, 1, m.Get(`first`))

	gm[`first`] = 3
	require.Equal(t, 1, m.Get(`first`))

	require.Same(t, m, m.Copy(true))
}

func TestMutableMap(t *testing.T) {
	m := vf.MutableMap(map[string]string{})
	require.Equal(t, 0, m.Len())
	require.Equal(t, newtype.Map(typ.String, typ.String), m.Type())

	m = vf.MutableMap(m.Type())
	require.Equal(t, 0, m.Len())
	require.Equal(t, newtype.Map(typ.String, typ.String), m.Type())

	m = vf.MutableMap(`map[string]int`)
	require.Equal(t, 0, m.Len())
	require.Equal(t, newtype.Map(typ.String, typ.Integer), m.Type())

	m = vf.MapWithCapacity(7, vf.String(`map[string]int`))
	require.Equal(t, 0, m.Len())
	require.Equal(t, newtype.Map(typ.String, typ.Integer), m.Type())

	require.Panic(t, func() { vf.MutableMap(`[]int`) }, `does not evaluate to a map type`)
}

func TestMapFromReflected(t *testing.T) {
	m := vf.FromReflectedMap(reflect.ValueOf(map[string]string{}), false)
	require.Equal(t, 0, m.Len())
	require.Equal(t, newtype.Map(typ.String, typ.String), m.Type())
	m.Put(`hi`, `there`)
	require.Equal(t, 1, m.Len())
}

func TestMapType_KeyType(t *testing.T) {
	m := vf.Map(`hello`, `world`)
	mt := m.Type().(dgo.MapType)
	require.Instance(t, mt.KeyType(), `hello`)
	require.NotInstance(t, mt.KeyType(), `hi`)

	m = vf.Map(`hello`, `world`, 2, 2.0)
	mt = m.Type().(dgo.MapType)
	require.Assignable(t, newtype.AnyOf(typ.String, typ.Integer), mt.KeyType())
	require.Instance(t, newtype.Array(newtype.AnyOf(typ.String, typ.Integer)), m.Keys())
	require.Assignable(t, newtype.AnyOf(typ.String, typ.Float), mt.ValueType())
	require.Instance(t, newtype.Array(newtype.AnyOf(typ.String, typ.Float)), m.Values())
}

func TestMapType_ValueType(t *testing.T) {
	m := vf.Map(`hello`, `world`)
	mt := m.Type().(dgo.MapType)
	require.Instance(t, mt.ValueType(), `world`)
	require.NotInstance(t, mt.ValueType(), `earth`)
}

func TestMapNilKey(t *testing.T) {
	m := vf.Map(nil, 5)
	require.Instance(t, typ.Map, m)

	require.Nil(t, m.Get(0))
	require.Equal(t, 5, m.Get(nil))
}

func TestMap_Any(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	require.False(t, m.Any(func(e dgo.MapEntry) bool {
		return e.Key().Equals(`fourth`)
	}))
	require.True(t, m.Any(func(e dgo.MapEntry) bool {
		return e.Key().Equals(`second`)
	}))
}

func TestMap_AllKeys(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	require.False(t, m.AllKeys(func(k dgo.Value) bool {
		return len(k.String()) == 5
	}))
	require.True(t, m.AnyKey(func(k dgo.Value) bool {
		return len(k.String()) >= 5
	}))
}

func TestMap_AnyKey(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	require.False(t, m.AnyKey(func(k dgo.Value) bool {
		return k.Equals(`fourth`)
	}))
	require.True(t, m.AnyKey(func(k dgo.Value) bool {
		return k.Equals(`second`)
	}))
}

func TestMap_AnyValue(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	require.False(t, m.AnyValue(func(v dgo.Value) bool {
		return v.Equals(`four`)
	}))
	require.True(t, m.AnyValue(func(v dgo.Value) bool {
		return v.Equals(`three`)
	}))
}

func TestMap_EachKey(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	var vs []dgo.Value
	m.EachKey(func(v dgo.Value) {
		vs = append(vs, v)
	})

	require.Equal(t, 3, len(vs))
	require.Equal(t, vf.Values(`first`, `second`, `third`), vs)
}

func TestMap_EachValue(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	var vs []dgo.Value
	m.EachValue(func(v dgo.Value) {
		vs = append(vs, v)
	})

	require.Equal(t, 3, len(vs))
	require.Equal(t, vf.Values(1, 2.0, `three`), vs)
}

func TestMap_Find(t *testing.T) {
	var entry dgo.MapEntry
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	found := m.Find(func(e dgo.MapEntry) bool {
		if e.Value().Equals(2.0) {
			entry = e
			return true
		}
		return false
	})
	require.Same(t, found, entry)

	found = m.Find(func(e dgo.MapEntry) bool { return false })
	require.Nil(t, found)
}

func TestMap_Put(t *testing.T) {
	m := vf.MutableMap(nil)
	m.Put(1, `hello`)
	require.Equal(t, m, map[int]string{1: `hello`})

	m.Put(1, `hello`)
	require.Equal(t, m, map[int]string{1: `hello`})

	m = vf.Map(`first`, 1)
	require.Panic(t, func() { m.Put(`second`, 2) }, `frozen`)
}

func TestMap_PutAll(t *testing.T) {
	m := vf.MutableMap(nil)
	m.PutAll(vf.Map(
		`first`, 1,
		`second`, 2))
	require.Equal(t, m, map[string]int{
		`first`:  1,
		`second`: 2,
	})

	m.PutAll(vf.Map())
	require.Equal(t, m, map[string]int{
		`first`:  1,
		`second`: 2,
	})

	m.PutAll(vf.Map(
		`first`, 1,
		`second`, 2))
	require.Equal(t, m, map[string]int{
		`first`:  1,
		`second`: 2,
	})

	m = vf.Map(`first`, 1)
	require.Panic(t, func() { m.PutAll(vf.Map(`first`, 1)) }, `frozen`)
}

func TestMap_StringKeys(t *testing.T) {
	m := vf.Map(`a`, 1, `b`, 2)
	require.True(t, m.StringKeys())
	m = vf.Map(`a`, 1, 2, `b`)
	require.False(t, m.StringKeys())
	m = vf.MutableMap(newtype.Map(typ.String, typ.String))
	require.True(t, m.StringKeys())
}

func TestMap_Freeze_recursive(t *testing.T) {
	m := vf.MutableMap(nil)
	mr := vf.MutableMap(nil)
	mr.Put(`hello`, `world`)
	m.Put(1, mr)
	m.Freeze()
	require.True(t, mr.Frozen(), `recursive freeze not applied`)
}

func TestMap_Copy_freeze_recursive(t *testing.T) {
	m := vf.MutableMap(nil)
	mr := vf.MutableMap(nil)
	k := vf.MutableValues(nil, `the`, `key`)
	mr.Put(1.0, `world`)
	m.Put(k, mr)
	m.Put(1, `one`)
	m.Put(2, vf.Values(`x`, `y`))
	m.Put(vf.Values(`a`, `b`), vf.MutableValues(nil, `x`, `y`))

	require.True(t, vf.Array(m).All(func(v dgo.Value) bool {
		return v.(dgo.MapEntry).Frozen()
	}), `not all entries in snapshot are frozen`)

	m.EachEntry(func(e dgo.MapEntry) {
		if e.Frozen() {
			require.True(t, typ.Integer.Instance(e.Key()))
		}
	})

	mc := m.Copy(true)
	require.False(t, m.All(func(e dgo.MapEntry) bool {
		return e.Frozen()
	}), `copy affected source`)
	require.True(t, mc.All(func(e dgo.MapEntry) bool {
		return e.Frozen()
	}), `map entries are not frozen in frozen copy`)

	mcr := mc.Get(k)
	require.True(t, mcr.(dgo.Map).Frozen(), `recursive copy freeze not applied`)
	require.False(t, k.Frozen(), `recursive freeze affected key`)
	require.False(t, mr.Frozen(), `recursive freeze affected original`)

	m.Freeze()
	require.True(t, m.All(func(e dgo.MapEntry) bool {
		return e.Frozen()
	}), `map entries are not frozen after freeze`)
}

func TestMap_selfReference(t *testing.T) {
	tp := newtype.Parse(`x=map[string](string|x)`)
	d := vf.MutableMap(nil)
	d.Put(`hello`, `world`)
	d.Put(`deep`, d)
	require.Instance(t, tp, d)
	require.Equal(t, `{"hello":"world","deep":<recursive self reference to map>}`, d.String())

	t2 := newtype.Parse(`x=map[string](string|map[string](string|x))`)
	require.Assignable(t, tp, t2)
}

func TestMap_Map(t *testing.T) {
	a := vf.Map(`a`, `value a`, `b`, `value b`, `c`, `value c`)
	require.Equal(t, vf.Map(map[string]string{`a`: `the a`, `b`: `the b`, `c`: `the c`}), a.Map(func(e dgo.MapEntry) interface{} {
		return strings.Replace(e.Value().String(), `value`, `the`, 1)
	}))
	require.Equal(t, vf.Map(`a`, nil, `b`, vf.Nil, `c`, nil), a.Map(func(e dgo.MapEntry) interface{} {
		return nil
	}))
}

func TestMap_ReflectTo(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	mr := map[string]interface{}{}
	m.ReflectTo(reflect.ValueOf(&mr).Elem())
	require.Equal(t, 1, mr[`first`])
	require.Equal(t, 2.0, mr[`second`])
	require.Equal(t, `three`, mr[`third`])

	mr = map[string]interface{}{}
	mp := &mr
	m.ReflectTo(reflect.ValueOf(&mp).Elem())
	mr = *mp
	require.Equal(t, 1, mr[`first`])
	require.Equal(t, 2.0, mr[`second`])
	require.Equal(t, `three`, mr[`third`])

	var mi interface{}
	mip := &mi
	m.ReflectTo(reflect.ValueOf(mip).Elem())

	mc, ok := mi.(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, 1, mc[`first`])
	require.Equal(t, 2.0, mc[`second`])
	require.Equal(t, `three`, mc[`third`])
}

func TestMap_Remove(t *testing.T) {
	mi := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	m := mi.Copy(false)
	m.Remove(`first`)
	require.Equal(t, m, map[string]interface{}{
		`second`: 2.0,
		`third`:  `three`,
	})

	m = mi.Copy(false)
	m.Remove(`second`)
	require.Equal(t, m, map[string]interface{}{
		`first`: 1,
		`third`: `three`,
	})

	m.Remove(`first`)
	require.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})

	require.Equal(t, m.Remove(`third`), `three`)
	require.Equal(t, m.Remove(`fourth`), nil)
	require.Equal(t, m, map[string]interface{}{})

	m = vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	require.Panic(t, func() { m.Remove(`first`) }, `frozen`)
}

func TestMap_RemoveAll(t *testing.T) {
	mi := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	m := mi.Copy(false)

	m.RemoveAll(vf.Strings(`second`, `third`))
	require.Equal(t, m, map[string]interface{}{
		`first`: 1,
	})

	m = mi.Copy(false)
	m.RemoveAll(vf.Strings(`first`, `second`))
	require.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})

	m.RemoveAll(vf.Strings())
	require.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})

	m.RemoveAll(vf.Strings(`first`, `third`))
	require.Equal(t, m, map[string]interface{}{})

	require.Panic(t, func() { mi.RemoveAll(vf.Strings(`first`, `second`)) }, `frozen`)
}

func TestMap_SetType(t *testing.T) {
	m := vf.FromReflectedMap(reflect.ValueOf(map[string]interface{}{
		`first`:  1,
		`second`: 2.0,
		`third`:  `three`,
	}), false)

	mt := newtype.Map(typ.String, newtype.AnyOf(typ.Integer, typ.Float, typ.String))
	m.SetType(mt)
	require.Same(t, mt, m.Type())

	require.Panic(t, func() { m.SetType(`map[string](float|string)`) },
		`cannot be assigned`)

	require.Panic(t, func() { m.SetType(vf.String(`map[string](float|string)`)) },
		`cannot be assigned`)

	require.Panic(t, func() { m.SetType(vf.String(`float|string`)) },
		`Map.SetType: argument does not evaluate to a MapType`)

	m.SetType(nil)
	require.Assignable(t, newtype.Parse(`{"first":1,"second":2.0,"third":"three"}`), m.Type())

	m.Freeze()
	require.Panic(t, func() { m.SetType(mt) }, `frozen`)
}

func TestMap_With(t *testing.T) {
	m := vf.Map()
	m = m.With(1, `a`)

	mb := vf.Map(1, `a`)
	require.Equal(t, m, mb)

	mb = mb.With(2, `b`)
	require.Equal(t, m, map[int]string{1: `a`})
	require.Equal(t, mb, map[int]string{1: `a`, 2: `b`})

	mc := m.With(1, `a`)
	require.Same(t, m, mc)

	mc = mb.With(3, `c`)
	require.Equal(t, mc, map[int]string{1: `a`, 2: `b`, 3: `c`})
}

func TestMap_Without(t *testing.T) {
	om := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	m := om.Without(`second`)
	require.Equal(t, m, map[string]interface{}{
		`first`: 1,
		`third`: `three`,
	})

	// Original is not modified
	require.Equal(t, om, map[string]interface{}{
		`first`:  1,
		`second`: 2.0,
		`third`:  `three`,
	})

	m = m.Without(`first`)
	require.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})

	require.Same(t, m, m.Without(`first`))

	m = m.Without(`third`)
	require.Equal(t, m, map[string]interface{}{})
}

func TestMap_WithoutAll(t *testing.T) {
	om := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	m := om.WithoutAll(vf.Strings(`first`, `second`))
	require.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})

	// Original is not modified
	require.Equal(t, om, map[string]interface{}{
		`first`:  1,
		`second`: 2.0,
		`third`:  `three`,
	})

	require.Same(t, m, m.WithoutAll(vf.Values()))
	require.Same(t, m, m.WithoutAll(vf.Strings(`first`)))

	m = m.WithoutAll(vf.Strings(`first`, `third`))
	require.Equal(t, m, map[string]interface{}{})
}

func TestMap_Merge(t *testing.T) {
	m1 := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	m2 := vf.Map(
		`third`, `tres`,
		`fourth`, `cuatro`)

	require.Equal(t, m1.Merge(m2), vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `tres`,
		`fourth`, `cuatro`))

	require.Same(t, m1, m1.Merge(m1))
	require.Same(t, m1, m1.Merge(vf.Map()))
	require.Same(t, m1, vf.Map().Merge(m1))
}

func TestMap_HashCode(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	require.Equal(t, m.HashCode(), m.HashCode())
	require.NotEqual(t, 0, m.HashCode())

	m = vf.MutableMap(map[string]interface{}{})
	m.Put(`first`, 1)
	m.Put(`self`, m)

	require.NotEqual(t, 0, m.HashCode())
	require.Equal(t, m.HashCode(), m.HashCode())
}

func TestMap_Equal(t *testing.T) {
	m1 := vf.MutableMap(map[string]interface{}{})
	m1.Put(`first`, 1)
	m1.Put(`self`, m1)

	require.NotEqual(t, m1, vf.Values(`first`, `self`))

	m2 := vf.MutableMap(map[string]interface{}{})
	m2.Put(`first`, 1)
	m2.Put(`self`, m2)

	require.Equal(t, m1, m2)

	m3 := vf.MutableMap(map[string]interface{}{})
	m3.Put(`second`, 1)
	m3.Put(`self`, m3)
	require.NotEqual(t, m1, m3)
}

func TestMap_Keys(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	require.True(t, m.Keys().SameValues(vf.Values(`first`, `second`, `third`)))
}

func TestMap_Values(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	require.True(t, m.Values().SameValues(vf.Values(1, 2.0, `three`)))
}

func TestMap_String(t *testing.T) {
	require.Equal(t, `{"a":1}`, vf.Map(`a`, 1).String())
}

func TestMapEntry_Equal(t *testing.T) {
	vf.Map(`a`, 1).EachEntry(func(e dgo.MapEntry) {
		require.Equal(t, e, internal.NewMapEntry(`a`, 1))
		require.NotEqual(t, e, vf.Values(`a`, 1))
		require.NotEqual(t, e, internal.NewMapEntry(`a`, 2))
		require.NotEqual(t, e, internal.NewMapEntry(`b`, 1))
	})
}

func TestMapEntry_Frozen(t *testing.T) {
	e := internal.NewMapEntry(`a`, 1)
	require.Same(t, e, e.FrozenCopy())

	e = internal.NewMapEntry(`a`, vf.MutableValues(nil, `a`))
	require.NotSame(t, e, e.FrozenCopy())
}

func TestMapEntry_String(t *testing.T) {
	vf.Map(`a`, 1).EachEntry(func(e dgo.MapEntry) {
		require.Equal(t, `"a":1`, e.String())
	})
}

func TestMapEntry_Type(t *testing.T) {
	vf.Map(`a`, 1).EachEntry(func(e dgo.MapEntry) {
		require.Equal(t, e.Type(), internal.NewMapEntry(`a`, 1).Type())
		require.Assignable(t, e.Type(), internal.NewMapEntry(`a`, 1).Type())
		require.Instance(t, e.Type(), internal.NewMapEntry(`a`, 1))
		require.NotInstance(t, e.Type(), internal.NewMapEntry(`a`, 2))
	})
}
