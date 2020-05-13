package internal_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func TestMapType_New(t *testing.T) {
	mt := tf.Map(typ.String, typ.Integer)
	m := vf.Map(`first`, 1, `second`, 2)
	assert.Same(t, m, vf.New(typ.Map, m))
	assert.Same(t, m, vf.New(mt, m))
	assert.Same(t, m, vf.New(m.Type(), m))
	assert.Same(t, m, vf.New(tf.StructMapFromMap(false, vf.Map(`first`, typ.Integer, `second`, typ.Integer)), m))
	assert.Same(t, m, vf.New(mt, vf.Arguments(m)))
	assert.Equal(t, m, vf.New(mt, vf.Values(`first`, 1, `second`, 2)))
	assert.Panic(t, func() { vf.New(mt, vf.Values(`first`, 1, `second`, `two`)) }, `cannot be assigned`)
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
	assert.Assignable(t, m1, m2)
	assert.NotAssignable(t, m1, m3)
	assert.Assignable(t, tf.Integer64(1, 2, true), m1)
	assert.NotAssignable(t, tf.Integer64(2, 3, true), m1)
	assert.NotAssignable(t, m2, vf.Int64(2).Type())
	assert.Assignable(t, m4, vf.Int64(3).Type())
	assert.Equal(t, m1, m2)
	assert.NotEqual(t, m1, m3)
	assert.NotEqual(t, m4, vf.Int64(3).Type())
	assert.NotEqual(t, m1, m4)
	assert.True(t, m1.HashCode() > 0)
	assert.Equal(t, m1.HashCode(), m1.HashCode())
	vm := m1.Type()
	assert.Instance(t, vm, m1)
	assert.True(t, m1.String() == `1&2`)
}

func TestNewMapType_DefaultType(t *testing.T) {
	mt := tf.Map()
	assert.Same(t, typ.Map, mt)

	mt = tf.Map(typ.Any, typ.Any)
	assert.Same(t, typ.Map, mt)

	mt = tf.Map(typ.Any, typ.Any, 0, dgo.UnboundedSize)
	assert.Same(t, typ.Map, mt)

	m1 := vf.Map(
		`a`, 1,
		`b`, 2,
	)
	assert.Assignable(t, mt, mt)
	assert.NotAssignable(t, mt, typ.String)
	assert.Instance(t, mt, m1)
	assert.NotInstance(t, mt, `a`)
	assert.Equal(t, mt, mt)
	assert.NotEqual(t, mt, typ.String)
	assert.Equal(t, mt.KeyType(), typ.Any)
	assert.Equal(t, mt.ValueType(), typ.Any)
	assert.Equal(t, mt.Min(), 0)
	assert.Equal(t, mt.Max(), dgo.UnboundedSize)
	assert.True(t, mt.Unbounded())
	assert.True(t, mt.HashCode() > 0)
	assert.Equal(t, mt.HashCode(), mt.HashCode())

	vm := mt.Type()
	assert.Instance(t, vm, mt)
	assert.Equal(t, `map[any]any`, mt.String())
}

func TestMap_ExactType(t *testing.T) {
	m1 := vf.Map(
		`a`, 3,
		`b`, 4)
	t1 := m1.Type().(dgo.StructMapType)
	m2 := vf.Map(
		`a`, 1,
		`b`, 2)
	t2 := m2.Type().(dgo.StructMapType)
	t3 := vf.Map(`b`, 2).Type().(dgo.StructMapType)
	assert.Equal(t, 2, t1.Min())
	assert.Equal(t, 2, t1.Max())
	assert.False(t, t1.Unbounded())
	assert.Assignable(t, t1, t1)
	assert.NotAssignable(t, t1, t2)
	assert.Instance(t, t1, m1)
	assert.NotInstance(t, t1, m2)
	assert.NotInstance(t, t1, `a`)
	assert.Assignable(t, typ.Map, t1)
	assert.NotAssignable(t, t1, typ.Map)
	assert.NotAssignable(t, t1, t3)
	assert.Assignable(t, tf.Map(typ.String, typ.Integer), t1)
	assert.NotAssignable(t, tf.Map(typ.String, typ.String), t1)
	assert.NotAssignable(t, tf.Map(typ.String, typ.Integer, 3, 3), t1)
	assert.NotAssignable(t, t1, tf.Map(typ.String, typ.Integer))
	assert.NotAssignable(t, typ.String, t1)
	assert.NotEqual(t, t1, t2)
	assert.NotEqual(t, t1, t3)
	assert.NotEqual(t, t1, typ.String)
	assert.NotEqual(t, tf.Map(typ.String, typ.Integer), t1)
	assert.False(t, t1.Additional())
	assert.Equal(t, 2, t1.Len())
	assert.Equal(t, tf.Map(typ.String, typ.Integer), typ.Generic(t1))
	assert.True(t, t1.HashCode() > 0)
	assert.Equal(t, t1.HashCode(), t1.HashCode())
	vm := t1.Type()
	assert.Instance(t, vm, t1)
	assert.Equal(t, `{"a":3,"b":4}`, t1.String())
}

