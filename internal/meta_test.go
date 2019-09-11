package internal_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
)

func TestMeta(t *testing.T) {
	tp := typ.Any.Type()
	require.Instance(t, tp, typ.Any)
	require.Instance(t, tp, typ.Boolean)
	require.NotInstance(t, tp, 3)

	require.Equal(t, typ.Any.Type(), tp)
	require.NotEqual(t, typ.Boolean.Type(), tp)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Assignable(t, tp, typ.Any.Type())
	require.NotAssignable(t, tp, typ.Boolean.Type())
	require.Equal(t, dgo.OpMeta, tp.(dgo.UnaryType).Operator())
	require.Equal(t, `type`, tp.String())
	require.Equal(t, `type[int]`, typ.Integer.Type().String())
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
