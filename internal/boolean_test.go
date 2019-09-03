package internal_test

import (
	"testing"

	"github.com/lyraproj/got/dgo"
	require "github.com/lyraproj/got/dgo_test"
	"github.com/lyraproj/got/typ"
	"github.com/lyraproj/got/vf"
)

func TestBooleanDefault(t *testing.T) {
	tp := typ.Boolean
	meta := tp.Type()
	require.Instance(t, meta, tp)
	require.NotInstance(t, tp, tp)
	require.NotAssignable(t, meta, tp)
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, meta)
}

func TestBooleanType(t *testing.T) {
	v := vf.True
	tp := v.Type().(dgo.BooleanType)
	require.Instance(t, tp, v)
	require.Instance(t, typ.Boolean, v)
	require.Instance(t, typ.True, v)
	require.Instance(t, tp, true)
	require.Instance(t, typ.Boolean, true)
	require.Instance(t, typ.True, true)
	require.NotInstance(t, typ.False, v)
	require.NotInstance(t, typ.False, true)
	require.Assignable(t, typ.Boolean, tp)
	require.NotAssignable(t, tp, typ.Boolean)
	require.NotEqual(t, v, tp)
	require.True(t, tp.IsInstance(true))
	require.False(t, tp.IsInstance(false))
	require.Equal(t, `true`, tp.String())

	v = vf.False
	tp = v.Type().(dgo.BooleanType)
	require.Instance(t, tp, v)
	require.Instance(t, typ.Boolean, v)
	require.Instance(t, typ.False, v)
	require.NotInstance(t, typ.True, v)
	require.Assignable(t, typ.Boolean, tp)
	require.NotAssignable(t, tp, typ.Boolean)
	require.NotEqual(t, v, tp)
	require.True(t, tp.IsInstance(false))
	require.False(t, tp.IsInstance(true))
	require.Equal(t, `false`, tp.String())

	require.Equal(t, typ.Boolean.HashCode(), typ.Boolean.HashCode())
	require.NotEqual(t, 0, typ.Boolean.HashCode())
	require.NotEqual(t, 0, typ.True.HashCode())
	require.NotEqual(t, 0, typ.False.HashCode())
	require.NotEqual(t, typ.Boolean.HashCode(), typ.True.HashCode())
	require.NotEqual(t, typ.Boolean.HashCode(), typ.False.HashCode())
	require.NotEqual(t, typ.True.HashCode(), typ.False.HashCode())
	require.Equal(t, `bool`, typ.Boolean.String())
}

func TestBoolean(t *testing.T) {
	require.Equal(t, vf.True, vf.Boolean(true))
	require.Equal(t, vf.False, vf.Boolean(false))
}

func TestBoolean_Equals(t *testing.T) {
	require.True(t, vf.True.Equals(vf.True))
	require.True(t, vf.True.Equals(true))
	require.False(t, vf.True.Equals(vf.False))
	require.False(t, vf.True.Equals(false))
	require.True(t, vf.False.Equals(vf.False))
	require.True(t, vf.False.Equals(false))
	require.False(t, vf.False.Equals(vf.True))
	require.False(t, vf.False.Equals(true))
	require.True(t, vf.True.GoBool())
	require.False(t, vf.False.GoBool())
}

func TestBoolean_HashCode(t *testing.T) {
	require.NotEqual(t, 0, vf.True.HashCode())
	require.NotEqual(t, 0, vf.False.HashCode())
	require.Equal(t, vf.True.HashCode(), vf.True.HashCode())
	require.NotEqual(t, vf.True.HashCode(), vf.False.HashCode())
}

func TestBoolean_CompareTo(t *testing.T) {
	c, ok := vf.True.CompareTo(vf.True)
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.True.CompareTo(vf.False)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.False.CompareTo(vf.True)
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.True.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.False.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	_, ok = vf.True.CompareTo(vf.Integer(1))
	require.False(t, ok)
}

func TestBoolean_String(t *testing.T) {
	require.Equal(t, `true`, vf.True.String())
	require.Equal(t, `false`, vf.False.String())
}
