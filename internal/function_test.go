package internal_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/vf"
)

func TestFunction(t *testing.T) {
	f, ok := vf.Value(reflect.ValueOf).(dgo.Callable)
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

func TestFunction_Call(t *testing.T) {
	f, ok := vf.Value(testA).(dgo.Callable)
	require.True(t, ok)

	err := errors.New(`if this is it`)
	s := `please let me know`
	rs := f.Call(s, err)
	require.Equal(t, 2, len(rs))
	require.Equal(t, s, rs[0])
	require.Equal(t, err, rs[1])
}

func TestFunction_Call_panic(t *testing.T) {
	f, ok := vf.Value(testA).(dgo.Callable)
	require.True(t, ok)

	require.Panic(t, func() { f.Call(`hey`) }, `illegal number of arguments. Expected 2, got 1`)
}

func TestFunction_Call_variadic(t *testing.T) {
	f, ok := vf.Value(testB).(dgo.Callable)
	require.True(t, ok)

	rs := f.Call(`xyz`, 1, 2, 3)
	require.Equal(t, 1, len(rs))
	require.Equal(t, `xyz with [1 2 3]`, rs[0])

	require.Equal(t, `<func(string, ...int) string Value>`, f.String())
	require.Equal(t, `native["func(string, ...int) string"]`, f.Type().String())
}

func TestFunction_Call_variadic_panic(t *testing.T) {
	f, ok := vf.Value(testB).(dgo.Callable)
	require.True(t, ok)

	require.Panic(t, func() { f.Call() }, `illegal number of arguments. Expected at least 1, got 0`)
}
