package internal_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tada/dgo/tf"

	"github.com/tada/dgo/test/require"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

type testint int

type testuint uint

type testfloat float64

type testbool bool

func (t testint) String() string {
	return fmt.Sprintf(`int[%d]`, t)
}

func (t testuint) String() string {
	return fmt.Sprintf(`uint[%d]`, t)
}

func (t testfloat) String() string {
	return fmt.Sprintf(`float[%g]`, t)
}

func (t testbool) String() string {
	return fmt.Sprintf(`bool[%t]`, t)
}

func TestNative_String(t *testing.T) {
	assert.Equal(t, `int[0]`, vf.Value(testint(0)).String())
	assert.Equal(t, `uint[0]`, vf.Value(testuint(0)).String())
	assert.Equal(t, `float[0]`, vf.Value(testfloat(0)).String())
	assert.Equal(t, `bool[true]`, vf.Value(testbool(true)).String())

	b := testbool(false)
	bp := &b
	assert.Equal(t, `bool[false]`, vf.Value(b).String())
	assert.Equal(t, `bool[false]`, vf.Value(bp).String())
	assert.Equal(t, `bool[false]`, vf.Value(&bp).String())
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
	assert.Equal(t, &testStruct{`a`}, s)
	assert.Equal(t, reflect.TypeOf(&testStruct{`a`}), s.(dgo.NativeType).GoType())
	assert.NotEqual(t, c, s)
	assert.NotEqual(t, &testStruct{`a`}, reflect.TypeOf)
	assert.NotEqual(t, &testStruct{`a`}, `ValueOf`)
	assert.NotEqual(t, &testStruct{`a`}, vf.Value(`ValueOf`))
	assert.NotEqual(t, 0, s.HashCode())
	assert.Equal(t, s.HashCode(), s.HashCode())
	assert.NotEqual(t, 0, c.HashCode())
	assert.Equal(t, c.HashCode(), c.HashCode())
	assert.Equal(t, `a`, s.String())

	// Pointer to struct not equal to struct
	assert.NotEqual(t, s.HashCode(), vf.Value(testStruct{`a`}).HashCode())

	s, ok = vf.Value(&testStructNoStringer{`a`}).(dgo.Native)
	require.True(t, ok)
	assert.Equal(t, `&internal_test.testStructNoStringer{A:"a"}`, s.String())
}

func TestNative_Format(t *testing.T) {
	assert.Equal(t, `hello from testStruct Format`, fmt.Sprintf("%v", vf.Value(&testStructFormat{})))
}

func TestNative_Format_from_String(t *testing.T) {
	assert.Equal(t, `hello from testStruct Format`, vf.Value(&testStructFormat{}).String())
}

func TestNativeType(t *testing.T) {
	c, ok := vf.Value(make(chan bool)).(dgo.Native)
	require.True(t, ok)

	c2, ok := vf.Value(make(chan struct{})).(dgo.Native)
	require.True(t, ok)

	ct := c.Type()
	ct2 := c2.Type()
	assert.Assignable(t, ct, ct)
	assert.Assignable(t, typ.Native, ct)
	assert.Assignable(t, typ.Native, tf.AnyOf(ct, ct2))
	assert.Assignable(t, typ.Generic(ct), ct)
	assert.NotAssignable(t, typ.Generic(ct).Type(), ct)
	assert.NotAssignable(t, ct, typ.Native)
	assert.NotAssignable(t, typ.Generic(ct), typ.Native)
	assert.NotAssignable(t, ct, typ.Generic(ct))
	assert.NotAssignable(t, ct, typ.String)
	assert.Instance(t, ct, c)
	assert.Instance(t, typ.Generic(ct), c)
	assert.Instance(t, typ.Native, c)
	assert.NotInstance(t, typ.Generic(ct), c2)
	assert.NotInstance(t, typ.Native, vf.String(`ValueOf`))
	assert.Instance(t, ct.Type(), ct)
	assert.Equal(t, ct, ct)
	assert.NotEqual(t, ct, ct2)
	assert.NotEqual(t, typ.Generic(ct), typ.Generic(ct2))
	assert.NotEqual(t, ct, typ.Any)
	assert.NotEqual(t, typ.Generic(ct), ct)
	assert.NotEqual(t, 0, ct.HashCode())
	assert.Equal(t, ct.HashCode(), ct.HashCode())
	assert.Equal(t, `native["chan bool"]`, typ.Generic(ct).String())

	// Force native of something that will never be native, just to
	// force a kind that is unknown to native
	n := 1
	assert.NotEqual(t, 1234, internal.Native(reflect.ValueOf(&n)).HashCode())
	assert.Equal(t, 1, internal.Native(reflect.ValueOf(n)).HashCode())
	assert.Equal(t, 1, internal.Native(reflect.ValueOf(n)).GoValue())
	assert.Equal(t, ct.ReflectType(), reflect.ValueOf(c.GoValue()).Type())
}

func TestNativeType_ReflectType(t *testing.T) {
	c, ok := vf.Value(make(chan bool)).(dgo.Native)
	require.True(t, ok)
	ct := typ.Generic(c.Type())
	assert.Equal(t, c.Type().ReflectType(), ct.ReflectType())
}

func TestNativeType_String(t *testing.T) {
	type X struct {
		A string
	}
	ct := vf.Value(&X{A: `a`}).Type()
	assert.Equal(t, `&internal_test.X{A:"a"}`, ct.String())
	assert.Equal(t, `native["*internal_test.X"]`, typ.Generic(ct).String())
}

func TestNativeType_CantInterface(t *testing.T) {
	ct := vf.Value(reflect.ValueOf(t).Elem().Field(0)).Type()
	assert.Equal(t, `native["testing.common"]`, ct.String())
	assert.Equal(t, `native["testing.common"]`, fmt.Sprintf(`%#v`, ct))
}

func TestNative_ReflectTo(t *testing.T) {
	var c chan bool
	v := vf.Value(make(chan bool)).(dgo.Native)
	vf.ReflectTo(v, reflect.ValueOf(&c).Elem())
	assert.NotNil(t, c)
	assert.Equal(t, v, c)

	var cp *chan bool
	vf.ReflectTo(v, reflect.ValueOf(&cp).Elem())
	assert.NotNil(t, cp)
	assert.NotNil(t, *cp)
	assert.Equal(t, v, *cp)

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	cc, ok := mi.(chan bool)
	require.True(t, ok)
	assert.Equal(t, v, cc)
}
