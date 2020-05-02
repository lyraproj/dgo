package internal_test

import (
	"math"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestTupleType(t *testing.T) {
	tt := typ.Tuple
	assert.Same(t, tt, tf.VariadicTuple(typ.Any))
	assert.Same(t, typ.Any, tt.ElementType())

	tt = typ.EmptyTuple
	assert.Same(t, tt, tf.Tuple())
	assert.Assignable(t, tt, tf.Array(0, 0))
	assert.Same(t, typ.Any, tt.ElementType())

	tt = tf.Tuple(typ.String, typ.Any, typ.Float)
	assert.Assignable(t, tt, tf.Tuple(typ.String, typ.Integer, typ.Float))
	assert.Assignable(t, tt, vf.Values(`one`, 2, 3.0).Type())
	assert.Equal(t, 3, tt.Min())
	assert.Equal(t, 3, tt.Max())
	assert.False(t, tt.Unbounded())

	et := tf.String(1, 100)
	tt = tf.Tuple(et)
	assert.Same(t, tt.ElementType(), et)

	tt = tf.Tuple(typ.String, typ.Integer, typ.Float)
	assert.Assignable(t, tt, tt)
	assert.Assignable(t, tt, tf.Tuple(typ.String, typ.Integer, typ.Float))
	assert.NotAssignable(t, tt, tf.Tuple(typ.String, typ.Integer, typ.Boolean))
	assert.Assignable(t, typ.Array, tt)
	assert.Assignable(t, tf.Array(0, 3), tt)
	assert.Assignable(t, tt, vf.Values(`one`, 2, 3.0).Type())
	assert.NotAssignable(t, tt, vf.Values(`one`, 2, 3.0, `four`).Type())
	assert.NotAssignable(t, tt, vf.Values(`one`, 2, 3).Type())
	assert.NotAssignable(t, tt, typ.Array)
	assert.NotAssignable(t, tt, typ.Tuple)
	assert.NotEqual(t, tt, typ.String)
	assert.Equal(t, 3, tt.Min())
	assert.Equal(t, 3, tt.Max())
	assert.False(t, tt.Unbounded())
	assert.Assignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer, typ.Float)), tt)
	assert.NotAssignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer)), tt)
	assert.Assignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer, typ.Float, typ.Boolean)), tt)
	assert.Assignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer, typ.Float), 0, 3), tt)
	assert.NotAssignable(t, tf.Array(tf.AnyOf(typ.String, typ.Integer, typ.Float), 0, 2), tt)

	okv := vf.Values(`hello`, 1, 2.0)
	assert.Instance(t, tt, okv)
	assert.NotInstance(t, tt, okv.Get(0))
	assert.Assignable(t, tt, okv.Type())
	assert.Assignable(t, typ.Array, tt)

	okv = vf.Values(`hello`, 1, 2)
	assert.NotInstance(t, tt, okv)

	okv = vf.Values(`hello`, 1, 2.0, true)
	assert.NotInstance(t, tt, okv)

	okm := vf.MutableValues(`world`, 2, 3.0)
	assert.Equal(t, `world`, okm.Set(0, `earth`))

	tt = tf.Tuple(typ.String, typ.String)
	assert.Assignable(t, tt, tf.Array(typ.String, 2, 2))
	assert.Equal(t, tt.ReflectType(), tf.Array(typ.String).ReflectType())
	assert.Assignable(t, tt, tf.Array(typ.String, 2, 2))
	assert.NotAssignable(t, tt, tf.AnyOf(typ.Nil, tf.Array(typ.String, 2, 2)))
	tt = tf.Tuple(typ.String, typ.Integer)
	assert.NotAssignable(t, tt, tf.Array(typ.String, 2, 2))
	assert.Equal(t, tt.ReflectType(), typ.Array.ReflectType())
	assert.Equal(t, typ.Any, typ.Tuple.ElementType())
	assert.Equal(t, tf.AllOf(typ.String, typ.Integer), tt.ElementType())
	assert.Equal(t, vf.Values(typ.String, typ.Integer), tt.ElementTypes())
	assert.Equal(t, typ.Array, typ.Generic(tt))
	assert.Instance(t, tt.Type(), tt)
	assert.Equal(t, `{string,int}`, tt.String())
	assert.NotEqual(t, 0, tt.HashCode())
	assert.Equal(t, tt.HashCode(), tt.HashCode())

	tt = tf.Tuple(typ.String, typ.String)
	assert.Equal(t, tf.Array(typ.String), typ.Generic(tt))

	te := tf.Tuple(vf.String(`a`).Type(), vf.String(`b`).Type())
	assert.Assignable(t, tt, te)
	assert.NotAssignable(t, te, tt)
	assert.Equal(t, tf.Array(typ.String), typ.Generic(te))
}

func TestTupleType_selfReference(t *testing.T) {
	internal.ResetDefaultAliases()
	tp := tf.ParseType(`x={string,x}`).(dgo.ArrayType)
	d := vf.MutableValues()
	d.Add(`hello`)
	d.Add(d)
	assert.Instance(t, tp, d)

	internal.ResetDefaultAliases()
	t2 := tf.ParseType(`x={string,{string,x}}`)
	assert.Assignable(t, tp, t2)
	assert.Equal(t, `x`, t2.String())

	internal.ResetDefaultAliases()
	assert.Equal(t, `{string,<recursive self reference to tuple type>}`, tp.String())
}

func TestVariadicTupleType(t *testing.T) {
	tt := tf.ParseType(`{string,...string}`).(dgo.TupleType)
	assert.Instance(t, tt, vf.Values(`one`))
	assert.Instance(t, tt, vf.Values(`one`, `two`))
	assert.NotInstance(t, tt, vf.Values(`one`, 2))
	assert.NotInstance(t, tt, vf.Values(1))

	tt = tf.VariadicTuple(typ.String, typ.String)
	assert.NotInstance(t, tt, vf.Values())
	assert.Instance(t, tt, vf.Values(`one`))
	assert.Instance(t, tt, vf.Values(`one`, `two`))
	assert.Instance(t, tt, vf.Values(`one`, `two`, `three`))
	assert.True(t, tt.Unbounded())
	assert.Equal(t, 1, tt.Min())
	assert.Equal(t, math.MaxInt64, tt.Max())
	assert.Panic(t, func() { tf.VariadicTuple() }, `must have at least one element`)
}
