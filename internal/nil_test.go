package internal_test

import (
	"errors"
	"reflect"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestNil_CompareTo(t *testing.T) {
	c, ok := vf.Nil.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.Nil.CompareTo(vf.Integer(1))
	require.True(t, ok)
	require.Equal(t, -1, c)
}

func TestNil_ReflectTo(t *testing.T) {
	err := errors.New(`some error`)
	vf.Nil.ReflectTo(reflect.ValueOf(&err).Elem())
	require.True(t, err == nil)

	err = errors.New(`some error`)
	ep := &err
	vf.Nil.ReflectTo(reflect.ValueOf(&ep).Elem())
	require.True(t, ep == nil)

	var mi interface{} = err
	mip := &mi
	vf.Nil.ReflectTo(reflect.ValueOf(mip).Elem())
	require.True(t, mi == nil)
}

func TestNilType(t *testing.T) {
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

	require.Equal(t, nt.ReflectType(), typ.Any.ReflectType())
}
