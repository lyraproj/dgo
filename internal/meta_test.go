package internal_test

import (
	"reflect"
	"testing"

	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/vf"

	"github.com/tada/dgo/dgo"
	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/typ"
)

func TestMeta(t *testing.T) {
	tp := typ.Any.Type()
	require.Instance(t, tp, typ.Any)
	require.Instance(t, tp, typ.Boolean)
	require.NotInstance(t, tp, 3)
	require.Same(t, typ.Any, tp.(dgo.Meta).Describes())
	require.Equal(t, tf.Meta(typ.Any), tp)

	require.Equal(t, typ.Any.Type(), tp)
	require.NotEqual(t, typ.Boolean.Type(), tp)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Assignable(t, tp, typ.Any.Type())
	require.NotAssignable(t, tp, typ.Boolean.Type())
	require.Equal(t, dgo.OpMeta, tp.(dgo.UnaryType).Operator())
	require.Equal(t, `type`, tp.String())
	require.Equal(t, `type[int]`, typ.Integer.Type().String())

	require.True(t, reflect.ValueOf(typ.Any).Type().AssignableTo(tp.ReflectType()))
}

func TestMetaType_New(t *testing.T) {
	m := tf.String(1, 10)
	require.Same(t, m, vf.New(typ.String.Type(), m))
	require.Same(t, m, vf.New(m.Type(), m))
	require.Same(t, m, vf.New(m.Type(), vf.Arguments(m)))
	require.Equal(t, m, vf.New(typ.String.Type(), vf.String(`string[1,10]`)))

	require.Equal(t, vf.String("x").Type(), vf.New(typ.String.Type(), vf.String(`"x"`)))

	require.Panic(t, func() { vf.New(typ.Integer.Type(), vf.String(`string[1, 10]`)) }, `cannot be assigned`)
}

func TestMetaMeta(t *testing.T) {
	mt := typ.Any.Type().Type()
	require.Equal(t, mt, typ.Boolean.Type().Type())
	require.NotEqual(t, mt, typ.Any)
	require.Equal(t, mt.Type(), mt)

	require.Instance(t, mt, typ.Any.Type())
	require.Instance(t, mt, typ.Boolean.Type())
	require.NotInstance(t, mt, typ.Any)

	require.Assignable(t, mt, mt)
	require.NotAssignable(t, mt, typ.Any.Type())

	require.Equal(t, dgo.OpMeta, mt.(dgo.UnaryType).Operator())
	require.Equal(t, `type[type]`, mt.String())
}
