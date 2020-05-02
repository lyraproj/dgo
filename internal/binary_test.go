package internal_test

import (
	"bytes"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/test/require"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestBinaryType(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.Binary(bs, false)
	tp := tf.Binary(3, 3)
	assert.Assignable(t, typ.Binary, tp)
	assert.NotAssignable(t, tp, typ.Binary)
	assert.NotAssignable(t, tp, typ.String)
	assert.Instance(t, typ.Binary, bs)
	assert.Instance(t, typ.Binary, v)
	assert.Instance(t, tp, v)
	assert.Instance(t, tp, []byte{1, 2, 3})
	assert.NotInstance(t, tp, []byte{1, 2})
	assert.NotInstance(t, tf.Binary(1, 5), `abc`)
	assert.Assignable(t, tf.Binary(3, 3), tp)
	assert.Assignable(t, tp, tf.Binary(3, 3))
	assert.Same(t, typ.Binary, tf.Binary())
	assert.Same(t, typ.Binary, tf.Binary(0, math.MaxInt64))
	assert.NotAssignable(t, tf.Binary(4), tp)
	assert.NotAssignable(t, tf.Binary(0, 2), tp)
	assert.Equal(t, tf.Binary(1, 2), tf.Binary(2, 1))
	assert.Equal(t, tf.Binary(-1, 2), tf.Binary(0, 2))
	assert.NotEqual(t, tf.Binary(0, 2), tf.String(0, 2))
	assert.Equal(t, tf.Binary(0, 2).HashCode(), tf.Binary(0, 2).HashCode())
	assert.NotEqual(t, tf.Binary(1, 2).HashCode(), tf.Binary(0, 2).HashCode())
	assert.NotEqual(t, tf.Binary(0, 1).HashCode(), tf.Binary(0, 2).HashCode())
	assert.Instance(t, tp.Type(), tp)
	assert.Panic(t, func() { tf.Binary(`blue`) }, `illegal argument`)
	assert.Panic(t, func() { tf.Binary(`blue`, 1) }, `illegal argument 1`)
	assert.Panic(t, func() { tf.Binary(1, `blue`) }, `illegal argument 2`)
	assert.Panic(t, func() { tf.Binary(1, 2, 3) }, `illegal number of arguments`)
	assert.Equal(t, `binary[3,3]`, tp.String())
	assert.Equal(t, reflect.TypeOf([]byte{}), tp.ReflectType())
}

func TestBinaryType_New(t *testing.T) {
	b := vf.New(typ.Binary, vf.Arguments(vf.Values(1, 2, 3)))
	assert.Equal(t, vf.Binary([]byte{1, 2, 3}, true), b)
}

func TestBinaryType_New_passThrough(t *testing.T) {
	b := vf.Binary([]byte{1, 2, 3}, true)
	assert.Same(t, b, vf.New(typ.Binary, b))
}

func TestBinaryType_New_stringRaw(t *testing.T) {
	assert.Equal(t, vf.Binary([]byte(`hello`), true), vf.New(typ.Binary, vf.Arguments(vf.String(`hello`), `%r`)))
}

func TestBinaryType_New_badCount(t *testing.T) {
	assert.Panic(t, func() { vf.New(typ.Binary, vf.Arguments(1, 2, 3)) }, `illegal number of arguments`)
}

func TestBinaryType_New_badArg(t *testing.T) {
	assert.Panic(t, func() { vf.New(typ.Binary, vf.Integer(3)) }, `illegal argument for binary`)
}

func TestBinaryType_New_badType(t *testing.T) {
	assert.Panic(t,
		func() { vf.New(tf.Binary(2, 2), vf.New(typ.Binary, vf.Value([]byte{1, 2, 3}))) },
		`cannot be assigned`)
}

func TestBinaryType_New_badBytes(t *testing.T) {
	assert.Panic(t,
		func() { vf.New(typ.Binary, vf.Values(1, 311, 3)) },
		`the value 311 cannot be assigned to a variable of type 0..255`)
}

func TestBinaryType_New_badFormat(t *testing.T) {
	assert.Panic(t, func() { vf.New(typ.Binary, vf.Arguments(`hello`, `%x`)) }, `illegal argument`)
}

