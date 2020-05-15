package internal_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestFunction(t *testing.T) {
	f, ok := vf.Value(reflect.ValueOf).(dgo.Function)
	require.True(t, ok)
	assert.Equal(t, reflect.ValueOf, f)
	assert.NotEqual(t, reflect.ValueOf, reflect.TypeOf)
	assert.NotEqual(t, reflect.ValueOf, `ValueOf`)
	assert.NotEqual(t, reflect.ValueOf, vf.Value(`ValueOf`))
	assert.NotEqual(t, 0, f.HashCode())
	assert.Equal(t, f.HashCode(), f.HashCode())
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

type structE struct {
	A string
	B int
}

func testE(s *structE) string {
	return fmt.Sprintf(`%s:%d`, s.A, s.B)
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
	assert.Equal(t, 2, len(rs))
	assert.Equal(t, s, rs[0])
	assert.Equal(t, err, rs[1])
}

func TestFunction_Call_struct(t *testing.T) {
	f, ok := vf.Value(testE).(dgo.Function)
	require.True(t, ok)
	rs := f.Call(vf.Values(vf.Map(`A`, `count`, `B`, 23)))
	assert.Equal(t, 1, len(rs))
	assert.Equal(t, `count:23`, rs[0])
}

func TestFunction_Call_nil(t *testing.T) {
	f, ok := vf.Value(testNil).(dgo.Function)
	require.True(t, ok)

	rs := f.Call(vf.Values())
	assert.Equal(t, 1, len(rs))
	assert.True(t, nil == rs[0])
}

func TestFunction_Call_panic(t *testing.T) {
	f, ok := vf.Value(testA).(dgo.Function)
	require.True(t, ok)
	assert.Panic(t, func() { f.Call(vf.Strings(`hey`)) }, `illegal number of arguments. Expected 2, got 1`)
}

func TestFunction_Call_variadic(t *testing.T) {
	f, ok := vf.Value(testB).(dgo.Function)
	require.True(t, ok)

	rs := f.Call(vf.Values(`xyz`, 1, 2, 3))
	assert.Equal(t, 1, len(rs))
	assert.Equal(t, `xyz with [1 2 3]`, rs[0])
	assert.Equal(t, `<func(string, ...int) string Value>`, f.String())
}

func TestFunction_Call_variadic_panic(t *testing.T) {
	f, ok := vf.Value(testB).(dgo.Function)
	require.True(t, ok)
	assert.Panic(t, func() { f.Call(vf.Values()) }, `illegal number of arguments. Expected at least 1, got 0`)
}

func TestFunction_GoFunc(t *testing.T) {
	f, ok := vf.Value(testC).(dgo.GoFunction)
	require.True(t, ok)
	_, ok = f.GoFunc().(func(fn func(string, ...int) string, vs ...int) string)
	assert.True(t, ok)
}

func TestFunctionType(t *testing.T) {
	ft := tf.Function(typ.Tuple, typ.Tuple)
	assert.Equal(t, tf.ParseType(`func(...any)(...any)`), ft)

	ft = tf.ParseType(`func(string, ...string)`).(dgo.FunctionType)
	assert.Assignable(t, typ.Any, ft)
	assert.NotAssignable(t, ft, typ.Any)
	assert.Assignable(t, ft, tf.ParseType(`func(string, ...string)`))
	assert.Instance(t, ft.Type(), ft)
	assert.Panic(t, func() { ft.ReflectType() }, `unable to build`)
	assert.Equal(t, ft, ft)
	assert.Equal(t, ft, tf.ParseType(`func(string, ...string)`))
	assert.Equal(t, ft.String(), `func(string,...string)`)
	assert.NotEqual(t, ft, tf.ParseType(`func(int, ...string)`))
	assert.NotEqual(t, ft, tf.ParseType(`func(string, []string)`))
	assert.NotEqual(t, ft, tf.ParseType(`func(string, ...string) string`))
	assert.NotEqual(t, ft, tf.ParseType(`{string, ...string}`))

	in := ft.In()
	assert.Assignable(t, in, tf.ParseType(`{string, ...string}`))
	assert.Assignable(t, in, tf.ParseType(`[1]string`))
	assert.NotAssignable(t, in, tf.ParseType(`[]string`))
	assert.Same(t, typ.String, in.ElementType())
	assert.Same(t, typ.String, tf.ParseType(`{...string}`).(dgo.TupleType).ElementType())
	assert.NotEqual(t, 0, ft.HashCode())
	assert.Equal(t, ft.HashCode(), tf.ParseType(`func(string, ...string)`).HashCode())

	ft = tf.ParseType(`func(string, ...int) string`).(dgo.FunctionType)
	assert.Instance(t, ft, testB)
	assert.NotInstance(t, ft, testA)
	assert.NotInstance(t, ft, ft)
	assert.Panic(t, func() { tf.Function(nil, tf.VariadicTuple(tf.Array(typ.Any))) },
		`return values cannot be variadic`)
}

func TestFunctionType_exact(t *testing.T) {
	ft := vf.Value(testA).Type().(dgo.FunctionType)
	assert.Equal(t, `func(string,error) (string,error)`, ft.String())

	ft = vf.Value(testB).Type().(dgo.FunctionType)
	assert.Equal(t, `func(string,...int) string`, ft.String())
	assert.Assignable(t, ft, ft)
	assert.Instance(t, ft, testB)
	assert.NotInstance(t, ft, testC)
	assert.NotInstance(t, ft, `func(string,...int) string`)
	assert.Instance(t, ft.Type(), ft)
	assert.NotEqual(t, 0, ft.HashCode())
	assert.Equal(t, tf.ParseType(`func(string,...int) string`).HashCode(), ft.HashCode())
	assert.Equal(t, ft.ReflectType(), reflect.TypeOf(testB))

	in := ft.In()
	assert.Equal(t, `{string,...int}`, in.String())
	assert.Equal(t, in, in)
	assert.Equal(t, tf.ParseType(`{string, ...int}`), in)
	assert.Equal(t, in, tf.ParseType(`{string, ...int}`))
	assert.NotEqual(t, in, tf.ParseType(`{string, int, ...int}`))
	assert.Assignable(t, in, tf.ParseType(`{string, int, ...int}`))
	assert.NotAssignable(t, in, tf.ParseType(`{string, string, ...int}`))
	assert.Equal(t, tf.ParseType(`string&int`), in.ElementType())
	assert.Instance(t, in, vf.Values(`hello`))
	assert.Instance(t, in, vf.Values(`hello`, 3))
	assert.Instance(t, in, vf.Values(`hello`, 3, 3))
	assert.Instance(t, in.Type(), in)
	assert.False(t, in.Unbounded())
	assert.Equal(t, in.HashCode(), tf.ParseType(`{string, ...int}`).HashCode())

	out := ft.Out()
	assert.Equal(t, out, tf.ParseType(`{string}`))
	assert.Equal(t, out.ReflectType(), reflect.SliceOf(reflect.TypeOf(``)))
	assert.Instance(t, out.Type(), out)

	ft = vf.Value(testC).Type().(dgo.FunctionType)
	assert.Equal(t, `func(func(string,...int) string,...int) string`, ft.String())
	assert.NotEqual(t, 0, ft.HashCode())
	assert.Equal(t, tf.ParseType(`func(func(string, ...int) string,...int) string`).HashCode(), ft.HashCode())
	assert.Equal(t, tf.Function(ft.In(), tf.Tuple(typ.String)).HashCode(), ft.HashCode())
	assert.Assignable(t, ft, tf.ParseType(`func(func(string, int, ...int) string, ...int) string`))
	assert.Equal(t, ft, tf.ParseType(`func(func(string, ...int) string, ...int) string`))
	assert.NotEqual(t, ft, tf.ParseType(`func(func(string, int, ...int) string, ...int) string`))
	assert.NotEqual(t, ft, tf.ParseType(`{func(string, int, ...int) string, ...int}`))

	ft = vf.Value(testD).Type().(dgo.FunctionType)
	assert.True(t, ft.In().Unbounded())
}