func TestMap_ExactType_Each(t *testing.T) {
	tp := vf.Map(`a`, 3, `b`, 4).Type().(dgo.StructMapType)
	cnt := 0
	tp.EachEntryType(func(e dgo.StructMapEntry) {
		assert.True(t, e.Required())
		assert.Assignable(t, typ.String, e.Key().(dgo.Type))
		assert.Assignable(t, typ.Integer, e.Value().(dgo.Type))
		cnt++
	})
	assert.Equal(t, 2, cnt)
}

func TestMap_ExactType_Get(t *testing.T) {
	tp := vf.Map(`a`, 3, `b`, 4).Type().(dgo.StructMapType)
	me := tp.GetEntryType(`a`)
	assert.Equal(t, tf.StructMapEntry(`a`, 3, true), me)

	me = tp.GetEntryType(vf.String(`a`).Type())
	assert.Equal(t, tf.StructMapEntry(`a`, 3, true), me)
	assert.Nil(t, tp.GetEntryType(`c`))
}

func TestMap_ExactType_Validate(t *testing.T) {
	tp := vf.Map(`a`, 3, `b`, 4).Type().(dgo.MapValidation)
	es := tp.Validate(nil, vf.Map(`a`, 3, `b`, 4))
	assert.Equal(t, 0, len(es))

	es = tp.Validate(nil, vf.Map(`a`, 2, `b`, 4))
	assert.Equal(t, 1, len(es))
}

func TestMap_ExactType_ValidateVerbose(t *testing.T) {
	tp := vf.Map(`a`, 3, `b`, 4).Type().(dgo.MapValidation)
	out := util.NewIndenter(``)
	assert.False(t, tp.ValidateVerbose(vf.Values(1, 2), out))
	assert.Equal(t, `value is not a Map`, out.String())
}

func TestMap_SizedType(t *testing.T) {
	mt := tf.Map(typ.String, typ.Integer)
	assert.Equal(t, mt, mt)
	assert.Equal(t, mt, tf.Map(typ.String, typ.Integer))
	assert.NotEqual(t, mt, tf.Map(typ.String, typ.Float))
	assert.NotEqual(t, mt, tf.Array(typ.String))

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
	assert.Assignable(t, mt, mt)
	assert.NotAssignable(t, mt, typ.Any)
	assert.Instance(t, mt, m1)
	assert.NotInstance(t, mt, m2)
	assert.NotInstance(t, mt, `a`)

	mtz := tf.Map(typ.String, typ.Integer, 3, 3)
	assert.Instance(t, mtz, m3)
	assert.NotInstance(t, mtz, m1)
	assert.Assignable(t, mt, mtz)
	assert.NotAssignable(t, mtz, mt)
	assert.NotEqual(t, mt, mtz)
	assert.NotEqual(t, mt, typ.Any)

	mta := tf.Map(typ.Any, typ.Any, 3, 3)
	assert.NotInstance(t, mta, m1)
	assert.Instance(t, mta, m3)

	mtva := tf.Map(typ.String, typ.Any)
	assert.Instance(t, mtva, m1)
	assert.Instance(t, mtva, m2)
	assert.Instance(t, mtva, m3)

	mtka := tf.Map(typ.Any, typ.Integer)
	assert.Instance(t, mtka, m1)
	assert.NotInstance(t, mtka, m2)
	assert.Instance(t, mtka, m3)
	assert.True(t, mt.HashCode() > 0)
	assert.Equal(t, mt.HashCode(), mt.HashCode())
	assert.NotEqual(t, mt.HashCode(), mtz.HashCode())
	vm := mt.Type()
	assert.Instance(t, vm, mt)
	assert.Equal(t, `map[string]int`, mt.String())
	assert.Equal(t, `map[string,3,3]int`, mtz.String())
	assert.Equal(t, mta.ReflectType(), typ.Map.ReflectType())
}

