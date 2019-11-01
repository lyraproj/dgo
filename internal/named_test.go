package internal_test

import (
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/vf"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/typ"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
)

type testNamed int

func (a testNamed) String() string {
	return a.Type().(dgo.NamedType).ValueString(a)
}

func (a testNamed) Type() dgo.Type {
	return newtype.ExactNamed(newtype.Named(`testNamed`), a)
}

func (a testNamed) Equals(other interface{}) bool {
	return a == other
}

func (a testNamed) HashCode() int {
	return int(a)
}

type testNamedB int
type testNamedC int

type testNamedDummy interface {
	Dummy()
}

func (testNamed) Dummy() {
}

func (testNamedB) Dummy() {
}

func TestNamedType(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	tp := newtype.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil)
	require.Equal(t, tp, tp)
	require.Equal(t, tp.Name(), `testNamed`)
	require.Equal(t, tp.String(), `testNamed`)
	require.NotEqual(t, tp, `testNamed`)

	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.Any)
	require.Instance(t, tp.Type(), tp)
	require.Instance(t, tp, testNamed(0))

	require.NotEqual(t, 0, tp.HashCode())
	require.Equal(t, tp.HashCode(), tp.HashCode())
}

func TestNamedType_redefined(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	newtype.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil)
	require.Panic(t, func() {
		newtype.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamedB(0)), nil)
	}, `attempt to redefine named type 'testNamed'`)
}

func TestNamedTypeFromReflected(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	tp := newtype.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil)
	require.Same(t, tp, newtype.NamedFromReflected(reflect.TypeOf(testNamed(0))))
	require.Nil(t, newtype.NamedFromReflected(reflect.TypeOf(testNamedC(0))))
}

func TestNamedType_Assignable(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	defer newtype.RemoveNamed(`testNamedB`)
	defer newtype.RemoveNamed(`testNamedC`)
	tp := newtype.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), reflect.TypeOf((*testNamedDummy)(nil)).Elem())
	require.Assignable(t, tp, newtype.NewNamed(`testNamedB`, nil, nil, reflect.TypeOf(testNamedB(0)), nil))
	require.NotAssignable(t, tp, newtype.NewNamed(`testNamedC`, nil, nil, reflect.TypeOf(testNamedC(0)), nil))
	require.NotAssignable(t, newtype.Named(`testNamedB`), newtype.Named(`testNamedC`))
}

func TestNamedType_New(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	tp := newtype.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil)

	v := tp.New(vf.Integer(3))
	require.Equal(t, v, testNamed(3))
	require.Equal(t, 3, tp.ExtractInitArg(v))
}

func TestNamedType_New_notApplicable(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	tp := newtype.NewNamed(`testNamed`, nil, nil, reflect.TypeOf(testNamed(0)), nil)

	require.Panic(t, func() { tp.New(vf.Integer(3)) }, `creating new instances of testNamed is not possible`)
	require.Panic(t, func() { tp.ExtractInitArg(testNamed(0)) }, `creating new instances of testNamed is not possible`)
}

func TestNamedType_ValueString(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	tp := newtype.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil)

	v := tp.New(vf.Integer(3))
	require.Equal(t, `testNamed 3`, tp.ValueString(v))
}

func TestNamedType_parse(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	tp := newtype.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil)
	require.Same(t, newtype.Parse(`testNamed`), tp)
}

func TestNamedType_exact(t *testing.T) {
	defer newtype.RemoveNamed(`testNamed`)
	tp := newtype.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		return testNamed(arg.(dgo.Integer).GoInt())
	}, func(value dgo.Value) dgo.Value {
		return vf.Integer(int64(value.(testNamed)))
	}, reflect.TypeOf(testNamed(0)), nil)

	v := tp.New(vf.Integer(3))
	et := v.Type()
	require.Same(t, tp, typ.Generic(et))
	require.Assignable(t, tp, et)
	require.NotAssignable(t, et, tp)
	require.NotAssignable(t, et, tp.New(vf.Integer(4)).Type())
	require.Instance(t, et, v)
	require.NotInstance(t, et, tp.New(vf.Integer(4)))
	require.Equal(t, `testNamed 3`, et.String())
	require.Equal(t, et, newtype.Parse(`testNamed 3`))

	require.Instance(t, et.Type(), et)
	require.Instance(t, tp.Type(), et)
	require.NotInstance(t, et.Type(), tp)

	require.NotEqual(t, tp, et)
	require.Equal(t, et, tp.New(vf.Integer(3)).Type())
	require.NotEqual(t, et, tp.New(vf.Integer(3)))
	require.NotEqual(t, et, tp.New(vf.Integer(4)).Type())
	require.NotEqual(t, tp.HashCode(), et.HashCode())
}