func TestExactBinaryType(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.Binary(bs, false)
	tp := v.Type().(dgo.BinaryType)
	assert.Assignable(t, typ.Binary, tp)
	assert.NotAssignable(t, tp, typ.Binary)
	assert.NotAssignable(t, tp, typ.String)
	assert.Instance(t, tp, v)
	assert.Instance(t, tp, bs)
	assert.True(t, tp.IsInstance(bs))
	assert.NotInstance(t, tp, []byte{1, 2})
	assert.NotInstance(t, tp, "AQID")
	assert.False(t, tp.Unbounded())
	assert.Assignable(t, tf.Binary(3, 3), tp)
	assert.NotAssignable(t, tp, tf.Binary(3, 3))
	assert.NotAssignable(t, tf.Binary(4), tp)
	assert.NotAssignable(t, tf.Binary(0, 2), tp)
	assert.Equal(t, v.HashCode(), tp.HashCode())
	assert.Instance(t, tp.Type(), tp)
	assert.Equal(t, `binary "AQID"`, tp.String())
	assert.Equal(t, reflect.TypeOf([]byte{}), tp.ReflectType())

	b := vf.New(tp, vf.Arguments(vf.Values(1, 2, 3)))
	assert.Equal(t, vf.Binary([]byte{1, 2, 3}, true), b)
	assert.Panic(t,
		func() { vf.New(tp, vf.New(typ.Binary, vf.Value([]byte{1, 2}))) },
		`cannot be assigned`)
}

type badReader int

func (badReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New(`oops`)
}

func TestBinary(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.Binary(bs, false)
	assert.Equal(t, `AQID`, v.Encode())

	// Test mutability
	bs[1] = 4
	assert.Equal(t, reflect.ValueOf(bs).Pointer(), reflect.ValueOf(v.GoBytes()).Pointer())
	assert.Equal(t, `AQQD`, v.Encode())

	// Test immutability
	v = vf.Binary(bs, true)
	bs[1] = 2
	assert.Equal(t, `AQQD`, v.Encode())
	assert.NotEqual(t, reflect.ValueOf(bs).Pointer(), reflect.ValueOf(v.GoBytes()).Pointer())

	// Test panic on bad read
	assert.Panic(t, func() { vf.BinaryFromData(badReader(0)) }, `oops`)
}

func TestBinaryString(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.BinaryFromString(`AQID`)
	assert.True(t, bytes.Equal(bs, v.GoBytes()))
	assert.Panic(t, func() { vf.BinaryFromString(`----`) }, `illegal base64`)
}

func TestBinaryFromEncoded(t *testing.T) {
	bs := []byte(`hello`)
	notStrict := `aGVsbG9=`
	v := vf.BinaryFromEncoded(notStrict, `%b`)
	assert.True(t, bytes.Equal(bs, v.GoBytes()))
	assert.Panic(t, func() { vf.BinaryFromEncoded(notStrict, `%B`) }, `illegal base64 data at input byte 7`)

	bs = []byte{'r', 0x82, 0xff}
	notUtf8 := string(bs)
	v = vf.BinaryFromEncoded(notUtf8, `%r`)
	assert.True(t, bytes.Equal(bs, v.GoBytes()))
	assert.Panic(t, func() { vf.BinaryFromEncoded(notUtf8, `%s`) }, `Expected valid utf8 string`)

	bs = []byte{3, 240, 126}
	assert.True(t, bytes.Equal(bs, vf.BinaryFromEncoded(`A/B+`, `%B`).GoBytes()))
	assert.True(t, bytes.Equal(bs, vf.BinaryFromEncoded(`A_B-`, `%u`).GoBytes()))
	assert.Panic(t, func() { vf.BinaryFromEncoded(`A/B+`, `%u`) }, `illegal base64 data at input byte 1`)
	assert.Panic(t, func() { vf.BinaryFromEncoded(`A/B+`, `%x`) }, `Expected one of the supported format specifiers`)
}

