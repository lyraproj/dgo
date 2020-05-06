package internal_test

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestToUint(t *testing.T) {
	_, ok := internal.ToUint(new(big.Int).SetUint64(0x8000000000000000))
	assert.True(t, ok)

	_, ok = internal.ToUint(uint(0x8000000000000000))
	assert.True(t, ok)

	_, ok = internal.ToUint(uint64(0x8000000000000000))
	assert.True(t, ok)

	_, ok = internal.ToUint(vf.Uint64(0x8000000000000000))
	assert.True(t, ok)

	bi := vf.New(typ.BigInt, vf.String(`0x8000000000000000`)).(dgo.BigInt)
	_, ok = internal.ToUint(bi)
	assert.True(t, ok)

	bi = vf.New(typ.BigInt, vf.String(`0x10000000000000000`)).(dgo.BigInt)
	_, ok = internal.ToUint(bi)
	assert.False(t, ok)

	_, ok = internal.ToUint(bi.GoBigInt())
	assert.False(t, ok)

	_, ok = internal.ToUint(-1)
	assert.False(t, ok)
}

func TestUint64_CompareTo(t *testing.T) {
	c, ok := vf.Uint64(3).CompareTo(vf.Int64(3))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Uint64(3).CompareTo(vf.Int64(-3))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(3).CompareTo(vf.Float(3))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Uint64(3).CompareTo(vf.Float(3.1))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Uint64(3).CompareTo(vf.Float(2.9))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(3).CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(3).CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	_, ok = vf.Uint64(3).CompareTo(vf.True)
	assert.False(t, ok)

	c, ok = vf.Uint64(3).CompareTo(2)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(3).CompareTo(int16(2))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(3).CompareTo(int64(2))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(3).CompareTo(uint(2))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(3).CompareTo(uint64(4))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Uint64(3).CompareTo(vf.Uint64(uint64(4)))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Uint64(math.MaxInt64).CompareTo(uint(math.MaxInt64 + 1))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Uint64(math.MaxInt64 + 1).CompareTo(uint(math.MaxInt64))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(math.MaxInt64 + 1).CompareTo(uint64(math.MaxInt64))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(math.MaxInt64 + 2).CompareTo(vf.Uint64(uint64(math.MaxInt64 + 1)))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Uint64(math.MaxInt64).CompareTo(vf.BigInt(new(big.Int).SetUint64(math.MaxInt64 + 1)))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Uint64(math.MaxUint64).CompareTo(new(big.Int).Neg(new(big.Int).SetUint64(math.MaxInt64 + 2)))
	assert.True(t, ok)
	assert.Equal(t, 1, c)
}

func TestUint64_Dec(t *testing.T) {
	assert.Equal(t, 41, vf.Uint64(42).Dec())
	assert.Equal(t, uint(math.MaxUint64), vf.Uint64(0).Dec())
}

func TestUint64_Inc(t *testing.T) {
	assert.Equal(t, 43, vf.Uint64(42).Inc())
	assert.Equal(t, 0, vf.Uint64(math.MaxUint64).Inc())
}

func TestUint64_exact(t *testing.T) {
	v := vf.Uint64(0x8000000000000000)
	tp := v.Type().(dgo.IntegerType)
	assert.Assignable(t, tp, v.Type())
	assert.NotAssignable(t, tp, v.Dec().Type())
	assert.Instance(t, tp, uint(0x8000000000000000))
	assert.NotInstance(t, tp, 2)
	assert.Equal(t, tp.Min(), uint(0x8000000000000000))
	assert.Equal(t, tp.Max(), uint(0x8000000000000000))
	assert.True(t, tp.Inclusive())
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Equal(t, `9223372036854775808`, tp.String())
	assert.Equal(t, `0x8000000000000000`, fmt.Sprintf(`%#x`, v))
	assert.Same(t, typ.Uint64, typ.Generic(tp))
	assert.Instance(t, tp.Type(), tp)
	assert.NotSame(t, tp.ReflectType(), typ.Uint64.ReflectType())
	assert.Same(t, v, v.Integer())

	assert.Panic(t, func() { v.GoInt() }, `cannot fit into an int64`)

	v = vf.Uint64(0x7000000000000000)
	assert.Equal(t, 0x7000000000000000, v.GoInt())
}