func TestMap_KeyType(t *testing.T) {
	m1 := vf.Map(`a`, 3, `b`, 4).Type().(dgo.MapType).KeyType()
	m2 := vf.Map(`a`, 1, `b`, 2).Type().(dgo.MapType).KeyType()
	m3 := vf.Map(`b`, 2).Type().(dgo.MapType).KeyType()
	assert.Assignable(t, m1, m2)
	assert.NotAssignable(t, m1, m3)
	assert.Assignable(t, tf.Enum(`a`, `b`), m1)
	assert.NotAssignable(t, tf.Enum(`b`, `c`), m1)
	assert.NotAssignable(t, m2, vf.String(`b`).Type())
	assert.Assignable(t, m3, vf.String(`b`).Type())
	assert.Equal(t, m1, m2)
	assert.NotEqual(t, m1, m3)
	assert.NotEqual(t, m3, vf.String(`b`).Type())
	assert.True(t, m1.HashCode() > 0)
	assert.Equal(t, m1.HashCode(), m1.HashCode())
	vm := m1.Type()
	assert.Instance(t, vm, m1)
	assert.Equal(t, `"a"&"b"`, m1.String())
}

func TestMap_EntryType(t *testing.T) {
	vf.Map(`a`, 3).EachEntry(func(v dgo.MapEntry) {
		assert.True(t, v.Frozen())
		assert.Same(t, v, v.FrozenCopy())
		assert.NotEqual(t, v, `a`)
		assert.True(t, v.HashCode() > 0)
		assert.Equal(t, v.HashCode(), v.HashCode())

		vt := v.Type()
		assert.Assignable(t, vt, vt)
		assert.NotAssignable(t, vt, typ.String)
		assert.Instance(t, vt, v)
		assert.Instance(t, vt, vt)
		assert.Equal(t, vt, vt)
		assert.NotEqual(t, vt, `a`)
		assert.True(t, vt.HashCode() > 0)
		assert.Equal(t, vt.HashCode(), vt.HashCode())
		assert.Equal(t, `"a":3`, vt.String())

		vm := vt.Type()
		assert.Instance(t, vm, vt)
		assert.Same(t, typ.Any, typ.Generic(vt))
		assert.True(t, reflect.ValueOf(v).Type().AssignableTo(vt.ReflectType()))
	})

	m := vf.MutableMap()
	m.Put(`a`, vf.MutableValues(1, 2))
	m.EachEntry(func(v dgo.MapEntry) {
		assert.False(t, v.Frozen())
		assert.NotSame(t, v, v.FrozenCopy())

		vt := v.Type()
		assert.Equal(t, `"a":{1,2}`, vt.String())
	})
	m = m.FrozenCopy().(dgo.Map)
	m.EachEntry(func(v dgo.MapEntry) {
		assert.True(t, v.Frozen())
		assert.NotSame(t, v, v.ThawedCopy())
	})
	m = m.ThawedCopy().(dgo.Map)
	m.EachEntry(func(v dgo.MapEntry) {
		assert.False(t, v.Frozen())
		assert.NotSame(t, v, v.ThawedCopy())
	})

	m = vf.MutableMap()
	m.Put(`a`, 1)
	m.Put(`b`, `string`)
	m.EachEntry(func(v dgo.MapEntry) {
		assert.True(t, v.Frozen())
		assert.Same(t, v, v.ThawedCopy())
	})
}

func TestNewMapType_max_min(t *testing.T) {
	tp := tf.Map(2, 1)
	assert.Equal(t, tp.Min(), 1)
	assert.Equal(t, tp.Max(), 2)
}

func TestNewMapType_negative_min(t *testing.T) {
	tp := tf.Map(-2, 3)
	assert.Equal(t, tp.Min(), 0)
	assert.Equal(t, tp.Max(), 3)
}

func TestNewMapType_negative_min_max(t *testing.T) {
	tp := tf.Map(-2, -3)
	assert.Equal(t, tp.Min(), 0)
	assert.Equal(t, tp.Max(), 0)
}

func TestNewMapType_explicit_unbounded(t *testing.T) {
	tp := tf.Map(0, -3)
	assert.Equal(t, tp.Min(), 0)
	assert.Equal(t, tp.Max(), 0)
}

func TestNewMapType_badOneArg(t *testing.T) {
	assert.Panic(t, func() { tf.Map(`bad`) }, `illegal argument`)
}

