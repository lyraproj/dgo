package internal_test

import (
	"bytes"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/tf"

	"github.com/lyraproj/dgo/typ"

	require "github.com/lyraproj/dgo/dgo_test"

	"github.com/lyraproj/dgo/vf"
)

func TestBinaryType(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.Binary(bs, false)
	tp := v.Type().(dgo.BinaryType)
	require.Assignable(t, typ.Binary, tp)
	require.NotAssignable(t, tp, typ.Binary)
	require.NotAssignable(t, tp, typ.String)
	require.Instance(t, typ.Binary, bs)
	require.Instance(t, typ.Binary, v)
	require.Instance(t, tp, v)
	require.Instance(t, tp, []byte{1, 2, 3})
	require.NotInstance(t, tp, []byte{1, 2})
	require.NotInstance(t, tp, []byte{1, 2})
	require.NotInstance(t, tf.Binary(1, 5), `abc`)

	require.Assignable(t, tf.Binary(3, 3), tp)
	require.Assignable(t, tp, tf.Binary(3, 3))

	require.Same(t, typ.Binary, tf.Binary())
	require.Same(t, typ.Binary, tf.Binary(0, math.MaxInt64))

	require.NotAssignable(t, tf.Binary(4), tp)
	require.NotAssignable(t, tf.Binary(0, 2), tp)

	require.Equal(t, tf.Binary(1, 2), tf.Binary(2, 1))
	require.Equal(t, tf.Binary(-1, 2), tf.Binary(0, 2))
	require.NotEqual(t, tf.Binary(0, 2), tf.String(0, 2))

	require.Equal(t, tf.Binary(0, 2).HashCode(), tf.Binary(0, 2).HashCode())
	require.NotEqual(t, tf.Binary(1, 2).HashCode(), tf.Binary(0, 2).HashCode())
	require.NotEqual(t, tf.Binary(0, 1).HashCode(), tf.Binary(0, 2).HashCode())

	require.Instance(t, tp.Type(), tp)

	require.Panic(t, func() { tf.Binary(`blue`) }, `illegal argument`)
	require.Panic(t, func() { tf.Binary(`blue`, 1) }, `illegal argument 1`)
	require.Panic(t, func() { tf.Binary(1, `blue`) }, `illegal argument 2`)
	require.Panic(t, func() { tf.Binary(1, 2, 3) }, `illegal number of arguments`)
	require.Equal(t, `binary[3,3]`, tp.String())

	require.Equal(t, reflect.TypeOf([]byte{}), tp.ReflectType())
}

func TestBinaryType_New(t *testing.T) {
	b := vf.New(typ.Binary, vf.Arguments(vf.Values(1, 2, 3)))
	require.Equal(t, vf.Binary([]byte{1, 2, 3}, true), b)
}

func TestBinaryType_New_passThrough(t *testing.T) {
	b := vf.Binary([]byte{1, 2, 3}, true)
	require.Same(t, b, vf.New(typ.Binary, b))
}

func TestBinaryType_New_stringRaw(t *testing.T) {
	require.Equal(t, vf.Binary([]byte(`hello`), true), vf.New(typ.Binary, vf.Arguments(vf.String(`hello`), `%r`)))
}

func TestBinaryType_New_badCount(t *testing.T) {
	require.Panic(t, func() { vf.New(typ.Binary, vf.Arguments(1, 2, 3)) }, `illegal number of arguments`)
}

func TestBinaryType_New_badArg(t *testing.T) {
	require.Panic(t, func() { vf.New(typ.Binary, vf.Integer(3)) }, `illegal argument for binary`)
}

func TestBinaryType_New_badType(t *testing.T) {
	require.Panic(t, func() { vf.New(tf.Binary(2, 2), vf.New(typ.Binary, vf.Value([]byte{1, 2, 3}))) }, `cannot be assigned`)
}

func TestBinaryType_New_badBytes(t *testing.T) {
	require.Panic(t, func() { vf.New(typ.Binary, vf.Values(1, 311, 3)) }, `the value 311 cannot be assigned to a variable of type 0..255`)
}

func TestBinaryType_New_badFormat(t *testing.T) {
	require.Panic(t, func() { vf.New(typ.Binary, vf.Arguments(`hello`, `%x`)) }, `illegal argument`)
}

type badReader int

func (badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New(`oops`)
}

func TestBinary(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.Binary(bs, false)
	require.Equal(t, `AQID`, v.String())

	// Test mutability
	bs[1] = 4
	require.Equal(t, reflect.ValueOf(bs).Pointer(), reflect.ValueOf(v.GoBytes()).Pointer())
	require.Equal(t, `AQQD`, v.String())

	// Test immutability
	v = vf.Binary(bs, true)
	bs[1] = 2
	require.Equal(t, `AQQD`, v.String())
	require.NotEqual(t, reflect.ValueOf(bs).Pointer(), reflect.ValueOf(v.GoBytes()).Pointer())

	// Test panic on bad read
	require.Panic(t, func() { vf.BinaryFromData(badReader(0)) }, `oops`)
}

func TestBinaryString(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.BinaryFromString(`AQID`)
	require.True(t, bytes.Equal(bs, v.GoBytes()))

	require.Panic(t, func() { vf.BinaryFromString(`----`) }, `illegal base64`)
}

