package internal_test

import (
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestInteger(t *testing.T) {
	assert.Instance(t, typ.Integer, 3)
	assert.NotInstance(t, typ.Integer, true)
	assert.Assignable(t, typ.Integer, typ.Integer)
	assert.Assignable(t, typ.Integer, tf.Integer64(3, 5, true))
	assert.Assignable(t, typ.Integer, vf.Integer(4).Type())
	assert.Equal(t, typ.Integer, typ.Integer)
	assert.Equal(t, typ.Integer, tf.Integer64(math.MinInt64, math.MaxInt64, true))
	assert.Instance(t, typ.Integer.Type(), typ.Integer)
	assert.True(t, typ.Integer.Instance(1234))
	assert.Equal(t, typ.Integer.Min(), nil)
	assert.Equal(t, typ.Integer.Max(), nil)
	assert.True(t, typ.Integer.Inclusive())
	assert.Equal(t, `int`, typ.Integer.String())
	assert.True(t, reflect.ValueOf(int64(3)).Type().AssignableTo(typ.Integer.ReflectType()))
	assert.Equal(t, tf.Integer(vf.Integer(3), vf.Integer(5), true), tf.Integer(vf.Integer(5), vf.Integer(3), true))
	assert.Equal(t, vf.Integer(5), tf.Integer(vf.Integer(5), vf.Integer(5), true))
	assert.Same(t, typ.Integer, tf.Integer(nil, nil, true))

	assert.Panic(t, func() { tf.Integer(vf.Integer(5), vf.Integer(5), false) }, `non inclusive`)
}

func TestIntegerExact(t *testing.T) {
	tp := vf.Integer(3).Type().(dgo.IntegerType)
	assert.Instance(t, tp, 3)
	assert.NotInstance(t, tp, 2)
	assert.NotInstance(t, tp, true)
	assert.Assignable(t, tf.Integer64(3, 5, true), tp)
	assert.Assignable(t, tp, tf.Integer64(3, 3, true))
	assert.NotAssignable(t, tp, typ.Integer)
	assert.Equal(t, tp, tf.Integer64(3, 3, true))
	assert.NotEqual(t, tp, tf.Integer64(2, 5, true))
	assert.Equal(t, tp.Min(), 3)
	assert.Equal(t, tp.Max(), 3)
	assert.True(t, tp.Inclusive())
	assert.True(t, tp.Instance(3))
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Equal(t, `3`, tp.String())
	assert.Same(t, typ.Integer, typ.Generic(tp))
	assert.Instance(t, tp.Type(), tp)
	assert.Same(t, tp.ReflectType(), typ.Integer.ReflectType())
}

func TestIntegerRange(t *testing.T) {
	tp := tf.Integer64(3, 5, true)
	assert.Instance(t, tp, 3)
	assert.NotInstance(t, tp, 2)
	assert.NotInstance(t, tp, true)
	assert.Assignable(t, tp, tf.Integer64(3, 5, true))
	assert.Assignable(t, tp, tf.Integer64(4, 4, true))
	assert.Assignable(t, tp, vf.Integer(4).Type())
	assert.NotAssignable(t, tp, typ.Integer)
	assert.NotAssignable(t, tp, tf.Integer(nil, vf.Integer(5), true))
	assert.NotAssignable(t, tp, tf.Integer(vf.Integer(3), nil, true))
	assert.NotAssignable(t, tp, tf.Integer64(2, 5, true))
	assert.NotAssignable(t, tp, tf.Integer64(3, 6, true))
	assert.NotAssignable(t, tp, vf.Integer(6).Type())
	assert.Equal(t, tp, tf.Integer64(5, 3, true))
	assert.NotEqual(t, tp, tf.Integer64(2, 5, true))
	assert.NotEqual(t, tp, tf.Integer64(3, 4, true))
	assert.NotEqual(t, tp, typ.Integer)
	assert.Equal(t, tp.Min(), 3)
	assert.Equal(t, tp.Max(), 5)
	assert.Equal(t, vf.Integer(4).Type(), tf.Integer64(4, 4, true))
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Equal(t, `3..5`, tp.String())
	assert.Instance(t, tp.Type(), tp)

	tp = tf.Integer64(3, 5, false)
	assert.Instance(t, tp, 4)
	assert.Instance(t, tp, big.NewInt(4))
	assert.NotInstance(t, tp, 5)
	assert.Assignable(t, tp, tf.Integer64(3, 5, false))
	assert.NotAssignable(t, tp, tf.Integer64(3, 5, true))
	assert.Assignable(t, tf.Integer64(3, 5, true), tp)
	assert.Assignable(t, tp, tf.Integer64(3, 4, true))
	assert.Assignable(t, tf.Integer64(3, 4, true), tp)
	assert.Panic(t, func() { tf.Integer64(4, 4, false) }, `cannot have equal min and max`)
	assert.Same(t, tp.ReflectType(), typ.Integer.ReflectType())
}

func TestIntegerType_New(t *testing.T) {
	assert.Equal(t, 17, vf.New(typ.Integer, vf.Arguments(`11`, 16)))
	assert.Equal(t, 17, vf.New(typ.Integer, vf.Float(17)))
	assert.Equal(t, 0, vf.New(typ.Integer, vf.False))
	assert.Equal(t, 1, vf.New(typ.Integer, vf.True))

	now := time.Now()
	assert.Equal(t, now.Unix(), vf.New(typ.Integer, vf.Arguments(vf.Time(now))))
	assert.Panic(t, func() { vf.New(typ.Integer, vf.String(`true`)) }, `cannot be converted`)
	assert.Panic(t, func() { vf.New(vf.Integer(4).Type(), vf.Integer(5)) }, `cannot be assigned`)
	assert.Panic(t, func() { vf.New(tf.Integer64(1, 4, true), vf.Integer(5)) }, `cannot be assigned`)
}

func TestInteger_CompareTo(t *testing.T) {
	c, ok := vf.Integer(3).CompareTo(vf.Integer(3))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Integer(3).CompareTo(vf.Integer(-3))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Integer(-3).CompareTo(vf.Integer(3))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(3))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(3.1))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(2.9))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	_, ok = vf.Integer(3).CompareTo(vf.True)
	assert.False(t, ok)

	c, ok = vf.Integer(3).CompareTo(2)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(int16(2))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(int64(2))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(uint(2))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(uint64(4))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Integer(math.MaxInt64).CompareTo(uint(math.MaxInt64 + 1))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Integer(math.MaxInt64).CompareTo(uint64(math.MaxInt64 + 1))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Integer(math.MaxInt64).CompareTo(new(big.Int).SetUint64(math.MaxInt64 + 1))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Integer(math.MinInt64).CompareTo(new(big.Int).Neg(new(big.Int).SetUint64(math.MaxInt64 + 2)))
	assert.True(t, ok)
	assert.Equal(t, 1, c)
}