func TestUint64_ReflectTo(t *testing.T) {
	var i int
	iv := vf.Uint64(42)
	iv.ReflectTo(reflect.ValueOf(&i).Elem())
	assert.Equal(t, i, iv)

	var ip *int
	iv.ReflectTo(reflect.ValueOf(&ip).Elem())
	assert.Equal(t, i, *ip)

	var mi interface{}
	mip := &mi
	iv.ReflectTo(reflect.ValueOf(mip).Elem())
	ic, ok := mi.(uint64)
	require.True(t, ok)
	assert.Equal(t, iv, ic)

	ix8 := int8(42)
	iv = vf.Uint64(uint64(ix8))
	var i8 int8
	iv.ReflectTo(reflect.ValueOf(&i8).Elem())
	assert.Equal(t, i8, iv)

	var ip8 *int8
	iv.ReflectTo(reflect.ValueOf(&ip8).Elem())
	assert.Equal(t, ix8, *ip8)

	ix16 := int16(42)
	iv = vf.Uint64(uint64(ix16))
	var i16 int16
	iv.ReflectTo(reflect.ValueOf(&i16).Elem())
	assert.Equal(t, i16, iv)

	var ip16 *int16
	iv.ReflectTo(reflect.ValueOf(&ip16).Elem())
	assert.Equal(t, ix16, *ip16)

	ix32 := int32(42)
	iv = vf.Uint64(uint64(ix32))
	var i32 int32
	iv.ReflectTo(reflect.ValueOf(&i32).Elem())
	assert.Equal(t, i32, iv)

	var ip32 *int32
	iv.ReflectTo(reflect.ValueOf(&ip32).Elem())
	assert.Equal(t, ix32, *ip32)

	ix64 := int64(42)
	iv = vf.Uint64(uint64(ix64))
	var i64 int64
	iv.ReflectTo(reflect.ValueOf(&i64).Elem())
	assert.Equal(t, i64, iv)

	var ip64 *int64
	iv.ReflectTo(reflect.ValueOf(&ip64).Elem())
	assert.Equal(t, ix64, *ip64)

	uix := uint(42)
	iv = vf.Uint64(uint64(uix))
	var ui uint
	iv.ReflectTo(reflect.ValueOf(&ui).Elem())
	assert.Equal(t, ui, iv)

	var uip *uint
	iv.ReflectTo(reflect.ValueOf(&uip).Elem())
	assert.Equal(t, uix, *uip)

	uix8 := uint8(42)
	iv = vf.Uint64(uint64(uix8))
	var ui8 uint8
	iv.ReflectTo(reflect.ValueOf(&ui8).Elem())
	assert.Equal(t, ui8, iv)

	var uip8 *uint8
	iv.ReflectTo(reflect.ValueOf(&uip8).Elem())
	assert.Equal(t, uix8, *uip8)

	uix16 := uint16(42)
	iv = vf.Uint64(uint64(uix16))
	var ui16 uint16
	iv.ReflectTo(reflect.ValueOf(&ui16).Elem())
	assert.Equal(t, ui16, iv)

	var uip16 *uint16
	iv.ReflectTo(reflect.ValueOf(&uip16).Elem())
	assert.Equal(t, uix16, *uip16)

	uix32 := uint32(42)
	iv = vf.Uint64(uint64(uix32))
	var ui32 uint32
	iv.ReflectTo(reflect.ValueOf(&ui32).Elem())
	assert.Equal(t, ui32, iv)

	var uip32 *uint32
	iv.ReflectTo(reflect.ValueOf(&uip32).Elem())
	assert.Equal(t, uix32, *uip32)

	uix64 := uint64(42)
	iv = vf.Uint64(uix64)
	var ui64 uint64
	iv.ReflectTo(reflect.ValueOf(&ui64).Elem())
	assert.Equal(t, ui64, iv)

	var uip64 *uint64
	iv.ReflectTo(reflect.ValueOf(&uip64).Elem())
	assert.Equal(t, uix64, *uip64)
}

func TestUint64_ToInt(t *testing.T) {
	_, ok := vf.Uint64(32).ToInt()
	assert.True(t, ok)
	_, ok = vf.Uint64(0x8000000000000000).ToInt()
	assert.False(t, ok)
}

func TestUint64_ToUint(t *testing.T) {
	_, ok := vf.Uint64(32).ToUint()
	assert.True(t, ok)
}

func TestUint64_ToFloat(t *testing.T) {
	f, ok := vf.Uint64(9223372036854775808).ToFloat()
	require.True(t, ok)
	assert.Equal(t, 9.223372036854776e+18, f)
}

func TestUint64_ToBigFloat(t *testing.T) {
	assert.Equal(t, vf.BigFloat(big.NewFloat(43)), vf.Uint64(43).ToBigFloat())
}

func TestUint64_ToBigInt(t *testing.T) {
	assert.Equal(t, vf.BigInt(new(big.Int).SetUint64(0x8000000000000000)), vf.Uint64(0x8000000000000000).ToBigInt())
}
