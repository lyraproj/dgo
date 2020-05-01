package internal_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/test/require"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func newBigFloat(t *testing.T, s string) dgo.BigFloat {
	t.Helper()
	bf, ok := vf.New(typ.BigFloat, vf.String(s)).(dgo.BigFloat)
	require.True(t, ok)
	return bf
}

func newBigFloatRange(t *testing.T, min string, max string, inclusive bool) dgo.FloatType {
	t.Helper()
	var minBf, maxBf dgo.BigFloat
	if min != "" {
		minBf = newBigFloat(t, min)
	}
	if max != "" {
		maxBf = newBigFloat(t, max)
	}
	return tf.Float(minBf, maxBf, inclusive)
}

func TestBigFloatType_New(t *testing.T) {
	ex, _ := new(big.Float).SetString(`3.14`)
	assert.True(t, ex.Cmp(newBigFloat(t, `3.14`).GoBigFloat()) == 0)
	assert.True(t, ex.Cmp(
		vf.New(newBigFloatRange(t, `2`, `4`, true), vf.String(`3.14`)).(dgo.BigFloat).GoBigFloat()) == 0)
	assert.True(t, ex.Cmp(
		vf.New(newBigFloat(t, `3.14`).Type(), vf.String(`3.14`)).(dgo.BigFloat).GoBigFloat()) == 0)

	assert.Equal(t, new(big.Float).SetInt64(int64(123)), vf.New(typ.BigFloat, vf.Integer(123)))
	assert.Equal(t, big.NewFloat(123), vf.New(typ.BigFloat, vf.Float(123)))
	assert.Equal(t, big.NewFloat(0), vf.New(typ.BigFloat, vf.Boolean(false)))
	assert.Equal(t, big.NewFloat(1), vf.New(typ.BigFloat, vf.Boolean(true)))

	assert.Panic(t, func() {
		vf.New(typ.BigFloat, vf.Binary([]byte{1, 2}, false))
	}, `cannot be converted`)

	assert.Panic(t, func() {
		vf.New(newBigFloatRange(t, `4`, `10`, true), vf.Float(3.14))
	}, `cannot be assigned`)
}

func TestBigFloatType_Range(t *testing.T) {
	tp := newBigFloatRange(t, `123.45`, `678.90`, true)
	assert.Assignable(t, tp, newBigFloatRange(t, `123.45`, `678.90`, true))
	assert.Assignable(t, tp, newBigFloatRange(t, `234.46`, `678.90`, true))
	assert.Assignable(t, tp, newBigFloatRange(t, `123.45`, `567.89`, true))
	assert.Assignable(t, tp, newBigFloatRange(t, `234.56`, `567.89`, true))
	assert.Assignable(t, tp, newBigFloatRange(t, `123.45`, `678.90`, false))
	assert.NotAssignable(t, newBigFloatRange(t, `123.45`, `678.90`, false), tp)

	assert.Assignable(t, newBigFloat(t, `345.67`).Type(), newBigFloat(t, `345.67`).Type())
	assert.Assignable(t, typ.Generic(newBigFloat(t, `345.67`).Type()), tp)
	assert.NotAssignable(t, newBigFloat(t, `345.67`).Type(), tp)

	assert.Instance(t, tp, newBigFloat(t, `345.67`))
	assert.Instance(t, newBigFloat(t, `345.67`).Type(), newBigFloat(t, `345.67`))
	assert.NotInstance(t, tp, newBigFloat(t, `120.2`))

	assert.Equal(t, `big 123.45..678.9`, tp.String())
}

func TestBigFloatType_ReflectType(t *testing.T) {
	tp := typ.BigFloat
	assert.Equal(t, tp.ReflectType(), reflect.TypeOf(new(big.Float)))
	tp = newBigFloatRange(t, `123.45`, `678.90`, true)
	assert.Equal(t, tp.ReflectType(), reflect.TypeOf(new(big.Float)))
	tp = newBigFloat(t, `123.45`).(dgo.FloatType)
	assert.Equal(t, tp.ReflectType(), reflect.TypeOf(new(big.Float)))
}

func TestBigFloatType_TypeIdentifier(t *testing.T) {
	tp := typ.BigFloat
	assert.Equal(t, tp.TypeIdentifier(), dgo.TiBigFloat)
	tp = newBigFloatRange(t, `123.45`, `678.90`, true)
	assert.Equal(t, tp.TypeIdentifier(), dgo.TiBigFloatRange)
	tp = newBigFloat(t, `123.45`).(dgo.FloatType)
	assert.Equal(t, tp.TypeIdentifier(), dgo.TiBigFloatExact)
}

