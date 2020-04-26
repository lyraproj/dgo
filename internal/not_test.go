package internal_test

import (
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestNotType(t *testing.T) {
	notNil := tf.Not(typ.Nil).(dgo.UnaryType)
	assert.Instance(t, notNil, `b`)
	assert.NotInstance(t, notNil, vf.Nil)
	assert.Same(t, notNil.Operand(), typ.Nil)
	assert.Same(t, tf.Not(tf.Not(typ.Nil)), typ.Nil)
	assert.Equal(t, notNil, tf.Not(typ.Nil))
	assert.NotEqual(t, notNil, tf.Not(typ.True))
	assert.NotEqual(t, notNil, `Not`)
	assert.Assignable(t, tf.Not(typ.String), typ.Float)
	assert.Assignable(t, tf.Not(typ.String), tf.Not(tf.AnyOf(typ.String, typ.Integer)))
	assert.NotAssignable(t, tf.Not(typ.Float), tf.Not(tf.AnyOf(typ.String, typ.Integer)))
	assert.Assignable(t, tf.Not(typ.Float), tf.AnyOf(typ.String, typ.Integer, typ.Boolean))
	assert.Assignable(t, tf.Not(typ.Float), tf.AllOf(typ.String, typ.Integer, typ.Boolean))
	assert.NotAssignable(t, tf.Not(typ.Float), tf.OneOf(typ.String, typ.Integer, typ.Boolean))
	assert.NotAssignable(t, tf.Not(typ.Float), tf.AllOf(typ.String, typ.Float))
	assert.Assignable(t, tf.Not(typ.Float), tf.OneOf(typ.String, typ.Float))
	assert.Assignable(t, tf.Not(typ.String), tf.AnyOf(tf.Not(typ.String), tf.Not(typ.Integer)))
	assert.NotAssignable(t, tf.Not(typ.Float), tf.AnyOf(tf.Not(typ.String), tf.Not(typ.Integer), tf.Not(typ.Boolean)))
	assert.NotAssignable(t, tf.Not(typ.Float), tf.OneOf(tf.Not(typ.String), tf.Not(typ.Integer), tf.Not(typ.Boolean)))
	assert.NotAssignable(t, tf.Not(typ.Float), tf.AllOf(tf.Not(typ.String), tf.Not(typ.Integer), tf.Not(typ.Boolean)))
	assert.Instance(t, notNil.Type(), notNil)
	assert.Equal(t, notNil.HashCode(), notNil.HashCode())
	assert.NotEqual(t, 0, notNil.HashCode())
	assert.Equal(t, `!nil`, notNil.String())
	assert.Equal(t, dgo.OpNot, notNil.Operator())
	assert.Equal(t, notNil.ReflectType(), typ.Any.ReflectType())
}
