package internal_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/typ"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/vf"
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
	B int
}

func TestNative(t *testing.T) {
	c, ok := vf.Value(make(chan bool)).(dgo.Native)
	require.True(t, ok)
	f, ok := vf.Value(reflect.ValueOf).(dgo.Native)
	require.True(t, ok)
	s, ok := vf.Value(testStruct{`a`, 2}).(dgo.Native)
	require.True(t, ok)
	require.Equal(t, reflect.ValueOf, f)
	require.NotEqual(t, c, f)
	require.NotEqual(t, reflect.ValueOf, reflect.TypeOf)
	require.NotEqual(t, reflect.ValueOf, `ValueOf`)
	require.NotEqual(t, reflect.ValueOf, vf.Value(`ValueOf`))

	require.True(t, f.Frozen())
	require.Same(t, f, f.FrozenCopy())
	f.Freeze() // No panic

	require.False(t, c.Frozen())
	require.Panic(t, func() { c.Freeze() }, `cannot be frozen`)
	require.Panic(t, func() { c.FrozenCopy() }, `cannot be frozen`)

	require.NotEqual(t, 0, f.HashCode())
	require.Equal(t, f.HashCode(), f.HashCode())

	require.NotEqual(t, 0, c.HashCode())
	require.Equal(t, c.HashCode(), c.HashCode())

	require.NotEqual(t, 0, s.HashCode())
	require.Equal(t, s.HashCode(), s.HashCode())

	nt := f.Type()
	ct := c.Type()
	require.Assignable(t, nt, nt)
	require.Assignable(t, typ.Native, nt)
	require.NotAssignable(t, nt, typ.Native)
	require.NotAssignable(t, nt, typ.String)
	require.Instance(t, nt, f)
	require.Instance(t, typ.Native, f)
	require.NotInstance(t, typ.Native, vf.String(`ValueOf`))
	require.Instance(t, nt.Type(), nt)

	require.Equal(t, nt, nt)
	require.NotEqual(t, nt, ct)
	require.NotEqual(t, nt, typ.Any)

	require.NotEqual(t, 0, nt.HashCode())
	require.Equal(t, nt.HashCode(), nt.HashCode())

	require.Equal(t, `chan bool`, ct.String())

	// Force native of something that will never be native, just to
	// force a kind that is unknown to native
	require.Equal(t, 1234, internal.Native(reflect.ValueOf(1)).HashCode())
	require.Equal(t, 1, internal.Native(reflect.ValueOf(1)).GoValue())
}
