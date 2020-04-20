package internal_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/tada/dgo/typ"

	"github.com/tada/dgo/tf"

	"github.com/tada/dgo/dgo"
	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/vf"
)

func TestFunction(t *testing.T) {
	f, ok := vf.Value(reflect.ValueOf).(dgo.Function)
	require.True(t, ok)
	require.Equal(t, reflect.ValueOf, f)
	require.NotEqual(t, reflect.ValueOf, reflect.TypeOf)
	require.NotEqual(t, reflect.ValueOf, `ValueOf`)
	require.NotEqual(t, reflect.ValueOf, vf.Value(`ValueOf`))

	require.NotEqual(t, 0, f.HashCode())
	require.Equal(t, f.HashCode(), f.HashCode())
}

func testA(m string, e error) (string, error) {
	return m, e
}

func testB(m string, vs ...int) string {
	return fmt.Sprintf(`%s with %v`, m, vs)
}

func testC(m func(string, ...int) string, vs ...int) string {
	return m(`b`, vs...)
}

func testD(vs ...int) string {
	return fmt.Sprintf(`%v`, vs)
}

func testNil() []string {
	return nil
}

func TestFunction_Call(t *testing.T) {
	f, ok := vf.Value(testA).(dgo.Function)
	require.True(t, ok)

	err := errors.New(`if this is it`)
	s := `please let me know`
	rs := f.Call(vf.Values(s, err))
	require.Equal(t, 2, len(rs))
	require.Equal(t, s, rs[0])
	require.Equal(t, err, rs[1])
}

func TestFunction_Call_nil(t *testing.T) {
	f, ok := vf.Value(testNil).(dgo.Function)
	require.True(t, ok)

	rs := f.Call(vf.Values())
	require.Equal(t, 1, len(rs))
	require.True(t, nil == rs[0])
}

func TestFunction_Call_panic(t *testing.T) {
	f, ok := vf.Value(testA).(dgo.Function)
	require.True(t, ok)

	require.Panic(t, func() { f.Call(vf.Strings(`hey`)) }, `illegal number of arguments. Expected 2, got 1`)
}

func TestFunction_Call_variadic(t *testing.T) {
	f, ok := vf.Value(testB).(dgo.Function)
	require.True(t, ok)

	rs := f.Call(vf.Values(`xyz`, 1, 2, 3))
	require.Equal(t, 1, len(rs))
	require.Equal(t, `xyz with [1 2 3]`, rs[0])

	require.Equal(t, `<func(string, ...int) string Value>`, f.String())
}

func TestFunction_Call_variadic_panic(t *testing.T) {
	f, ok := vf.Value(testB).(dgo.Function)
	require.True(t, ok)

	require.Panic(t, func() { f.Call(vf.Values()) }, `illegal number of arguments. Expected at least 1, got 0`)
}

func TestFunction_GoFunc(t *testing.T) {
	f, ok := vf.Value(testC).(dgo.GoFunction)
	require.True(t, ok)
	_, ok = f.GoFunc().(func(fn func(string, ...int) string, vs ...int) string)
	require.True(t, ok)
}

