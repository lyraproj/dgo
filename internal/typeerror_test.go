package internal_test

import (
	"testing"

	require "github.com/lyraproj/got/dgo_test"
	"github.com/lyraproj/got/newtype"
	"github.com/lyraproj/got/typ"
	"github.com/lyraproj/got/vf"
)

func TestTypeError(t *testing.T) {
	v := newtype.IllegalAssignment(typ.String, vf.Integer(3))

	require.Equal(t, v, newtype.IllegalAssignment(typ.String, vf.Integer(3)))
	require.NotEqual(t, v, newtype.IllegalAssignment(typ.String, vf.Integer(4)))
	require.NotEqual(t, v, `oops`)

	require.Instance(t, v.Type(), v)
	require.NotEqual(t, 0, v.HashCode())
	require.Equal(t, v.HashCode(), v.HashCode())

	v = newtype.IllegalSize(newtype.String(1, 10), 12)

	require.Equal(t, v, newtype.IllegalSize(newtype.String(1, 10), 12))
	require.NotEqual(t, v, newtype.IllegalSize(newtype.String(1, 10), 11))
	require.NotEqual(t, v, `oops`)

	require.Instance(t, v.Type(), v)
	require.NotEqual(t, 0, v.HashCode())
	require.Equal(t, v.HashCode(), v.HashCode())
}