func TestNewMapType_badTwoArg(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	n := testNamed(0)
	assert.Panic(t, func() { tf.Map(n, typ.String) }, `illegal argument 1`)
	assert.Panic(t, func() { tf.Map(1, typ.String) }, `illegal argument 2`)
	assert.Panic(t, func() { tf.Map(typ.String, n) }, `illegal argument 2`)
}

func TestNewMapType_badThreeArg(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	n := testNamed(0)
	assert.Panic(t, func() { tf.Map(n, typ.Integer, 2) }, `illegal argument 1`)
	assert.Panic(t, func() { tf.Map(typ.String, n, `bad`) }, `illegal argument 2`)
	assert.Panic(t, func() { tf.Map(typ.String, typ.Integer, `bad`) }, `illegal argument 3`)
}

func TestNewMapType_badFourArg(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	n := testNamed(0)
	assert.Panic(t, func() { tf.Map(n, typ.Integer, 2, 2) }, `illegal argument 1`)
	assert.Panic(t, func() { tf.Map(typ.String, n, 2, 2) }, `illegal argument 2`)
	assert.Panic(t, func() { tf.Map(typ.String, typ.Integer, `bad`, 2) }, `illegal argument 3`)
	assert.Panic(t, func() { tf.Map(typ.String, typ.Integer, 2, `bad`) }, `illegal argument 4`)
}

func TestNewMapType_badArgCount(t *testing.T) {
	assert.Panic(t, func() { tf.Map(typ.String, typ.Integer, 2, 2, true) }, `illegal number of arguments`)
}

func TestMap_illegalArgument(t *testing.T) {
	assert.Panic(t, func() { vf.Map('a', 23, 'b') }, `the number of arguments to Map must be 1 or an even number, got: 3`)
	assert.Panic(t, func() { vf.Map(23) }, `illegal argument`)
}

func TestMap_empty(t *testing.T) {
	m := vf.Map()
	assert.Equal(t, 0, m.Len())
	assert.True(t, m.Frozen())
}

func TestMap_fromArray(t *testing.T) {
	m := vf.Map(vf.Values(`a`, 1, `b`, `2`))
	assert.Equal(t, 2, m.Len())
	assert.True(t, m.Frozen())
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
	assert.Equal(t, 5, m.Len())
	assert.True(t, m.Frozen())
	m = m.Copy(true)
	assert.True(t, m.Frozen())
	assert.Equal(t, `Alpha`, m.Get(`A`))
	assert.Equal(t, 32, m.Get(`B`))
	assert.Equal(t, c, m.Get(`C`))
	assert.Equal(t, 42, m.Get(`D`))
	e, ok := m.Get(`E`).(dgo.Array)
	require.True(t, ok)
	assert.True(t, e.Frozen())

	// Pass by value. Should give us the same result
	m2 := vf.Map(s)
	assert.Equal(t, m, m2)
}

func TestMap_immutable(t *testing.T) {
	gm := map[string]int{
		`first`:  1,
		`second`: 2,
	}

	m := vf.Map(gm)
	assert.Equal(t, 1, m.Get(`first`))

	gm[`first`] = 3
	assert.Equal(t, 1, m.Get(`first`))
	assert.Same(t, m, m.Copy(true))
}

func TestMapFromReflected(t *testing.T) {
	m := vf.FromReflectedMap(reflect.ValueOf(map[string]string{}), false)
	assert.Equal(t, 0, m.Len())
	assert.Equal(t, tf.ParseType(`{}`), m.Type())
	assert.Equal(t, tf.ParseType(`map[any]any`), typ.Generic(m.Type()))
	m = vf.FromReflectedMap(reflect.ValueOf(map[string]string{`foo`: `bar`}), false)
	assert.Equal(t, 1, m.Len())
	assert.Equal(t, tf.ParseType(`{foo: "bar"}`), m.Type())
	assert.Equal(t, tf.ParseType(`map[string]string`), typ.Generic(m.Type()))
	m.Put(`hi`, `there`)
	assert.Equal(t, 2, m.Len())
}

func TestMapType_KeyType(t *testing.T) {
	m := vf.Map(`hello`, `world`)
	mt := m.Type().(dgo.MapType)
	assert.Instance(t, mt.KeyType(), `hello`)
	assert.NotInstance(t, mt.KeyType(), `hi`)

	m = vf.Map(`hello`, `world`, 2, 2.0)
	mt = m.Type().(dgo.MapType)
	assert.Assignable(t, tf.AnyOf(typ.String, typ.Integer), mt.KeyType())
	assert.Instance(t, tf.Array(tf.AnyOf(typ.String, typ.Integer)), m.Keys())
	assert.Assignable(t, tf.AnyOf(typ.String, typ.Float), mt.ValueType())
	assert.Instance(t, tf.Array(tf.AnyOf(typ.String, typ.Float)), m.Values())
}