func TestFunctionType(t *testing.T) {
	ft := tf.Function(typ.Tuple, typ.Tuple)
	require.Equal(t, tf.ParseType(`func(...any)(...any)`), ft)

	ft = tf.ParseType(`func(string, ...string)`).(dgo.FunctionType)
	require.Assignable(t, typ.Any, ft)
	require.NotAssignable(t, ft, typ.Any)
	require.Assignable(t, ft, tf.ParseType(`func(string, ...string)`))
	require.Instance(t, ft.Type(), ft)

	require.Panic(t, func() { ft.ReflectType() }, `unable to build`)
	require.Equal(t, ft, ft)
	require.Equal(t, ft, tf.ParseType(`func(string, ...string)`))
	require.Equal(t, ft.String(), `func(string,...string)`)
	require.NotEqual(t, ft, tf.ParseType(`func(int, ...string)`))
	require.NotEqual(t, ft, tf.ParseType(`func(string, []string)`))
	require.NotEqual(t, ft, tf.ParseType(`func(string, ...string) string`))
	require.NotEqual(t, ft, tf.ParseType(`{string, ...string}`))

	in := ft.In()
	require.Assignable(t, in, tf.ParseType(`{string, ...string}`))
	require.Assignable(t, in, tf.ParseType(`[1]string`))
	require.NotAssignable(t, in, tf.ParseType(`[]string`))
	require.Same(t, typ.String, in.ElementType())
	require.Same(t, typ.String, tf.ParseType(`{...string}`).(dgo.TupleType).ElementType())

	require.NotEqual(t, 0, ft.HashCode())
	require.Equal(t, ft.HashCode(), tf.ParseType(`func(string, ...string)`).HashCode())

	ft = tf.ParseType(`func(string, ...int) string`).(dgo.FunctionType)
	require.Instance(t, ft, testB)
	require.NotInstance(t, ft, testA)
	require.NotInstance(t, ft, ft)

	require.Panic(t, func() { tf.Function(nil, tf.VariadicTuple(tf.Array(typ.Any))) },
		`return values cannot be variadic`)
}

func TestFunctionType_exact(t *testing.T) {
	ft := vf.Value(testA).Type().(dgo.FunctionType)
	require.Equal(t, `func(string,error) (string,error)`, ft.String())

	ft = vf.Value(testB).Type().(dgo.FunctionType)
	require.Equal(t, `func(string,...int) string`, ft.String())
	require.Assignable(t, ft, ft)
	require.Instance(t, ft, testB)
	require.NotInstance(t, ft, testC)
	require.NotInstance(t, ft, `func(string,...int) string`)
	require.Instance(t, ft.Type(), ft)
	require.NotEqual(t, 0, ft.HashCode())
	require.Equal(t, tf.ParseType(`func(string,...int) string`).HashCode(), ft.HashCode())
	require.Equal(t, ft.ReflectType(), reflect.TypeOf(testB))

	in := ft.In()
	require.Equal(t, `{string,...int}`, in.String())
	require.Equal(t, in, in)
	require.Equal(t, tf.ParseType(`{string, ...int}`), in)
	require.Equal(t, in, tf.ParseType(`{string, ...int}`))
	require.NotEqual(t, in, tf.ParseType(`{string, int, ...int}`))
	require.Assignable(t, in, tf.ParseType(`{string, int, ...int}`))
	require.NotAssignable(t, in, tf.ParseType(`{string, string, ...int}`))
	require.Equal(t, tf.ParseType(`string&int`), in.ElementType())

	require.Instance(t, in, vf.Values(`hello`))
	require.Instance(t, in, vf.Values(`hello`, 3))
	require.Instance(t, in, vf.Values(`hello`, 3, 3))
	require.Instance(t, in.Type(), in)

	require.False(t, in.Unbounded())
	require.Equal(t, in.HashCode(), tf.ParseType(`{string, ...int}`).HashCode())

	out := ft.Out()
	require.Equal(t, out, tf.ParseType(`{string}`))
	require.Equal(t, out.ReflectType(), reflect.SliceOf(reflect.TypeOf(``)))
	require.Instance(t, out.Type(), out)

	ft = vf.Value(testC).Type().(dgo.FunctionType)
	require.Equal(t, `func(func(string,...int) string,...int) string`, ft.String())
	require.NotEqual(t, 0, ft.HashCode())
	require.Equal(t, tf.ParseType(`func(func(string, ...int) string,...int) string`).HashCode(), ft.HashCode())
	require.Equal(t, tf.Function(ft.In(), tf.Tuple(typ.String)).HashCode(), ft.HashCode())
	require.Assignable(t, ft, tf.ParseType(`func(func(string, int, ...int) string, ...int) string`))

	require.Equal(t, ft, tf.ParseType(`func(func(string, ...int) string, ...int) string`))
	require.NotEqual(t, ft, tf.ParseType(`func(func(string, int, ...int) string, ...int) string`))
	require.NotEqual(t, ft, tf.ParseType(`{func(string, int, ...int) string, ...int}`))

	ft = vf.Value(testD).Type().(dgo.FunctionType)
	require.True(t, ft.In().Unbounded())
}
