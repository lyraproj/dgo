package internal_test

import (
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/tf"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/vf"

	"github.com/lyraproj/dgo/typ"

	require "github.com/lyraproj/dgo/dgo_test"
)

func TestSensitiveType(t *testing.T) {
	var tp dgo.Type = typ.Sensitive
	s := vf.Sensitive(vf.Integer(0))
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.Any)
	require.Assignable(t, tp, s.Type())
	require.NotAssignable(t, s.Type(), tp)
	require.Assignable(t, tf.Sensitive(typ.Integer), s.Type())
	require.Instance(t, tp.Type(), s.Type())

	tp = s.Type()
	require.Equal(t, tp, tp)
	require.Equal(t, tp, vf.Sensitive(vf.Integer(0)).Type())
	require.Equal(t, tp, vf.Sensitive(vf.Integer(1)).Type()) // type uses generic of wrapped
	require.NotEqual(t, tp, typ.Any)
	require.NotEqual(t, tp, tf.Array(typ.String))

	require.NotEqual(t, 0, tp.HashCode())
	require.Equal(t, tp.HashCode(), tp.HashCode())

	require.Instance(t, tp, s)
	require.NotInstance(t, tp, vf.Integer(0))
	require.True(t, reflect.TypeOf(s).AssignableTo(tp.ReflectType()))
	require.Equal(t, dgo.TiSensitive, tp.TypeIdentifier())
	require.Equal(t, dgo.OpSensitive, tp.(dgo.UnaryType).Operator())
	require.Equal(t, `sensitive[int]`, s.Type().String())

	require.Equal(t, tf.Sensitive(), tf.ParseType(`sensitive`))
	require.Equal(t, tf.Sensitive(typ.Integer), tf.ParseType(`sensitive[int]`))
	require.Equal(t, vf.Sensitive(typ.Integer).Type(), tf.ParseType(`sensitive int`))
	require.Equal(t, vf.Sensitive(34).Type(), tf.ParseType(`sensitive 34`))
	require.Panic(t, func() { tf.ParseType(`sensitive[34]`) }, `illegal argument`)
	require.Panic(t, func() { tf.ParseType(`sensitive[int, string]`) }, `illegal number of arguments`)
}

func TestSensitiveType_New(t *testing.T) {
	s := vf.Sensitive(`hide me`)
	require.Equal(t, s, vf.New(typ.Sensitive, vf.Arguments(`hide me`)))
	require.Equal(t, s, vf.New(typ.Sensitive, vf.String(`hide me`)))
	require.Same(t, s, vf.New(typ.Sensitive, vf.Arguments(s)))
	require.Same(t, s, vf.New(typ.Sensitive, s))

	require.Panic(t, func() { vf.New(typ.Sensitive, vf.Arguments(`hide me`, `and me`)) }, `illegal number of arguments`)
}

func TestSensitive(t *testing.T) {
	s := vf.Sensitive(vf.String(`a`))
	require.Equal(t, s, s)
	require.Equal(t, s, vf.Sensitive(vf.String(`a`)))
	require.NotEqual(t, s, vf.Sensitive(vf.String(`b`)))
	require.NotEqual(t, s, vf.Strings(`a`))

	require.True(t, s.Frozen())
	a := vf.MutableValues(`a`)
	s = vf.Sensitive(a)
	require.False(t, s.Frozen())
	s.Freeze()
	require.True(t, s.Frozen())
	require.True(t, a.Frozen())
	require.Same(t, s.Unwrap(), a)

	a = vf.MutableValues(`a`)
	s = vf.Sensitive(a)
	c := s.FrozenCopy().(dgo.Sensitive)
	require.False(t, s.Frozen())
	require.True(t, c.Frozen())
	require.Equal(t, s.Unwrap(), c.Unwrap())
	require.NotSame(t, s.Unwrap(), c.Unwrap())

	s = vf.Sensitive(vf.MutableValues(`a`))
	require.False(t, s.Frozen())
	c = s.FrozenCopy().(dgo.Sensitive)
	require.NotSame(t, s, c)
	require.True(t, c.Frozen())
	s = c.FrozenCopy().(dgo.Sensitive)
	require.Same(t, s, c)

	c = s.ThawedCopy().(dgo.Sensitive)
	require.NotSame(t, s, c)
	require.False(t, c.Frozen()) // string is frozen regardless

	require.Equal(t, `sensitive [value redacted]`, s.String())

	require.NotEqual(t, typ.Sensitive.HashCode(), s.HashCode())
	require.Equal(t, s.HashCode(), s.HashCode())
}