func TestBinaryFromEncoded(t *testing.T) {
	bs := []byte(`hello`)
	notStrict := `aGVsbG9=`
	v := vf.BinaryFromEncoded(notStrict, `%b`)
	require.True(t, bytes.Equal(bs, v.GoBytes()))
	require.Panic(t, func() { vf.BinaryFromEncoded(notStrict, `%B`) }, `illegal base64 data at input byte 7`)

	bs = []byte{'r', 0x82, 0xff}
	notUtf8 := string(bs)
	v = vf.BinaryFromEncoded(notUtf8, `%r`)
	require.True(t, bytes.Equal(bs, v.GoBytes()))
	require.Panic(t, func() { vf.BinaryFromEncoded(notUtf8, `%s`) }, `Expected valid utf8 string`)

	bs = []byte{3, 240, 126}
	require.True(t, bytes.Equal(bs, vf.BinaryFromEncoded(`A/B+`, `%B`).GoBytes()))
	require.True(t, bytes.Equal(bs, vf.BinaryFromEncoded(`A_B-`, `%u`).GoBytes()))
	require.Panic(t, func() { vf.BinaryFromEncoded(`A/B+`, `%u`) }, `illegal base64 data at input byte 1`)

	require.Panic(t, func() { vf.BinaryFromEncoded(`A/B+`, `%x`) }, `Expected one of the supported format specifiers`)
}

func TestBinary_CompareTo(t *testing.T) {
	a := vf.Binary([]byte{'a', 'b', 'c'}, false)

	c, ok := a.CompareTo(a)
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = a.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	_, ok = a.CompareTo(vf.Value('a'))
	require.False(t, ok)

	b := vf.Binary([]byte{'a', 'b', 'c'}, false)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 0, c)

	b = vf.Binary([]byte{'a', 'b', 'c', 'd'}, false)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, -1, c)

	b = vf.Binary([]byte{'a', 'b', 'd', 'd'}, false)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, -1, c)

	b = vf.Binary([]byte{'a', 'b'}, false)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 1, c)

	b = vf.Binary([]byte{'a', 'b', 'd'}, false)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, -1, c)

	b = vf.Binary([]byte{'a', 'b', 'b'}, false)
	c, ok = a.CompareTo(b)
	require.True(t, ok)
	require.Equal(t, 1, c)
}

func TestBinary_Copy(t *testing.T) {
	a := vf.Binary([]byte{'a', 'b'}, true)
	require.Same(t, a, a.Copy(true))

	c := a.Copy(false)
	require.False(t, c.Frozen())
	require.NotSame(t, c, c.Copy(false))

	c = c.Copy(true)
	require.True(t, c.Frozen())
	require.Same(t, c, c.Copy(true))
}

func TestBinary_Equal(t *testing.T) {
	a := vf.Binary([]byte{1, 2}, true)
	require.True(t, a.Equals(a))

	b := vf.Binary([]byte{1}, true)
	require.False(t, a.Equals(b))
	b = vf.Binary([]byte{1, 2}, true)
	require.True(t, a.Equals(b))
	require.True(t, a.Equals([]byte{1, 2}))
	require.False(t, a.Equals(`12`))
	require.Equal(t, a.HashCode(), b.HashCode())
	require.True(t, a.Equals(vf.Value([]uint8{1, 2})))
	require.True(t, a.Equals(vf.Value(reflect.ValueOf([]uint8{1, 2}))))
}

func TestBinary_Freeze(t *testing.T) {
	a := vf.Binary([]byte{1, 2}, false)
	require.False(t, a.Frozen())
	a.Freeze()
	require.True(t, a.Frozen())
}

func TestBinary_FrozenCopy(t *testing.T) {
	a := vf.Binary([]byte{1, 2}, false)
	b := a.FrozenCopy().(dgo.Binary)
	require.False(t, a.Frozen())
	require.True(t, b.Frozen())
	a.Freeze()

	b = a.FrozenCopy().(dgo.Binary)
	require.Same(t, a, b)
	require.True(t, a.Frozen())
}

func TestBinary_FrozenEqual(t *testing.T) {
	f := vf.Binary([]byte{1, 2}, true)
	require.True(t, f.Frozen(), `not frozen`)

	a := f.Copy(false)
	require.False(t, a.Frozen(), `frozen`)

	require.Equal(t, f, a)
	require.Equal(t, a, f)

	a.Freeze()
	require.True(t, a.Frozen(), `not frozen`)
	require.Same(t, a, a.Copy(true))

	b := a.Copy(false)
	require.NotSame(t, a, b)
	require.NotSame(t, b, b.Copy(true))
	require.NotSame(t, b, b.Copy(false))
}

func TestBinary_ReflectTo(t *testing.T) {
	var s []byte
	b := vf.Binary([]byte{1, 2}, true)
	b.ReflectTo(reflect.ValueOf(&s).Elem())
	require.Equal(t, b, s)

	s = nil
	sp := &s
	b.ReflectTo(reflect.ValueOf(&sp).Elem())
	s = *sp
	require.Equal(t, b, s)

	var mi interface{}
	mip := &mi
	b.ReflectTo(reflect.ValueOf(mip).Elem())

	bc, ok := mi.([]byte)
	require.True(t, ok)
	require.Equal(t, b, bc)
}