func TestMapType_ValueType(t *testing.T) {
	m := vf.Map(`hello`, `world`)
	mt := m.Type().(dgo.MapType)
	assert.Instance(t, mt.ValueType(), `world`)
	assert.NotInstance(t, mt.ValueType(), `earth`)
}

func TestMapNilKey(t *testing.T) {
	m := vf.Map(nil, 5)
	assert.Instance(t, typ.Map, m)
	assert.Nil(t, m.Get(0))
	assert.Equal(t, 5, m.Get(nil))
}

func TestMap_Any(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	assert.False(t, m.Any(func(e dgo.MapEntry) bool {
		return e.Key().Equals(`fourth`)
	}))
	assert.True(t, m.Any(func(e dgo.MapEntry) bool {
		return e.Key().Equals(`second`)
	}))
}

func TestMap_AllKeys(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	assert.False(t, m.AllKeys(func(k dgo.Value) bool {
		return len(k.String()) == 5
	}))
	assert.True(t, m.AnyKey(func(k dgo.Value) bool {
		return len(k.String()) >= 5
	}))
}

func TestMap_AnyKey(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	assert.False(t, m.AnyKey(func(k dgo.Value) bool {
		return k.Equals(`fourth`)
	}))
	assert.True(t, m.AnyKey(func(k dgo.Value) bool {
		return k.Equals(`second`)
	}))
}

func TestMap_AnyValue(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	assert.False(t, m.AnyValue(func(v dgo.Value) bool {
		return v.Equals(`four`)
	}))
	assert.True(t, m.AnyValue(func(v dgo.Value) bool {
		return v.Equals(`three`)
	}))
}

func TestMap_AppendAt(t *testing.T) {
	m := vf.MutableMap(
		`first`, 1,
		`second`, vf.MutableValues(2))

	assert.Equal(t, vf.Values(10), m.AppendAt(`first`, 10))
	assert.Equal(t, vf.Values(2, 3), m.AppendAt(`second`, 3))
}

func TestMap_AppendAt_immutable(t *testing.T) {
	m := vf.Map(`first`, 1)
	assert.Panic(t, func() { m.AppendAt(`first`, 10) }, `AppendAt called on a frozen Map`)
}

func TestMap_ContainsKey(t *testing.T) {
	assert.True(t, vf.Map(`a`, `the a`).ContainsKey(`a`))
	assert.False(t, vf.Map(`a`, `the a`).ContainsKey(`b`))
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
	assert.Equal(t, 3, len(vs))
	assert.Equal(t, vf.Values(`first`, `second`, `third`), vs)
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
	assert.Equal(t, 3, len(vs))
	assert.Equal(t, vf.Values(1, 2.0, `three`), vs)
}

func TestMap_Format(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	assert.Equal(t, `map[string]any{"first":1, "second":2, "third":"three"}`, fmt.Sprintf("%#v", m))
}

func TestMap_Format_nativeKey(t *testing.T) {
	m := vf.Map(
		&testStruct{A: `a`}, 1,
		&testStruct{A: `b`}, 2)
	assert.Equal(t, `map[native["*internal_test.testStruct"]]int{&internal_test.testStruct{A:"a"}:1, &internal_test.testStruct{A:"b"}:2}`, fmt.Sprintf("%#v", m))
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
	assert.Same(t, found, entry)

	found = m.Find(func(e dgo.MapEntry) bool { return false })
	assert.Nil(t, found)
}

func TestMap_Put(t *testing.T) {
	m := vf.MutableMap(vf.Values(1, `hello`))
	assert.Equal(t, m, map[int]string{1: `hello`})

	m.Put(1, `hello`)
	assert.Equal(t, m, map[int]string{1: `hello`})

	m = vf.Map(`first`, 1)
	assert.Panic(t, func() { m.Put(`second`, 2) }, `frozen`)
}

