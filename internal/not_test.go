package internal_test

import (
	"testing"

	"github.com/lyraproj/got/dgo"
	require "github.com/lyraproj/got/dgo_test"
	"github.com/lyraproj/got/newtype"
	"github.com/lyraproj/got/typ"
	"github.com/lyraproj/got/vf"
)

func TestNotType(t *testing.T) {
	notNil := newtype.Not(typ.Nil).(dgo.UnaryType)
	require.Instance(t, notNil, `b`)
	require.NotInstance(t, notNil, vf.Nil)
	require.Same(t, notNil.Operand(), typ.Nil)
	require.Same(t, newtype.Not(newtype.Not(typ.Nil)), typ.Nil)

	require.Equal(t, notNil, newtype.Not(typ.Nil))
	require.NotEqual(t, notNil, newtype.Not(typ.True))
	require.NotEqual(t, notNil, `Not`)

	require.Assignable(t, newtype.Not(typ.String), typ.Float)
	require.Assignable(t, newtype.Not(typ.String), newtype.Not(newtype.AnyOf(typ.String, typ.Integer)))
	require.NotAssignable(t, newtype.Not(typ.Float), newtype.Not(newtype.AnyOf(typ.String, typ.Integer)))

	require.Assignable(t, newtype.Not(typ.Float), newtype.AnyOf(typ.String, typ.Integer, typ.Boolean))

	require.Assignable(t, newtype.Not(typ.Float), newtype.AllOf(typ.String, typ.Integer, typ.Boolean))
	require.NotAssignable(t, newtype.Not(typ.Float), newtype.OneOf(typ.String, typ.Integer, typ.Boolean))

	require.NotAssignable(t, newtype.Not(typ.Float), newtype.AllOf(typ.String, typ.Float))
	require.Assignable(t, newtype.Not(typ.Float), newtype.OneOf(typ.String, typ.Float))

	require.Assignable(t, newtype.Not(typ.String), newtype.AnyOf(newtype.Not(typ.String), newtype.Not(typ.Integer)))

	require.NotAssignable(t, newtype.Not(typ.Float), newtype.AnyOf(newtype.Not(typ.String), newtype.Not(typ.Integer), newtype.Not(typ.Boolean)))
	require.NotAssignable(t, newtype.Not(typ.Float), newtype.OneOf(newtype.Not(typ.String), newtype.Not(typ.Integer), newtype.Not(typ.Boolean)))
	require.NotAssignable(t, newtype.Not(typ.Float), newtype.AllOf(newtype.Not(typ.String), newtype.Not(typ.Integer), newtype.Not(typ.Boolean)))

	require.Instance(t, notNil.Type(), notNil)

	require.Equal(t, notNil.HashCode(), notNil.HashCode())
	require.NotEqual(t, 0, notNil.HashCode())

	require.Equal(t, `!nil`, notNil.String())
	require.Equal(t, dgo.OpNot, notNil.Operator())
}