func TestInteger_Dec(t *testing.T) {
	assert.Equal(t, 41, vf.Integer(42).Dec())
	assert.Equal(t, math.MaxInt64, vf.Integer(math.MinInt64).Dec())
}

func TestInteger_Inc(t *testing.T) {
	assert.Equal(t, 43, vf.Integer(42).Inc())
	assert.Equal(t, math.MinInt64, vf.Integer(math.MaxInt64).Inc())
}

func TestInteger_ReflectTo(t *testing.T) {
	var i int
	iv := vf.Integer(42)
	iv.ReflectTo(reflect.ValueOf(&i).Elem())
	assert.Equal(t, i, iv)

	var ip *int
	iv.ReflectTo(reflect.ValueOf(&ip).Elem())
	assert.Equal(t, i, *ip)

	var mi interface{}
	mip := &mi
	iv.ReflectTo(reflect.ValueOf(mip).Elem())
	ic, ok := mi.(int64)
	require.True(t, ok)
	assert.Equal(t, iv, ic)

	ix8 := int8(42)
	iv = vf.Integer(int64(ix8))
	var i8 int8
	iv.ReflectTo(reflect.ValueOf(&i8).Elem())
	assert.Equal(t, i8, iv)

	var ip8 *int8
	iv.ReflectTo(reflect.ValueOf(&ip8).Elem())
	assert.Equal(t, ix8, *ip8)

	ix16 := int16(42)
	iv = vf.Integer(int64(ix16))
	var i16 int16
	iv.ReflectTo(reflect.ValueOf(&i16).Elem())
	assert.Equal(t, i16, iv)

	var ip16 *int16
	iv.ReflectTo(reflect.ValueOf(&ip16).Elem())
	assert.Equal(t, ix16, *ip16)

	ix32 := int32(42)
	iv = vf.Integer(int64(ix32))
	var i32 int32
	iv.ReflectTo(reflect.ValueOf(&i32).Elem())
	assert.Equal(t, i32, iv)

	var ip32 *int32
	iv.ReflectTo(reflect.ValueOf(&ip32).Elem())
	assert.Equal(t, ix32, *ip32)

	ix64 := int64(42)
	iv = vf.Integer(ix64)
	var i64 int64
	iv.ReflectTo(reflect.ValueOf(&i64).Elem())
	assert.Equal(t, i64, iv)

	var ip64 *int64
	iv.ReflectTo(reflect.ValueOf(&ip64).Elem())
	assert.Equal(t, ix64, *ip64)

	uix := uint(42)
	iv = vf.Integer(int64(uix))
	var ui uint
	iv.ReflectTo(reflect.ValueOf(&ui).Elem())
	assert.Equal(t, ui, iv)

	var uip *uint
	iv.ReflectTo(reflect.ValueOf(&uip).Elem())
	assert.Equal(t, uix, *uip)

	uix8 := uint8(42)
	iv = vf.Integer(int64(uix8))
	var ui8 uint8
	iv.ReflectTo(reflect.ValueOf(&ui8).Elem())
	assert.Equal(t, ui8, iv)

	var uip8 *uint8
	iv.ReflectTo(reflect.ValueOf(&uip8).Elem())
	assert.Equal(t, uix8, *uip8)

	uix16 := uint16(42)
	iv = vf.Integer(int64(uix16))
	var ui16 uint16
	iv.ReflectTo(reflect.ValueOf(&ui16).Elem())
	assert.Equal(t, ui16, iv)

	var uip16 *uint16
	iv.ReflectTo(reflect.ValueOf(&uip16).Elem())
	assert.Equal(t, uix16, *uip16)

	uix32 := uint32(42)
	iv = vf.Integer(int64(uix32))
	var ui32 uint32
	iv.ReflectTo(reflect.ValueOf(&ui32).Elem())
	assert.Equal(t, ui32, iv)

	var uip32 *uint32
	iv.ReflectTo(reflect.ValueOf(&uip32).Elem())
	assert.Equal(t, uix32, *uip32)

	uix64 := uint64(42)
	iv = vf.Integer(int64(uix64))
	var ui64 uint64
	iv.ReflectTo(reflect.ValueOf(&ui64).Elem())
	assert.Equal(t, ui64, iv)

	var uip64 *uint64
	iv.ReflectTo(reflect.ValueOf(&uip64).Elem())
	assert.Equal(t, uix64, *uip64)
}

func TestInteger_String(t *testing.T) {
	assert.Equal(t, `1234`, vf.Integer(1234).String())
	assert.Equal(t, `-4321`, vf.Integer(-4321).String())
}

func TestInteger_ToBigFloat(t *testing.T) {
	assert.Equal(t, vf.BigFloat(big.NewFloat(43)), vf.Integer(43).ToBigFloat())
}

func TestToInt(t *testing.T) {
	_, ok := internal.ToInt(new(big.Int).SetUint64(0x8000000000000000))
	assert.False(t, ok)

	_, ok = internal.ToInt(uint(0x8000000000000000))
	assert.False(t, ok)

	_, ok = internal.ToInt(uint64(0x8000000000000000))
	assert.False(t, ok)

	_, ok = internal.ToInt(new(big.Int).SetUint64(0x7000000000000000))
	assert.True(t, ok)

	_, ok = internal.ToInt(uint(0x7000000000000000))
	assert.True(t, ok)

	_, ok = internal.ToInt(uint64(0x7000000000000000))
	assert.True(t, ok)
}