func TestMap_PutAll(t *testing.T) {
	m := vf.MutableMap(
		`first`, 1,
		`second`, 2)
	assert.Equal(t, m, map[string]int{
		`first`:  1,
		`second`: 2,
	})

	m.PutAll(vf.Map())
	assert.Equal(t, m, map[string]int{
		`first`:  1,
		`second`: 2,
	})

	m.PutAll(vf.Map(
		`first`, 1,
		`second`, 2))
	assert.Equal(t, m, map[string]int{
		`first`:  1,
		`second`: 2,
	})

	m = vf.Map(`first`, 1)
	assert.Panic(t, func() { m.PutAll(vf.Map(`first`, 1)) }, `frozen`)
}

func TestMap_StringKeys(t *testing.T) {
	m := vf.Map(`a`, 1, `b`, 2)
	assert.True(t, m.StringKeys())
	m = vf.Map(`a`, 1, 2, `b`)
	assert.False(t, m.StringKeys())
	m = vf.MutableMap()
	assert.True(t, m.StringKeys())
}

func TestHashMap_ThawedCopy(t *testing.T) {
	mr := vf.MutableMap()
	mr.Put(`hello`, `world`)
	bs := vf.Binary([]byte{1, 2, 3}, false)
	m := vf.Map(1, mr, 2, bs, 3, `string`)
	assert.True(t, m.Get(1).(dgo.Map).Frozen())
	assert.True(t, m.Get(2).(dgo.Binary).Frozen())

	mt := m.ThawedCopy().(dgo.Map)
	assert.True(t, m.Get(1).(dgo.Map).Frozen())
	assert.True(t, m.Get(2).(dgo.Binary).Frozen())

	assert.False(t, mt.Get(1).(dgo.Map).Frozen())
	assert.False(t, mt.Get(2).(dgo.Binary).Frozen())
}

func TestMap_Copy_freeze_recursive(t *testing.T) {
	m := vf.MutableMap()
	mr := vf.MutableMap()
	k := vf.MutableValues(`the`, `key`)
	mr.Put(1.0, `world`)
	m.Put(k, mr)
	m.Put(1, `one`)
	m.Put(2, vf.Values(`x`, `y`))
	m.Put(vf.Values(`a`, `b`), vf.MutableValues(`x`, `y`))
	assert.True(t, vf.Array(m).All(func(v dgo.Value) bool {
		return v.(dgo.MapEntry).Frozen()
	}), `not all entries in snapshot are frozen`)

	m.EachEntry(func(e dgo.MapEntry) {
		if e.Frozen() {
			assert.True(t, typ.Integer.Instance(e.Key()))
		}
	})

	mc := m.Copy(true)
	assert.False(t, m.All(func(e dgo.MapEntry) bool {
		return e.Frozen()
	}), `copy affected source`)
	assert.True(t, mc.All(func(e dgo.MapEntry) bool {
		return e.Frozen()
	}), `map entries are not frozen in frozen copy`)

	mcr := mc.Get(k)
	assert.True(t, mcr.(dgo.Map).Frozen(), `recursive copy freeze not applied`)
	assert.False(t, k.Frozen(), `recursive freeze affected key`)
	assert.False(t, mr.Frozen(), `recursive freeze affected original`)
}

func TestMap_selfReference(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`x=map[string](string|x)`)
	d := vf.MutableMap()
	d.Put(`hello`, `world`)
	d.Put(`deep`, d)
	assert.Instance(t, tp, d)
	assert.Equal(t, `{"hello":"world","deep":<recursive self reference to map type>}`, d.String())

	internal.ResetDefaultAliases()
	t2 := tf.ParseType(`x=map[string](string|map[string](string|x))`)
	assert.Assignable(t, tp, t2)
}

