package internal_test

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestFloat(t *testing.T) {
	require.Instance(t, typ.Float, 3.1415927)
	require.NotInstance(t, typ.Float, true)
	require.Assignable(t, typ.Float, typ.Float)
	require.Assignable(t, typ.Float, tf.FloatRange(3.1, 5.8, true))
	require.Assignable(t, typ.Float, vf.Float(4.2).Type())
	require.Equal(t, typ.Float, typ.Float)
	require.Equal(t, typ.Float, tf.FloatRange(-math.MaxFloat64, math.MaxFloat64, true))
	require.Instance(t, typ.Float.Type(), typ.Float)
	require.True(t, typ.Float.IsInstance(1234))
	require.Equal(t, typ.Float.Min(), -math.MaxFloat64)
	require.Equal(t, typ.Float.Max(), math.MaxFloat64)
	require.True(t, typ.Float.Inclusive())

	require.Equal(t, typ.Float.HashCode(), typ.Float.HashCode())
	require.NotEqual(t, 0, typ.Float.HashCode())

	require.Equal(t, `float`, typ.Float.String())

	require.Equal(t, reflect.ValueOf(3.14).Type(), typ.Float.ReflectType())
}

func TestFloatExact(t *testing.T) {
	tp := vf.Float(3.1415927).Type().(dgo.FloatRangeType)
	require.Instance(t, tp, 3.1415927)
	require.NotInstance(t, tp, 2.05)
	require.NotInstance(t, tp, true)
	require.Assignable(t, tf.FloatRange(3.0, 5.0, true), tp)
	require.Assignable(t, tp, tf.FloatRange(3.1415927, 3.1415927, true))
	require.NotAssignable(t, tp, typ.Float)
	require.Equal(t, tp, tf.FloatRange(3.1415927, 3.1415927, true))
	require.NotEqual(t, tp, tf.FloatRange(2.1, 5.5, true))
	require.Equal(t, tp.Min(), 3.1415927)
	require.Equal(t, tp.Max(), 3.1415927)
	require.True(t, tp.Inclusive())
	require.True(t, tp.IsInstance(3.1415927))
	require.Equal(t, `3.1415927`, tp.String())
	require.Equal(t, `3.0`, vf.Float(3.0).Type().String())

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `3.1415927`, tp.String())

	require.Same(t, typ.Float, typ.Generic(tp))

	require.Instance(t, tp.Type(), tp)

	require.Same(t, tp.ReflectType(), typ.Float.ReflectType())
}

func TestFloatRange(t *testing.T) {
	tp := tf.FloatRange(3.1, 5.8, true)
	require.Instance(t, tp, 3.1415927)
	require.NotInstance(t, tp, 3.01)
	require.NotInstance(t, tp, true)
	require.Assignable(t, tp, tf.FloatRange(3.1, 5.8, true))
	require.Assignable(t, tp, tf.FloatRange(4.2, 4.2, true))
	require.Assignable(t, tp, vf.Float(4.2).Type())
	require.NotAssignable(t, tp, tf.FloatRange(2.5, 5.5, true))
	require.NotAssignable(t, tp, tf.FloatRange(3.1, 6.2, true))
	require.NotAssignable(t, tp, vf.Float(6.0).Type())
	require.NotAssignable(t, tp, typ.Float)
	require.Equal(t, tp, tf.FloatRange(5.8, 3.1, true))
	require.NotEqual(t, tp, tf.FloatRange(2.5, 5.5, true))
	require.NotEqual(t, tp, typ.Float)
	require.Equal(t, tp.Min(), 3.1)
	require.Equal(t, tp.Max(), 5.8)
	require.Equal(t, vf.Float(4.2).Type(), tf.FloatRange(4.2, 4.2, true))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `3.1..5.8`, tp.String())

	require.Instance(t, tp.Type(), tp)

	tp = tf.FloatRange(3.1, 5.8, false)
	require.Instance(t, tp, 3.1415927)
	require.NotInstance(t, tp, 5.8)
	require.Assignable(t, tp, tf.FloatRange(3.1, 5.8, false))
	require.NotAssignable(t, tp, tf.FloatRange(3.1, 5.8, true))
	require.Assignable(t, tf.FloatRange(3.1, 5.8, true), tp)
	require.Equal(t, `3.1...5.8`, tp.String())

	require.Panic(t, func() { tf.FloatRange(3.1, 3.1, false) }, `cannot have equal min and max`)

	require.Same(t, tp.ReflectType(), typ.Float.ReflectType())
}

func TestFloatType_New(t *testing.T) {
	require.Equal(t, 17.3, vf.New(typ.Float, vf.Float(17.3)))
	require.Equal(t, 17.3, vf.New(typ.Float, vf.String(`17.3`)))
	require.Equal(t, 17.0, vf.New(typ.Float, vf.Integer(17)))
	require.Equal(t, 0.0, vf.New(typ.Float, vf.False))
	require.Equal(t, 1.0, vf.New(typ.Float, vf.True))

	now := time.Now()
	require.Equal(t, vf.New(typ.Float, vf.Time(now)), vf.New(typ.Float, vf.Arguments(vf.Time(now))))

	require.Panic(t, func() { vf.New(typ.Float, vf.String(`true`)) }, `cannot be converted`)
	require.Panic(t, func() { vf.New(vf.Float(4).Type(), vf.Float(5)) }, `cannot be assigned`)
	require.Panic(t, func() { vf.New(tf.FloatRange(1, 4, true), vf.Float(5)) }, `cannot be assigned`)
}

func TestNumber(t *testing.T) {
	require.Equal(t, vf.Float(3.14).ToInt(), vf.Integer(3).ToInt())
	require.Equal(t, vf.Integer(3).ToFloat(), vf.Float(3.0).ToFloat())
	require.NotEqual(t, vf.Float(3.14).ToFloat(), vf.Integer(3).ToFloat())
}

func TestFloat_CompareTo(t *testing.T) {
	c, ok := vf.Float(3.14).CompareTo(vf.Float(3.14))
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.Float(3.14).CompareTo(vf.Float(-3.14))
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Float(-3.14).CompareTo(vf.Float(3.14))
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.Float(3).CompareTo(vf.Integer(3))
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.Float(2.9).CompareTo(vf.Integer(3))
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.Float(3.1).CompareTo(vf.Integer(3))
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Float(3.14).CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Float(3.14).CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	_, ok = vf.Float(3.14).CompareTo(vf.True)
	require.False(t, ok)
}

func TestFloat_ReflectTo(t *testing.T) {
	var f float64
	fv := vf.Float(0.5403)
	fv.ReflectTo(reflect.ValueOf(&f).Elem())
	require.Equal(t, f, fv)

	var fp *float64
	fv.ReflectTo(reflect.ValueOf(&fp).Elem())
	require.Equal(t, f, *fp)

	var mi interface{}
	mip := &mi
	fv.ReflectTo(reflect.ValueOf(mip).Elem())
	fc, ok := mi.(float64)
	require.True(t, ok)
	require.Equal(t, fv, fc)

	fx := float32(0.5403)
	fv = vf.Float(float64(fx))
	var f32 float32
	fv.ReflectTo(reflect.ValueOf(&f32).Elem())
	require.Equal(t, f32, fv)

	var fp32 *float32
	fv.ReflectTo(reflect.ValueOf(&fp32).Elem())
	require.Equal(t, fx, *fp32)
}

func TestFloat_String(t *testing.T) {
	require.Equal(t, `1234.0`, vf.Float(1234).String())
	require.Equal(t, `-4321.0`, vf.Float(-4321).String())
}
