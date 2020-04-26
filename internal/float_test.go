package internal_test

import (
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/tada/dgo/internal"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/test/require"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestFloat(t *testing.T) {
	assert.Instance(t, typ.Float, 3.1415927)
	assert.NotInstance(t, typ.Float, true)
	assert.Assignable(t, typ.Float, typ.Float)
	assert.Assignable(t, typ.Float, tf.Float64(3.1, 5.8, true))
	assert.Assignable(t, typ.Float, vf.Float(4.2).Type())
	assert.Equal(t, typ.Float, typ.Float)
	assert.Equal(t, typ.Float, tf.Float64(-math.MaxFloat64, math.MaxFloat64, true))
	assert.Equal(t, typ.Float, tf.Float(nil, nil, true))
	assert.Instance(t, typ.Float.Type(), typ.Float)
	assert.True(t, typ.Float.Instance(1234.0))
	assert.Equal(t, typ.Float.Min(), nil)
	assert.Equal(t, typ.Float.Max(), nil)
	assert.True(t, typ.Float.Inclusive())
	assert.Equal(t, typ.Float.HashCode(), typ.Float.HashCode())
	assert.NotEqual(t, 0, typ.Float.HashCode())
	assert.Equal(t, `float`, typ.Float.String())
	assert.Equal(t, reflect.ValueOf(3.14).Type(), typ.Float.ReflectType())
}

func TestFloatExact(t *testing.T) {
	tp := vf.Float(3.1415927).Type().(dgo.FloatType)
	assert.Instance(t, tp, 3.1415927)
	assert.NotInstance(t, tp, 2.05)
	assert.NotInstance(t, tp, true)
	assert.Assignable(t, tf.Float64(3.0, 5.0, true), tp)
	assert.Assignable(t, tp, tf.Float64(3.1415927, 3.1415927, true))
	assert.NotAssignable(t, tp, typ.Float)
	assert.Equal(t, tp, tf.Float64(3.1415927, 3.1415927, true))
	assert.NotEqual(t, tp, tf.Float64(2.1, 5.5, true))
	assert.Equal(t, tp.Min(), 3.1415927)
	assert.Equal(t, tp.Max(), 3.1415927)
	assert.True(t, tp.Inclusive())
	assert.True(t, tp.Instance(3.1415927))
	assert.Equal(t, `3.1415927`, tp.String())
	assert.Equal(t, `3.0`, vf.Float(3.0).Type().String())
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())

	assert.Equal(t, `3.1415927`, tp.String())

	assert.Same(t, typ.Float, typ.Generic(tp))

	assert.Instance(t, tp.Type(), tp)

	assert.Same(t, tp.ReflectType(), typ.Float.ReflectType())
}

func TestFloatRange(t *testing.T) {
	tp := tf.Float64(3.1, 5.8, true)
	assert.Instance(t, tp, 3.1415927)
	assert.NotInstance(t, tp, 3.01)
	assert.NotInstance(t, tp, true)
	assert.Assignable(t, tp, tf.Float64(3.1, 5.8, true))
	assert.Assignable(t, tp, tf.Float64(4.2, 4.2, true))
	assert.Assignable(t, tp, vf.Float(4.2).Type())
	assert.Assignable(t, tp, tf.AnyOf(vf.Float(4.2).Type(), vf.Float(5.2).Type()))
	assert.NotAssignable(t, tp, tf.Float64(2.5, 5.5, true))
	assert.NotAssignable(t, tp, tf.Float64(3.1, 6.2, true))
	assert.NotAssignable(t, tp, vf.Float(6.0).Type())
	assert.NotAssignable(t, tp, typ.Float)
	assert.NotAssignable(t, tp, tf.Float(vf.Float(3.1), nil, true))
	assert.NotAssignable(t, tp, tf.Float(nil, vf.Float(5.8), true))
	assert.Equal(t, tp, tf.Float64(5.8, 3.1, true))
	assert.Equal(t, tp, tf.Float(vf.Float(5.8), vf.Float(3.1), true))
	assert.NotEqual(t, tp, tf.Float64(2.5, 5.5, true))
	assert.NotEqual(t, tp, typ.Float)
	assert.Equal(t, tp.Min(), 3.1)
	assert.Equal(t, tp.Max(), 5.8)
	assert.Equal(t, vf.Float(4.2).Type(), tf.Float64(4.2, 4.2, true))
	assert.Equal(t, vf.Float(4.2).Type(), tf.Float(vf.Float(4.2), vf.Float(4.2), true))

	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())

	assert.Equal(t, `3.1..5.8`, tp.String())

	assert.Instance(t, tp.Type(), tp)

	tp = tf.Float64(3.1, 5.8, false)
	assert.Instance(t, tp, 3.1415927)
	assert.Instance(t, tp, float32(3.1415927))
	assert.Instance(t, tp, big.NewFloat(3.1415927))
	assert.NotInstance(t, tp, 5.8)
	assert.Assignable(t, tp, tf.Float64(3.1, 5.8, false))
	assert.NotAssignable(t, tp, tf.Float64(3.1, 5.8, true))
	assert.Assignable(t, tf.Float64(3.1, 5.8, true), tp)
	assert.Equal(t, `3.1...5.8`, tp.String())

	assert.Panic(t, func() { tf.Float64(3.1, 3.1, false) }, `cannot have equal min and max`)
	assert.Panic(t, func() { tf.Float(vf.Float(3.1), vf.Float(3.1), false) }, `cannot have equal min and max`)

	assert.Same(t, tf.Float64(-math.MaxFloat64, math.MaxFloat64, true), tf.Float(nil, nil, true))

	assert.Same(t, tp.ReflectType(), typ.Float.ReflectType())
}

