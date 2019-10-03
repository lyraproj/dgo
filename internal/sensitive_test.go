package internal_test

import (
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/newtype"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/vf"

	"github.com/lyraproj/dgo/typ"

	require "github.com/lyraproj/dgo/dgo_test"
)

func TestSensitiveType(t *testing.T) {
	tp := typ.Sensitive
	s := vf.Sensitive(vf.Integer(0))
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.Any)
	require.Assignable(t, tp, s.Type())
	require.NotAssignable(t, s.Type(), tp)
	require.Assignable(t, newtype.Sensitive(typ.Integer), s.Type())
	require.Instance(t, tp.Type(), s.Type())

	require.Equal(t, tp, tp)
	require.NotEqual(t, tp, typ.Any)

	require.NotEqual(t, 0, tp.HashCode())
	require.Equal(t, tp.HashCode(), tp.HashCode())

	require.Instance(t, tp, s)
	require.NotInstance(t, tp, vf.Integer(0))
	require.True(t, reflect.TypeOf(s).AssignableTo(tp.ReflectType()))
	require.Equal(t, dgo.TiSensitive, tp.TypeIdentifier())
	require.Equal(t, dgo.OpSensitive, tp.Operator())
	require.Equal(t, `sensitive[int]`, s.Type().String())
}

func TestSensitive(t *testing.T) {
	s := vf.Sensitive(vf.String(`a`))
	require.Equal(t, s, s)
	require.Equal(t, s, vf.Sensitive(vf.String(`a`)))
	require.NotEqual(t, s, vf.Sensitive(vf.String(`b`)))
	require.NotEqual(t, s, vf.String(`a`))

	require.True(t, s.Frozen())
	a := vf.MutableValues(nil, `a`)
	s = vf.Sensitive(a)
	require.False(t, s.Frozen())
	s.Freeze()
	require.True(t, s.Frozen())
	require.True(t, a.Frozen())
	require.Same(t, s.Unwrap(), a)

	a = vf.MutableValues(nil, `a`)
	s = vf.Sensitive(a)
	c := s.FrozenCopy().(dgo.Sensitive)
	require.False(t, s.Frozen())
	require.True(t, c.Frozen())
	require.Equal(t, s.Unwrap(), c.Unwrap())
	require.NotSame(t, s.Unwrap(), c.Unwrap())

	s = vf.Sensitive(vf.String(`a`))
	c = s.FrozenCopy().(dgo.Sensitive)
	require.Same(t, s, c)

	require.Equal(t, `sensitive [value redacted]`, s.String())

	require.NotEqual(t, typ.Sensitive.HashCode(), s.HashCode())
	require.Equal(t, s.HashCode(), s.HashCode())
}
