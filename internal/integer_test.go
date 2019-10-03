package internal_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestInteger(t *testing.T) {
	require.Instance(t, typ.Integer, 3)
	require.NotInstance(t, typ.Integer, true)
	require.Assignable(t, typ.Integer, typ.Integer)
	require.Assignable(t, typ.Integer, newtype.IntegerRange(3, 5, true))
	require.Assignable(t, typ.Integer, vf.Integer(4).Type())
	require.Equal(t, typ.Integer, typ.Integer)
	require.Equal(t, typ.Integer, newtype.IntegerRange(math.MinInt64, math.MaxInt64, true))
	require.Instance(t, typ.Integer.Type(), typ.Integer)
	require.True(t, typ.Integer.IsInstance(1234))
	require.Equal(t, typ.Integer.Min(), math.MinInt64)
	require.Equal(t, typ.Integer.Max(), math.MaxInt64)
	require.True(t, typ.Integer.Inclusive())

	require.Equal(t, `int`, typ.Integer.String())

	require.True(t, reflect.ValueOf(int64(3)).Type().AssignableTo(typ.Integer.ReflectType()))
}

func TestIntegerExact(t *testing.T) {
	tp := vf.Integer(3).Type().(dgo.IntegerRangeType)
	require.Instance(t, tp, 3)
	require.NotInstance(t, tp, 2)
	require.NotInstance(t, tp, true)
	require.Assignable(t, newtype.IntegerRange(3, 5, true), tp)
	require.Assignable(t, tp, newtype.IntegerRange(3, 3, true))
	require.NotAssignable(t, tp, typ.Integer)
	require.Equal(t, tp, newtype.IntegerRange(3, 3, true))
	require.NotEqual(t, tp, newtype.IntegerRange(2, 5, true))
	require.Equal(t, tp.Min(), 3)
	require.Equal(t, tp.Max(), 3)
	require.True(t, tp.Inclusive())
	require.True(t, tp.IsInstance(3))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `3`, tp.String())

	require.Same(t, typ.Integer, typ.Generic(tp))

	require.Instance(t, tp.Type(), tp)

	require.Same(t, tp.ReflectType(), typ.Integer.ReflectType())
}

func TestIntegerRange(t *testing.T) {
	tp := newtype.IntegerRange(3, 5, true)
	require.Instance(t, tp, 3)
	require.NotInstance(t, tp, 2)
	require.NotInstance(t, tp, true)
	require.Assignable(t, tp, newtype.IntegerRange(3, 5, true))
	require.Assignable(t, tp, newtype.IntegerRange(4, 4, true))
	require.Assignable(t, tp, vf.Integer(4).Type())
	require.NotAssignable(t, tp, newtype.IntegerRange(2, 5, true))
	require.NotAssignable(t, tp, newtype.IntegerRange(3, 6, true))
	require.NotAssignable(t, tp, vf.Integer(6).Type())
	require.Equal(t, tp, newtype.IntegerRange(5, 3, true))
	require.NotEqual(t, tp, newtype.IntegerRange(2, 5, true))
	require.NotEqual(t, tp, newtype.IntegerRange(3, 4, true))
	require.NotEqual(t, tp, typ.Integer)
	require.Equal(t, tp.Min(), 3)
	require.Equal(t, tp.Max(), 5)
	require.Equal(t, vf.Integer(4).Type(), newtype.IntegerRange(4, 4, true))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `3..5`, tp.String())

	require.Instance(t, tp.Type(), tp)

	tp = newtype.IntegerRange(3, 5, false)
	require.Instance(t, tp, 4)
	require.NotInstance(t, tp, 5)
	require.Assignable(t, tp, newtype.IntegerRange(3, 5, false))
	require.NotAssignable(t, tp, newtype.IntegerRange(3, 5, true))
	require.Assignable(t, newtype.IntegerRange(3, 5, true), tp)
	require.Assignable(t, tp, newtype.IntegerRange(3, 4, true))
	require.Assignable(t, newtype.IntegerRange(3, 4, true), tp)

	require.Panic(t, func() { newtype.IntegerRange(4, 4, false) }, `cannot have equal min and max`)

	require.Same(t, tp.ReflectType(), typ.Integer.ReflectType())
}