func TestBigFloat_CompareTo(t *testing.T) {
	b := newBigFloat(t, `123.45`)
	c, ok := b.CompareTo(b)
	require.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = b.CompareTo(big.NewFloat(140.12))
	require.True(t, ok)
	assert.Equal(t, -1, c)

	_, ok = b.CompareTo(`123.45`)
	assert.False(t, ok)

	c, ok = b.CompareTo(nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(120)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(vf.Integer(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(120.12)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(vf.Float(140.14))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(float32(120.12))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(float32(140.14))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(int8(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(int16(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(int8(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(int32(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(int64(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(uint8(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(uint16(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(uint32(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(uint64(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(uint(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = b.CompareTo(big.NewInt(120))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = b.CompareTo(big.NewInt(140))
	assert.True(t, ok)
	assert.Equal(t, -1, c)
}

func TestBigFloat_Equals(t *testing.T) {
	b := vf.BigFloat(big.NewFloat(140.12))
	assert.True(t, b.Equals(b))
	assert.True(t, b.Equals(big.NewFloat(140.12)))
	assert.True(t, b.Equals(vf.BigFloat(big.NewFloat(140.12))))
	assert.True(t, b.Equals(140.12))
	assert.True(t, b.Equals(vf.Float(140.12)))

	f32 := float32(140.12)
	b = vf.BigFloat(big.NewFloat(float64(f32)))
	assert.True(t, b.Equals(float32(140.12)))
}

func TestBigFloat_Float(t *testing.T) {
	b := vf.BigFloat(big.NewFloat(140.12))
	assert.Same(t, b, b.Float())
}

func TestBigFloat_GoFloat(t *testing.T) {
	assert.Equal(t, 123.45, vf.BigFloat(big.NewFloat(123.45)).GoFloat())
}

func TestBigFloat_HashCode(t *testing.T) {
	assert.Equal(t, vf.BigFloat(big.NewFloat(123.45)).HashCode(), vf.BigFloat(big.NewFloat(123.45)).HashCode())
	assert.NotEqual(t, vf.BigFloat(big.NewFloat(123.46)).HashCode(), vf.BigFloat(big.NewFloat(123.45)).HashCode())
}

func TestBigFloat_GoFloat_out_of_bounds(t *testing.T) {
	assert.Panic(t, func() {
		newBigFloat(t, `1e10000`).GoFloat()
	}, `cannot fit`)
}

func TestBigFloat_Integer(t *testing.T) {
	bi := new(big.Int).Exp(big.NewInt(10), big.NewInt(1000), nil)
	bf := vf.BigFloat(new(big.Float).SetInt(bi))
	iv := bf.Integer()
	assert.Equal(t, iv, bi)

	_, ok := bf.ToInt()
	assert.False(t, ok)

	bf = vf.BigFloat(big.NewFloat(123.45))
	i, ok := bf.ToInt()
	assert.True(t, ok)
	assert.Equal(t, i, 123)
}

func TestBigFloat_ReflectTo(t *testing.T) {
	bf := newBigFloat(t, `1234.5678`)
	var mi interface{}
	bf.ReflectTo(reflect.ValueOf(&mi).Elem())
	fc, ok := mi.(*big.Float)
	require.True(t, ok)
	assert.True(t, fc.Cmp(bf.GoBigFloat()) == 0)

	var mf *big.Float
	bf.ReflectTo(reflect.ValueOf(&mf).Elem())
	assert.True(t, mf.Cmp(bf.GoBigFloat()) == 0)

	var m big.Float
	bf.ReflectTo(reflect.ValueOf(&m).Elem())
	assert.True(t, m.Cmp(bf.GoBigFloat()) == 0)
}

func TestBigFloat_String(t *testing.T) {
	// Number that requires more than 53 bits precision
	b := vf.BigFloat(big.NewFloat(123.45))
	assert.Equal(t, `big 123.45`, b.String())
}

func TestBigFloat_ToBigFloat(t *testing.T) {
	// Number that requires more than 53 bits precision
	b := big.NewFloat(123.45)
	assert.True(t, b == vf.BigFloat(b).ToBigFloat())
}

func TestBigFloat_ToBigInt(t *testing.T) {
	// Number that requires more than 53 bits precision
	b := big.NewFloat(123.45)
	assert.Equal(t, big.NewInt(123), vf.BigFloat(b).ToBigInt())
}

func TestBigFloat_ToFloat(t *testing.T) {
	// Number that requires more than 53 bits precision
	f := vf.Float(123.45)
	bf := vf.BigFloat(big.NewFloat(123.45))
	f2, ok := bf.ToFloat()
	assert.True(t, ok)
	assert.Equal(t, f, f2)
}

func TestBigFloat_ToFloat_out_of_bounds(t *testing.T) {
	// Number that requires more than 53 bits precision
	_, ok := newBigFloat(t, `1e10000`).ToFloat()
	assert.False(t, ok)
	_, ok = newBigFloat(t, `1e-10000`).ToFloat()
	assert.False(t, ok)
	_, ok = newBigFloat(t, `-1e10000`).ToFloat()
	assert.False(t, ok)
	_, ok = newBigFloat(t, `-1e-10000`).ToFloat()
	assert.False(t, ok)
}

func TestBigFloat_range_specific(t *testing.T) {
	ft := newBigFloat(t, `1e10000`).(dgo.FloatType)
	assert.True(t, ft.Inclusive())
	assert.Same(t, ft, ft.Min())
	assert.Same(t, ft, ft.Max())
}