func TestMap_Map(t *testing.T) {
	a := vf.Map(`a`, `value a`, `b`, `value b`, `c`, `value c`)
	assert.Equal(t,
		vf.Map(map[string]string{`a`: `the a`, `b`: `the b`, `c`: `the c`}),
		a.Map(func(e dgo.MapEntry) interface{} {
			return strings.Replace(e.Value().(dgo.String).GoString(), `value`, `the`, 1)
		}))
	assert.Equal(t, vf.Map(`a`, nil, `b`, vf.Nil, `c`, nil), a.Map(func(e dgo.MapEntry) interface{} {
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
	assert.Equal(t, 1, mr[`first`])
	assert.Equal(t, 2.0, mr[`second`])
	assert.Equal(t, `three`, mr[`third`])

	mr = map[string]interface{}{}
	mp := &mr
	m.ReflectTo(reflect.ValueOf(&mp).Elem())
	mr = *mp
	assert.Equal(t, 1, mr[`first`])
	assert.Equal(t, 2.0, mr[`second`])
	assert.Equal(t, `three`, mr[`third`])

	var mi interface{}
	mip := &mi
	m.ReflectTo(reflect.ValueOf(mip).Elem())

	mc, ok := mi.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, 1, mc[`first`])
	assert.Equal(t, 2.0, mc[`second`])
	assert.Equal(t, `three`, mc[`third`])
}

func TestMap_Remove(t *testing.T) {
	mi := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	m := mi.Copy(false)
	m.Remove(`first`)
	assert.Equal(t, m, map[string]interface{}{
		`second`: 2.0,
		`third`:  `three`,
	})

	m = mi.Copy(false)
	m.Remove(`second`)
	assert.Equal(t, m, map[string]interface{}{
		`first`: 1,
		`third`: `three`,
	})

	m.Remove(`first`)
	assert.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})
	assert.Equal(t, m.Remove(`third`), `three`)
	assert.Equal(t, m.Remove(`fourth`), nil)
	assert.Equal(t, m, map[string]interface{}{})

	m = vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	assert.Panic(t, func() { m.Remove(`first`) }, `frozen`)
}

func TestMap_RemoveAll(t *testing.T) {
	mi := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	m := mi.Copy(false)

	m.RemoveAll(vf.Strings(`second`, `third`))
	assert.Equal(t, m, map[string]interface{}{
		`first`: 1,
	})

	m = mi.Copy(false)
	m.RemoveAll(vf.Strings(`first`, `second`))
	assert.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})

	m.RemoveAll(vf.Strings())
	assert.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})

	m.RemoveAll(vf.Strings(`first`, `third`))
	assert.Equal(t, m, map[string]interface{}{})
	assert.Panic(t, func() { mi.RemoveAll(vf.Strings(`first`, `second`)) }, `frozen`)
}

func TestMap_With(t *testing.T) {
	m := vf.Map()
	m = m.With(1, `a`)

	mb := vf.Map(1, `a`)
	assert.Equal(t, m, mb)

	mb = mb.With(2, `b`)
	assert.Equal(t, m, map[int]string{1: `a`})
	assert.Equal(t, mb, map[int]string{1: `a`, 2: `b`})

	mc := m.With(1, `a`)
	assert.Same(t, m, mc)

	mc = mb.With(3, `c`)
	assert.Equal(t, mc, map[int]string{1: `a`, 2: `b`, 3: `c`})
}

func TestMap_Without(t *testing.T) {
	om := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	m := om.Without(`second`)
	assert.Equal(t, m, map[string]interface{}{
		`first`: 1,
		`third`: `three`,
	})

	// Original is not modified
	assert.Equal(t, om, map[string]interface{}{
		`first`:  1,
		`second`: 2.0,
		`third`:  `three`,
	})

	m = m.Without(`first`)
	assert.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})
	assert.Same(t, m, m.Without(`first`))

	m = m.Without(`third`)
	assert.Equal(t, m, map[string]interface{}{})
}

func TestMap_WithoutAll(t *testing.T) {
	om := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	m := om.WithoutAll(vf.Strings(`first`, `second`))
	assert.Equal(t, m, map[string]interface{}{
		`third`: `three`,
	})

	// Original is not modified
	assert.Equal(t, om, map[string]interface{}{
		`first`:  1,
		`second`: 2.0,
		`third`:  `three`,
	})
	assert.Same(t, m, m.WithoutAll(vf.Values()))
	assert.Same(t, m, m.WithoutAll(vf.Strings(`first`)))

	m = m.WithoutAll(vf.Strings(`first`, `third`))
	assert.Equal(t, m, map[string]interface{}{})
}

func TestMap_Merge(t *testing.T) {
	m1 := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	m2 := vf.Map(
		`third`, `tres`,
		`fourth`, `cuatro`)
	assert.Equal(t, m1.Merge(m2), vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `tres`,
		`fourth`, `cuatro`))
	assert.Same(t, m1, m1.Merge(m1))
	assert.Same(t, m1, m1.Merge(vf.Map()))
	assert.Same(t, m1, vf.Map().Merge(m1))
}

func TestMap_HashCode(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	assert.Equal(t, m.HashCode(), m.HashCode())
	assert.NotEqual(t, 0, m.HashCode())

	m2 := vf.Map(
		`second`, 2.0,
		`first`, 1,
		`third`, `three`)
	assert.Equal(t, m.HashCode(), m2.HashCode())

	// Self containing map
	m = vf.MutableMap()
	m.Put(`first`, 1)
	m.Put(`self`, m)
	assert.NotEqual(t, 0, m.HashCode())
	assert.Equal(t, m.HashCode(), m.HashCode())
}

