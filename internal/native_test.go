package internal_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
	"github.com/tada/dgo/typ"

	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/vf"
)

type teststruct int

func (t teststruct) String() string {
	return fmt.Sprintf(`test[%d]`, t)
}

func TestNative_String(t *testing.T) {
	n := vf.Value(teststruct(0))
	require.Equal(t, `test[0]`, n.String())
}

type testStruct struct {
	A string
}

type testStructNoStringer struct {
	A string
}

func (t *testStruct) String() string {
	return t.A
}

type testStructFormat struct {
	A string
}

func (t *testStructFormat) Format(s fmt.State, _ rune) {
	_, _ = s.Write([]byte(`hello from testStruct Format`))
}

func TestNative(t *testing.T) {
	c, ok := vf.Value(make(chan bool)).(dgo.Native)
	require.True(t, ok)
	s, ok := vf.Value(&testStruct{`a`}).(dgo.Native)
	require.True(t, ok)
	require.Equal(t, &testStruct{`a`}, s)
	require.NotEqual(t, c, s)
	require.NotEqual(t, &testStruct{`a`}, reflect.TypeOf)
	require.NotEqual(t, &testStruct{`a`}, `ValueOf`)
	require.NotEqual(t, &testStruct{`a`}, vf.Value(`ValueOf`))

	require.NotEqual(t, 0, s.HashCode())
	require.Equal(t, s.HashCode(), s.HashCode())

	require.NotEqual(t, 0, c.HashCode())
	require.Equal(t, c.HashCode(), c.HashCode())

	require.Equal(t, `a`, s.String())

	// Pointer to struct not equal to struct
	require.NotEqual(t, s.HashCode(), vf.Value(testStruct{`a`}).HashCode())

	s, ok = vf.Value(&testStructNoStringer{`a`}).(dgo.Native)
	require.True(t, ok)
	require.Equal(t, `<*internal_test.testStructNoStringer Value>`, s.String())
}

func TestNative_Format(t *testing.T) {
	require.Equal(t, `hello from testStruct Format`, fmt.Sprintf("%v", vf.Value(&testStructFormat{})))
}

func TestNative_Format_from_String(t *testing.T) {
	require.Equal(t, `hello from testStruct Format`, vf.Value(&testStructFormat{}).String())
}

func TestNativeType(t *testing.T) {
	c, ok := vf.Value(make(chan bool)).(dgo.Native)
	require.True(t, ok)

	ct := c.Type()
	require.Assignable(t, ct, ct)
	require.Assignable(t, typ.Native, ct)
	require.NotAssignable(t, ct, typ.Native)
	require.NotAssignable(t, ct, typ.String)
	require.Instance(t, ct, c)
	require.Instance(t, typ.Native, c)
	require.NotInstance(t, typ.Native, vf.String(`ValueOf`))
	require.Instance(t, ct.Type(), ct)

	require.Equal(t, ct, ct)
	require.NotEqual(t, ct, typ.Any)

	require.NotEqual(t, 0, ct.HashCode())
	require.Equal(t, ct.HashCode(), ct.HashCode())

	require.Equal(t, `native["chan bool"]`, ct.String())

	// Force native of something that will never be native, just to
	// force a kind that is unknown to native
	n := 1
	require.NotEqual(t, 1234, internal.Native(reflect.ValueOf(&n)).HashCode())
	require.Equal(t, 1234, internal.Native(reflect.ValueOf(n)).HashCode())
	require.Equal(t, 1, internal.Native(reflect.ValueOf(n)).GoValue())

	require.Equal(t, ct.ReflectType(), reflect.ValueOf(c.GoValue()).Type())
}

func TestNative_ReflectTo(t *testing.T) {
	var c chan bool
	v := vf.Value(make(chan bool)).(dgo.Native)
	vf.ReflectTo(v, reflect.ValueOf(&c).Elem())
	require.NotNil(t, c)
	require.Equal(t, v, c)

	var cp *chan bool
	vf.ReflectTo(v, reflect.ValueOf(&cp).Elem())
	require.NotNil(t, cp)
	require.NotNil(t, *cp)
	require.Equal(t, v, *cp)

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	cc, ok := mi.(chan bool)
	require.True(t, ok)
	require.Equal(t, v, cc)
}
