package internal_test

import (
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestBooleanDefault(t *testing.T) {
	tp := typ.Boolean
	meta := tp.Type()
	assert.Instance(t, meta, tp)
	assert.NotInstance(t, tp, tp)
	assert.NotAssignable(t, meta, tp)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, meta)
	assert.True(t, tp.IsInstance(true))
	assert.True(t, tp.IsInstance(false))
	assert.Equal(t, reflect.TypeOf(true), tp.ReflectType())
	assert.Equal(t, reflect.TypeOf(false), tp.ReflectType())
	assert.Same(t, tp, typ.Generic(tp))
}

func TestBooleanType(t *testing.T) {
	v := vf.True
	tp := v.Type().(dgo.BooleanType)
	assert.Instance(t, tp, v)
	assert.Instance(t, typ.Boolean, v)
	assert.Instance(t, typ.True, v)
	assert.Instance(t, tp, true)
	assert.Instance(t, typ.Boolean, true)
	assert.Instance(t, typ.True, true)
	assert.NotInstance(t, typ.False, v)
	assert.NotInstance(t, typ.False, true)
	assert.Assignable(t, typ.Boolean, tp)
	assert.NotAssignable(t, tp, typ.Boolean)
	assert.Equal(t, v, tp)
	assert.True(t, tp.IsInstance(true))
	assert.False(t, tp.IsInstance(false))
	assert.Equal(t, `true`, tp.String())

	v = vf.False
	tp = v.Type().(dgo.BooleanType)
	assert.Instance(t, tp, v)
	assert.Instance(t, typ.Boolean, v)
	assert.Instance(t, typ.False, v)
	assert.NotInstance(t, typ.True, v)
	assert.Assignable(t, typ.Boolean, tp)
	assert.NotAssignable(t, tp, typ.Boolean)
	assert.Equal(t, v, tp)
	assert.True(t, tp.IsInstance(false))
	assert.False(t, tp.IsInstance(true))
	assert.Equal(t, `false`, tp.String())
	assert.Equal(t, typ.Boolean.HashCode(), typ.Boolean.HashCode())
	assert.NotEqual(t, 0, typ.Boolean.HashCode())
	assert.NotEqual(t, 0, typ.True.HashCode())
	assert.NotEqual(t, 0, typ.False.HashCode())
	assert.NotEqual(t, typ.Boolean.HashCode(), typ.True.HashCode())
	assert.NotEqual(t, typ.Boolean.HashCode(), typ.False.HashCode())
	assert.NotEqual(t, typ.True.HashCode(), typ.False.HashCode())
	assert.Equal(t, `bool`, typ.Boolean.String())
	assert.Equal(t, reflect.TypeOf(false), typ.True.ReflectType())
	assert.NotSame(t, tp, typ.Generic(tp))
}

func TestNew_bool(t *testing.T) {
	assert.Equal(t, vf.True, vf.New(typ.Boolean, vf.String(`y`)))
	assert.Equal(t, vf.True, vf.New(typ.Boolean, vf.String(`Yes`)))
	assert.Equal(t, vf.True, vf.New(typ.Boolean, vf.String(`TRUE`)))
	assert.Equal(t, vf.False, vf.New(typ.Boolean, vf.String(`N`)))
	assert.Equal(t, vf.False, vf.New(typ.Boolean, vf.String(`no`)))
	assert.Equal(t, vf.False, vf.New(typ.Boolean, vf.String(`False`)))
	assert.Equal(t, vf.True, vf.New(typ.Boolean, vf.Float(1)))
	assert.Equal(t, vf.False, vf.New(typ.Boolean, vf.Float(0)))
	assert.Equal(t, vf.True, vf.New(typ.Boolean, vf.Float(1)))
	assert.Equal(t, vf.False, vf.New(typ.Boolean, vf.Int64(0)))
	assert.Equal(t, vf.True, vf.New(typ.Boolean, vf.Int64(1)))
	assert.Equal(t, vf.False, vf.New(typ.Boolean, vf.False))
	assert.Equal(t, vf.True, vf.New(typ.Boolean, vf.True))
	assert.Equal(t, vf.True, vf.New(typ.Boolean, vf.Arguments(vf.True)))
	assert.Panic(t, func() { vf.New(typ.Boolean, vf.String(`unhappy`)) }, `unable to create a bool from unhappy`)
	assert.Panic(t, func() { vf.New(typ.Boolean, vf.Arguments(vf.True, vf.True)) }, `illegal number of arguments`)
	assert.Panic(t, func() { vf.New(typ.False, vf.True) }, `the value true cannot be assigned to a variable of type false`)
}

func TestBoolean(t *testing.T) {
	assert.Equal(t, vf.True, vf.Boolean(true))
	assert.Equal(t, vf.False, vf.Boolean(false))
}

func TestBoolean_Equals(t *testing.T) {
	assert.True(t, vf.True.Equals(vf.True))
	assert.True(t, vf.True.Equals(true))
	assert.False(t, vf.True.Equals(vf.False))
	assert.False(t, vf.True.Equals(false))
	assert.True(t, vf.False.Equals(vf.False))
	assert.True(t, vf.False.Equals(false))
	assert.False(t, vf.False.Equals(vf.True))
	assert.False(t, vf.False.Equals(true))
	assert.True(t, vf.True.GoBool())
	assert.False(t, vf.False.GoBool())
}

func TestBoolean_HashCode(t *testing.T) {
	assert.NotEqual(t, 0, vf.True.HashCode())
	assert.NotEqual(t, 0, vf.False.HashCode())
	assert.Equal(t, vf.True.HashCode(), vf.True.HashCode())
	assert.NotEqual(t, vf.True.HashCode(), vf.False.HashCode())
}

func TestBoolean_CompareTo(t *testing.T) {
	c, ok := vf.True.CompareTo(vf.True)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.True.CompareTo(vf.False)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.False.CompareTo(vf.True)
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	c, ok = vf.True.CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	c, ok = vf.False.CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	_, ok = vf.True.CompareTo(vf.Int64(1))
	assert.False(t, ok)
}

func TestBoolean_ReflectTo(t *testing.T) {
	var b bool
	vf.True.ReflectTo(reflect.ValueOf(&b).Elem())
	assert.True(t, b)

	var bp *bool
	vf.True.ReflectTo(reflect.ValueOf(&bp).Elem())
	assert.True(t, *bp)

	var mi interface{}
	mip := &mi
	vf.True.ReflectTo(reflect.ValueOf(mip).Elem())
	bc, ok := mi.(bool)
	require.True(t, ok)
	assert.True(t, bc)
}

func TestBoolean_String(t *testing.T) {
	assert.Equal(t, `true`, vf.True.String())
	assert.Equal(t, `false`, vf.False.String())
}