func TestBinary_CompareTo(t *testing.T) {
	a := vf.Binary([]byte{'a', 'b', 'c'}, false)

	c, ok := a.CompareTo(a)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = a.CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	_, ok = a.CompareTo(vf.Value('a'))
	assert.False(t, ok)

	b := vf.Binary([]byte{'a', 'b', 'c'}, false)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = a.CompareTo(b.FrozenCopy())
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	bs := []byte{'a', 'b', 'c'}
	c, ok = a.CompareTo(bs)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	b = vf.Binary([]byte{'a', 'b', 'c', 'd'}, false)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	b = vf.Binary([]byte{'a', 'b', 'd', 'd'}, false)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	b = vf.Binary([]byte{'a', 'b'}, false)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 1, c)

	b = vf.Binary([]byte{'a', 'b', 'd'}, false)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, -1, c)

	b = vf.Binary([]byte{'a', 'b', 'b'}, false)
	c, ok = a.CompareTo(b)
	assert.True(t, ok)
	assert.Equal(t, 1, c)
}

func TestBinary_Copy(t *testing.T) {
	a := vf.Binary([]byte{'a', 'b'}, true)
	assert.Same(t, a, a.Copy(true))

	c := a.Copy(false)
	assert.False(t, c.Frozen())
	assert.NotSame(t, c, c.Copy(false))

	c = c.Copy(true)
	assert.True(t, c.Frozen())
	assert.Same(t, c, c.Copy(true))

	c = a.ThawedCopy().(dgo.Binary)
	assert.False(t, c.Frozen())
	assert.NotSame(t, c, c.ThawedCopy())
}

func TestBinary_Equal(t *testing.T) {
	a := vf.Binary([]byte{1, 2}, true)
	assert.True(t, a.Equals(a))

	b := vf.Binary([]byte{1}, true)
	assert.False(t, a.Equals(b))
	b = vf.Binary([]byte{1, 2}, true)
	assert.True(t, a.Equals(b))
	assert.True(t, a.Equals([]byte{1, 2}))
	assert.False(t, a.Equals(`12`))
	assert.Equal(t, a.HashCode(), b.HashCode())
	assert.True(t, a.Equals(vf.Value([]uint8{1, 2})))
	assert.True(t, a.Equals(vf.Value(reflect.ValueOf([]uint8{1, 2}))))
}

func TestBinary_FrozenCopy(t *testing.T) {
	a := vf.Binary([]byte{1, 2}, false)
	b := a.FrozenCopy().(dgo.Binary)
	assert.False(t, a.Frozen())
	assert.True(t, b.Frozen())

	a = a.FrozenCopy().(dgo.Binary)
	b = a.FrozenCopy().(dgo.Binary)
	assert.Same(t, a, b)
	assert.True(t, a.Frozen())
}

func TestBinary_FrozenEqual(t *testing.T) {
	f := vf.Binary([]byte{1, 2}, true)
	assert.True(t, f.Frozen(), `not frozen`)

	a := f.Copy(false)
	assert.False(t, a.Frozen(), `frozen`)
	assert.Equal(t, f, a)
	assert.Equal(t, a, f)

	a = a.FrozenCopy().(dgo.Binary)
	assert.True(t, a.Frozen(), `not frozen`)
	assert.Same(t, a, a.Copy(true))

	b := a.Copy(false)
	assert.NotSame(t, a, b)
	assert.NotSame(t, b, b.Copy(true))
	assert.NotSame(t, b, b.Copy(false))
}

func TestBinary_ReflectTo(t *testing.T) {
	var s []byte
	ba := []byte{1, 2}
	b := vf.Binary(ba, true)
	b.ReflectTo(reflect.ValueOf(&s).Elem())
	assert.True(t, bytes.Equal(ba, s))
	s[0] = 2
	assert.Equal(t, ba[0], 1)

	bm := vf.Binary(ba, false)
	bm.ReflectTo(reflect.ValueOf(&s).Elem())
	assert.True(t, bytes.Equal(ba, s))
	s[0] = 2
	assert.Equal(t, ba[0], 2)

	s = nil
	sp := &s
	b.ReflectTo(reflect.ValueOf(&sp).Elem())
	s = *sp
	assert.Equal(t, b, s)

	var mi interface{}
	mip := &mi
	b.ReflectTo(reflect.ValueOf(mip).Elem())

	bc, ok := mi.([]byte)
	require.True(t, ok)
	assert.Equal(t, b, bc)
}