func TestFloatType_New(t *testing.T) {
	assert.Equal(t, 17.3, vf.New(typ.Float, vf.Float(17.3)))
	assert.Equal(t, 17.3, vf.New(typ.Float, vf.String(`17.3`)))
	assert.Equal(t, 17.0, vf.New(typ.Float, vf.Integer(17)))
	assert.Equal(t, 0.0, vf.New(typ.Float, vf.False))
	assert.Equal(t, 1.0, vf.New(typ.Float, vf.True))

	now := time.Now()
	assert.Equal(t, vf.New(typ.Float, vf.Time(now)), vf.New(typ.Float, vf.Arguments(vf.Time(now))))

	assert.Panic(t, func() { vf.New(typ.Float, vf.String(`true`)) }, `cannot be converted`)
	assert.Panic(t, func() { vf.New(vf.Float(4).Type(), vf.Float(5)) }, `cannot be assigned`)
	assert.Panic(t, func() { vf.New(tf.Float64(1, 4, true), vf.Float(5)) }, `cannot be assigned`)
}

func TestNumber(t *testing.T) {
	ai, _ := vf.Float(3.14).ToInt()
	bi, _ := vf.Integer(3).ToInt()
	assert.Equal(t, ai, bi)

	af, _ := vf.Integer(3).ToFloat()
	bf, _ := vf.Float(3.0).ToFloat()
	assert.Equal(t, af, bf)

	af, _ = vf.Float(3.14).ToFloat()
	bf, _ = vf.Integer(3).ToFloat()
	assert.NotEqual(t, af, bf)
}

func TestFloat_CompareTo(t *testing.T) {
	c, ok := vf.Float(3.14).CompareTo(vf.Float(3.14))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Float(3.14).CompareTo(3.14)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Float(3.14).CompareTo(vf.Float(3.14).ToBigFloat())
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Float(3.14).CompareTo(vf.Float(3.14).ToBigInt())
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Float(float64(float32(3.14))).CompareTo(float32(3.14))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Float(3.14).CompareTo(vf.Float(-3.14))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Float(-3.14).CompareTo(vf.Float(3.14))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Float(math.MaxFloat64).CompareTo(newBigFloat(t, `1e1000`))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Float(-math.MaxFloat64).CompareTo(newBigFloat(t, `-1e1000`))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Float(0).CompareTo(newBigFloat(t, `1e-1000`))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Float(0).CompareTo(newBigFloat(t, `-1e-1000`))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Float(-1e-50).CompareTo(newBigFloat(t, `1e-1000`))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Float(1e-50).CompareTo(newBigFloat(t, `-1e-1000`))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Float(3).CompareTo(vf.Integer(3))
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Float(2.9).CompareTo(vf.Integer(3))
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.Float(3.1).CompareTo(vf.Integer(3))
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Float(3.14).CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.Float(3.14).CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	_, ok = vf.Float(3.14).CompareTo(vf.True)
	assert.False(t, ok)
}

func TestFloat_ReflectTo(t *testing.T) {
	var f float64
	fv := vf.Float(0.5403)
	fv.ReflectTo(reflect.ValueOf(&f).Elem())
	assert.Equal(t, f, fv)

	var fp *float64
	fv.ReflectTo(reflect.ValueOf(&fp).Elem())
	assert.Equal(t, f, *fp)

	var mi interface{}
	mip := &mi
	fv.ReflectTo(reflect.ValueOf(mip).Elem())
	fc, ok := mi.(float64)
	require.True(t, ok)
	assert.Equal(t, fv, fc)

	fx := float32(0.5403)
	fv = vf.Float(float64(fx))
	var f32 float32
	fv.ReflectTo(reflect.ValueOf(&f32).Elem())
	assert.Equal(t, f32, fv)

	var fp32 *float32
	fv.ReflectTo(reflect.ValueOf(&fp32).Elem())
	assert.Equal(t, fx, *fp32)
}

func TestFloat_String(t *testing.T) {
	assert.Equal(t, `1234.0`, vf.Float(1234).String())
	assert.Equal(t, `-4321.0`, vf.Float(-4321).String())
}

func TestToFloat(t *testing.T) {
	bf := big.NewFloat(3.14159265)
	f, ok := internal.ToFloat(bf)
	require.True(t, ok)
	assert.Equal(t, 3.14159265, f)
	f, ok = internal.ToFloat(vf.BigFloat(bf))
	require.True(t, ok)
	assert.Equal(t, 3.14159265, f)
}
