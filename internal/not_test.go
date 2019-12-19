package internal_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestNotType(t *testing.T) {
	notNil := tf.Not(typ.Nil).(dgo.UnaryType)
	require.Instance(t, notNil, `b`)
	require.NotInstance(t, notNil, vf.Nil)
	require.Same(t, notNil.Operand(), typ.Nil)
	require.Same(t, tf.Not(tf.Not(typ.Nil)), typ.Nil)

	require.Equal(t, notNil, tf.Not(typ.Nil))
	require.NotEqual(t, notNil, tf.Not(typ.True))
	require.NotEqual(t, notNil, `Not`)

	require.Assignable(t, tf.Not(typ.String), typ.Float)
	require.Assignable(t, tf.Not(typ.String), tf.Not(tf.AnyOf(typ.String, typ.Integer)))
	require.NotAssignable(t, tf.Not(typ.Float), tf.Not(tf.AnyOf(typ.String, typ.Integer)))

	require.Assignable(t, tf.Not(typ.Float), tf.AnyOf(typ.String, typ.Integer, typ.Boolean))

	require.Assignable(t, tf.Not(typ.Float), tf.AllOf(typ.String, typ.Integer, typ.Boolean))
	require.NotAssignable(t, tf.Not(typ.Float), tf.OneOf(typ.String, typ.Integer, typ.Boolean))

	require.NotAssignable(t, tf.Not(typ.Float), tf.AllOf(typ.String, typ.Float))
	require.Assignable(t, tf.Not(typ.Float), tf.OneOf(typ.String, typ.Float))

	require.Assignable(t, tf.Not(typ.String), tf.AnyOf(tf.Not(typ.String), tf.Not(typ.Integer)))

	require.NotAssignable(t, tf.Not(typ.Float), tf.AnyOf(tf.Not(typ.String), tf.Not(typ.Integer), tf.Not(typ.Boolean)))
	require.NotAssignable(t, tf.Not(typ.Float), tf.OneOf(tf.Not(typ.String), tf.Not(typ.Integer), tf.Not(typ.Boolean)))
	require.NotAssignable(t, tf.Not(typ.Float), tf.AllOf(tf.Not(typ.String), tf.Not(typ.Integer), tf.Not(typ.Boolean)))

	require.Instance(t, notNil.Type(), notNil)

	require.Equal(t, notNil.HashCode(), notNil.HashCode())
	require.NotEqual(t, 0, notNil.HashCode())

	require.Equal(t, `!nil`, notNil.String())
	require.Equal(t, dgo.OpNot, notNil.Operator())

	require.Equal(t, notNil.ReflectType(), typ.Any.ReflectType())
}
