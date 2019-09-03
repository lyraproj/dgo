package internal_test

import (
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestNil(t *testing.T) {
	c, ok := vf.Nil.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.Nil.CompareTo(vf.Integer(1))
	require.True(t, ok)
	require.Equal(t, -1, c)

	nt := typ.Nil
	require.Assignable(t, nt, nt)
	require.Instance(t, nt, vf.Nil)
	require.Assignable(t, typ.Any, nt)
	require.NotAssignable(t, nt, typ.Any)

	require.Equal(t, nt, nt)
	require.NotEqual(t, nt, typ.Any)

	require.Instance(t, nt.Type(), nt)
	require.NotInstance(t, nt.Type(), typ.Any)

	require.Equal(t, `nil`, nt.String())

	require.NotEqual(t, 0, nt.HashCode())
	require.Equal(t, nt.HashCode(), nt.HashCode())
}