func TestInteger_CompareTo(t *testing.T) {
	c, ok := vf.Integer(3).CompareTo(vf.Integer(3))
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.Integer(3).CompareTo(vf.Integer(-3))
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Integer(-3).CompareTo(vf.Integer(3))
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(3))
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(3.1))
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(2.9))
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	_, ok = vf.Integer(3).CompareTo(vf.True)
	require.False(t, ok)
}

func TestInteger_ReflectTo(t *testing.T) {
	var i int
	iv := vf.Integer(42)
	iv.ReflectTo(reflect.ValueOf(&i).Elem())
	require.Equal(t, i, iv)

	var ip *int
	iv.ReflectTo(reflect.ValueOf(&ip).Elem())
	require.Equal(t, i, *ip)

	var mi interface{}
	mip := &mi
	iv.ReflectTo(reflect.ValueOf(mip).Elem())
	ic, ok := mi.(int64)
	require.True(t, ok)
	require.Equal(t, iv, ic)

	ix8 := int8(42)
	iv = vf.Integer(int64(ix8))
	var i8 int8
	iv.ReflectTo(reflect.ValueOf(&i8).Elem())
	require.Equal(t, i8, iv)

	var ip8 *int8
	iv.ReflectTo(reflect.ValueOf(&ip8).Elem())
	require.Equal(t, ix8, *ip8)

	ix16 := int16(42)
	iv = vf.Integer(int64(ix16))
	var i16 int16
	iv.ReflectTo(reflect.ValueOf(&i16).Elem())
	require.Equal(t, i16, iv)

	var ip16 *int16
	iv.ReflectTo(reflect.ValueOf(&ip16).Elem())
	require.Equal(t, ix16, *ip16)

	ix32 := int32(42)
	iv = vf.Integer(int64(ix32))
	var i32 int32
	iv.ReflectTo(reflect.ValueOf(&i32).Elem())
	require.Equal(t, i32, iv)

	var ip32 *int32
	iv.ReflectTo(reflect.ValueOf(&ip32).Elem())
	require.Equal(t, ix32, *ip32)

	ix64 := int64(42)
	iv = vf.Integer(ix64)
	var i64 int64
	iv.ReflectTo(reflect.ValueOf(&i64).Elem())
	require.Equal(t, i64, iv)

	var ip64 *int64
	iv.ReflectTo(reflect.ValueOf(&ip64).Elem())
	require.Equal(t, ix64, *ip64)

	uix := uint(42)
	iv = vf.Integer(int64(uix))
	var ui uint
	iv.ReflectTo(reflect.ValueOf(&ui).Elem())
	require.Equal(t, ui, iv)

	var uip *uint
	iv.ReflectTo(reflect.ValueOf(&uip).Elem())
	require.Equal(t, uix, *uip)

	uix8 := uint8(42)
	iv = vf.Integer(int64(uix8))
	var ui8 uint8
	iv.ReflectTo(reflect.ValueOf(&ui8).Elem())
	require.Equal(t, ui8, iv)

	var uip8 *uint8
	iv.ReflectTo(reflect.ValueOf(&uip8).Elem())
	require.Equal(t, uix8, *uip8)

	uix16 := uint16(42)
	iv = vf.Integer(int64(uix16))
	var ui16 uint16
	iv.ReflectTo(reflect.ValueOf(&ui16).Elem())
	require.Equal(t, ui16, iv)

	var uip16 *uint16
	iv.ReflectTo(reflect.ValueOf(&uip16).Elem())
	require.Equal(t, uix16, *uip16)

	uix32 := uint32(42)
	iv = vf.Integer(int64(uix32))
	var ui32 uint32
	iv.ReflectTo(reflect.ValueOf(&ui32).Elem())
	require.Equal(t, ui32, iv)

	var uip32 *uint32
	iv.ReflectTo(reflect.ValueOf(&uip32).Elem())
	require.Equal(t, uix32, *uip32)

	uix64 := uint64(42)
	iv = vf.Integer(int64(uix64))
	var ui64 uint64
	iv.ReflectTo(reflect.ValueOf(&ui64).Elem())
	require.Equal(t, ui64, iv)

	var uip64 *uint64
	iv.ReflectTo(reflect.ValueOf(&uip64).Elem())
	require.Equal(t, uix64, *uip64)
}

func TestInteger_String(t *testing.T) {
	require.Equal(t, `1234`, vf.Integer(1234).String())
	require.Equal(t, `-4321`, vf.Integer(-4321).String())
}