func TestMap_Equal(t *testing.T) {
	m1 := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)

	m2 := vf.Map(
		`second`, 2.0,
		`first`, 1,
		`third`, `three`)
	assert.Equal(t, m1, m2)

	m1 = vf.MutableMap()
	m1.Put(`first`, 1)
	m1.Put(`self`, m1)
	assert.Equal(t, m1.Keys(), vf.Values(`first`, `self`))
	assert.NotEqual(t, m1, vf.Values(`first`, `self`))

	m2 = vf.MutableMap()
	m2.Put(`first`, 1)
	m2.Put(`self`, m2)
	assert.Equal(t, m1, m2)

	m3 := vf.MutableMap()
	m3.Put(`second`, 1)
	m3.Put(`self`, m3)
	assert.NotEqual(t, m1, m3)
}

func TestMap_Keys(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	assert.True(t, m.Keys().SameValues(vf.Values(`first`, `second`, `third`)))
}

func TestMap_Values(t *testing.T) {
	m := vf.Map(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	a := m.Values()
	assert.True(t, a.Frozen())
	assert.Equal(t, a, vf.Values(1, 2.0, `three`))
}

func TestMap_Values_mutable(t *testing.T) {
	m := vf.MutableMap(
		`first`, 1,
		`second`, 2.0,
		`third`, `three`)
	a := m.Values()
	assert.False(t, a.Frozen())
	assert.Equal(t, a, vf.Values(1, 2.0, `three`))
}

func TestMap_String(t *testing.T) {
	assert.Equal(t, `{"a":1}`, vf.Map(`a`, 1).String())
}

func TestMap_Resolve(t *testing.T) {
	n := vf.String(`b`)
	am := tf.BuiltInAliases().Collect(func(a dgo.AliasAdder) {
		a.Add(tf.Integer64(0, 255, true), n)
	})
	assert.Equal(t, n, am.GetName(tf.Integer64(0, 255, true)))
}

func TestMapEntry_Equal(t *testing.T) {
	vf.Map(`a`, 1).EachEntry(func(e dgo.MapEntry) {
		assert.Equal(t, e, internal.NewMapEntry(`a`, 1))
		assert.NotEqual(t, e, vf.Values(`a`, 1))
		assert.NotEqual(t, e, internal.NewMapEntry(`a`, 2))
		assert.NotEqual(t, e, internal.NewMapEntry(`b`, 1))
	})
}

func TestMapEntry_Frozen(t *testing.T) {
	e := internal.NewMapEntry(`a`, 1)
	assert.Same(t, e, e.FrozenCopy())

	e = internal.NewMapEntry(`a`, vf.MutableValues(`a`))
	assert.NotSame(t, e, e.FrozenCopy())
}

func TestMapEntry_Thawed(t *testing.T) {
	e := internal.NewMapEntry(`a`, vf.MutableValues(1))
	assert.NotSame(t, e, e.ThawedCopy())

	e = e.FrozenCopy().(dgo.MapEntry)
	assert.NotSame(t, e, e.ThawedCopy())

	e = internal.NewMapEntry(`a`, vf.MutableValues(`a`))
	c := e.ThawedCopy().(dgo.MapEntry)
	c.Value().(dgo.Array).Set(0, `b`)
	assert.Equal(t, vf.Values(`a`), e.Value())
	assert.Equal(t, vf.Values(`b`), c.Value())

	e = internal.NewMapEntry(`a`, `b`)
	assert.Same(t, e, e.ThawedCopy())
}

func TestMapEntry_String(t *testing.T) {
	vf.Map(`a`, 1).EachEntry(func(e dgo.MapEntry) {
		assert.Equal(t, `"a":1`, e.String())
	})
}

func TestMapEntry_Type(t *testing.T) {
	vf.Map(`a`, 1).EachEntry(func(e dgo.MapEntry) {
		assert.Equal(t, e.Type(), internal.NewMapEntry(`a`, 1).Type())
		assert.Assignable(t, e.Type(), internal.NewMapEntry(`a`, 1).Type())
		assert.Instance(t, e.Type(), internal.NewMapEntry(`a`, 1))
		assert.NotInstance(t, e.Type(), internal.NewMapEntry(`a`, 2))
	})
}
