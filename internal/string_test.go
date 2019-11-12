package internal_test

import (
	"math"
	"reflect"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestPattern(t *testing.T) {
	tp := tf.Pattern(regexp.MustCompile(`^doh$`)).(dgo.StringType)
	require.Instance(t, tp, `doh`)
	require.Instance(t, tp, vf.Value(`doh`))
	require.NotInstance(t, tp, `dog`)
	require.NotInstance(t, tp, vf.Value(`dog`))
	require.NotInstance(t, tp, 3)
	require.Assignable(t, typ.String, tp)
	require.Assignable(t, tf.Pattern(regexp.MustCompile(`^doh$`)), tp)
	require.NotAssignable(t, tf.String(3, 3), tp)
	require.NotAssignable(t, tf.Enum(`doh`), tp)
	require.NotAssignable(t, tf.Pattern(regexp.MustCompile(`doh`)), tp)

	require.True(t, tp.Unbounded())
	require.Equal(t, 0, tp.Min())
	require.Equal(t, math.MaxInt64, tp.Max())

	require.Equal(t, tp, tf.Pattern(regexp.MustCompile(`^doh$`)))
	require.NotEqual(t, tp, typ.String)

	require.NotEqual(t, 0, tp.HashCode())
	require.Equal(t, tp.HashCode(), tf.Pattern(regexp.MustCompile(`^doh$`)).HashCode())

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, `/^doh$/`, tp.String())

	require.Same(t, typ.String, typ.Generic(tp))

	s := "a\tb"
	require.Equal(t, `/a\tb/`, tf.Pattern(regexp.MustCompile(s)).String())

	s = "a\nb"
	require.Equal(t, `/a\nb/`, tf.Pattern(regexp.MustCompile(s)).String())

	s = "a\rb"
	require.Equal(t, `/a\rb/`, tf.Pattern(regexp.MustCompile(s)).String())

	s = "a\\b"
	require.Equal(t, `/a\\b/`, tf.Pattern(regexp.MustCompile(s)).String())

	s = "a\u0014b"
	require.Equal(t, `/a\u{14}b/`, tf.Pattern(regexp.MustCompile(s)).String())

	require.Equal(t, `/a\/b/`, tf.Pattern(regexp.MustCompile(`a/b`)).String())

	require.Equal(t, typ.String.ReflectType(), tp.ReflectType())
}

func TestStringDefault(t *testing.T) {
	tp := typ.String
	require.Instance(t, tp, `doh`)
	require.NotInstance(t, tp, 1)
	require.Assignable(t, tp, tp)
	require.Assignable(t, tp, typ.DgoString)
	require.Instance(t, tp.Type(), tp)
	require.Assignable(t, tp, tf.Pattern(regexp.MustCompile(`^doh$`)))
	require.NotAssignable(t, tf.String(3, 3), tp)
	require.NotAssignable(t, tf.Enum(`doh`), tp)
	require.NotAssignable(t, tf.Pattern(regexp.MustCompile(`doh`)), tp)

	require.Equal(t, 0, tp.Min())
	require.Equal(t, math.MaxInt64, tp.Max())
	require.True(t, tp.Unbounded())

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `string`, tp.String())

	require.True(t, reflect.ValueOf(`hello`).Type().AssignableTo(typ.String.ReflectType()))
}

func TestStringExact(t *testing.T) {
	tp := vf.Value(`doh`).Type().(dgo.StringType)
	require.Instance(t, tp, `doh`)
	require.NotInstance(t, tp, `duh`)
	require.NotInstance(t, tp, 3)
	require.Instance(t, tp.Type(), tp)
	require.Assignable(t, typ.String, tp)
	require.Assignable(t, tp, tp)
	require.Assignable(t, tf.String(3, 3), tp)
	require.Assignable(t, tf.Enum(`doh`, `duh`), tp)
	require.Assignable(t, tf.Pattern(regexp.MustCompile(`^doh$`)), tp)
	require.NotAssignable(t, tf.Pattern(regexp.MustCompile(`^duh$`)), tp)
	require.NotAssignable(t, tp, tf.Enum(`doh`, `duh`))
	require.NotAssignable(t, tp, tf.Pattern(regexp.MustCompile(`^doh$`)))
	require.NotAssignable(t, tp, typ.String)
	require.NotAssignable(t, tp, typ.Integer)
	require.Equal(t, tp, vf.Value(`doh`).Type())
	require.NotEqual(t, tp, vf.Value(`duh`).Type())
	require.NotEqual(t, tp, vf.Value(3).Type())

	require.Equal(t, 3, tp.Min())
	require.Equal(t, 3, tp.Max())
	require.False(t, tp.Unbounded())

	require.Same(t, typ.String, typ.Generic(tp))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `"doh"`, tp.String())
	require.Equal(t, typ.String.ReflectType(), tp.ReflectType())
}

func TestCiString(t *testing.T) {
	tp := tf.CiString(`abc`)
	require.Equal(t, tp, tp)
	require.Equal(t, tp, tf.CiString(`ABC`))
	require.NotEqual(t, tp, vf.String(`abc`).Type())

	require.Instance(t, tp, `abc`)
	require.Instance(t, tp, `ABC`)
	require.Instance(t, tp, vf.String(`aBc`))
	require.NotInstance(t, tp, `cde`)
	require.NotInstance(t, tp, []byte(`abc`))

	require.Instance(t, tp.Type(), tp)
}

func TestString_badOneArg(t *testing.T) {
	require.Panic(t, func() { tf.String(true) }, `illegal argument`)
}

func TestString_badTwoArg(t *testing.T) {
	require.Panic(t, func() { tf.String(`bad`, 2) }, `illegal argument 1`)
	require.Panic(t, func() { tf.String(1, `bad`) }, `illegal argument 2`)
}

func TestString_badArgCount(t *testing.T) {
	require.Panic(t, func() { tf.String(2, 2, true) }, `illegal number of arguments`)
}

func TestStringType(t *testing.T) {
	tp := tf.String()
	require.Same(t, tp, typ.String)

	tp = tf.String(0, math.MaxInt64)
	require.Same(t, tp, typ.String)

	tp = tf.String(`hello`)
	require.Equal(t, tp, vf.String(`hello`).Type())

	tp = tf.String(1)
	require.Equal(t, tp, tf.String(1, math.MaxInt64))
	require.Assignable(t, tp, typ.DgoString)

	tp = tf.String(2)
	require.NotAssignable(t, tp, typ.DgoString)

	tp = tf.String(3, 5)
	require.Instance(t, tp, `doh`)
	require.NotInstance(t, tp, `do`)
	require.Instance(t, tp, `dudoh`)
	require.Instance(t, tp, vf.Value(`dudoh`))
	require.NotInstance(t, tp, `duhdoh`)
	require.NotInstance(t, tp, 3)
	require.Instance(t, tp.Type(), tp)
	require.Assignable(t, typ.String, tp)
	require.Assignable(t, tp, tp)
	require.Assignable(t, tp, tf.String(3, 3))
	require.NotAssignable(t, tp, tf.String(2, 3))
	require.NotAssignable(t, tp, tf.String(3, 6))
	require.Assignable(t, tp, tf.Enum(`doh`, `duh`))
	require.NotAssignable(t, tp, tf.Enum(`doh`, `duhduh`))
	require.NotAssignable(t, tf.Enum(`doh`, `duh`), tp)
	require.NotAssignable(t, tf.Pattern(regexp.MustCompile(`^doh$`)), tp)
	require.NotAssignable(t, tp, tf.Pattern(regexp.MustCompile(`^doh$`)))
	require.NotAssignable(t, tp, typ.String)
	require.NotAssignable(t, tp, typ.Integer)
	require.Equal(t, tp, tf.String(3, 5))
	require.Equal(t, tp, tf.String(5, 3))
	require.Equal(t, tf.String(-3, 3), tf.String(0, 3))
	require.NotEqual(t, tp, tf.String(3, 4))
	require.NotEqual(t, tp, tf.String(2, 5))
	require.NotEqual(t, tp, vf.Value(3).Type())

	require.Equal(t, 3, tp.Min())
	require.Equal(t, 5, tp.Max())
	require.False(t, tp.Unbounded())

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `string[3,5]`, tp.String())
	require.Equal(t, typ.String.ReflectType(), tp.ReflectType())
}

func TestStringType_New(t *testing.T) {
	require.Equal(t, `0xc`, vf.New(typ.String, vf.Arguments(12, `%#x`)))
	require.Equal(t, `23`, vf.New(typ.String, vf.Arguments(23)))

	require.Equal(t, `a`, vf.New(tf.CiString(`a`), vf.String(`a`)))
	require.Equal(t, `A`, vf.New(tf.CiString(`a`), vf.String(`A`)))
	require.Equal(t, `string`, vf.New(typ.DgoString, vf.String(`string`)))
	require.Panic(t, func() { vf.New(tf.CiString(`a`), vf.String(`b`)) }, `cannot be assigned`)
	require.Panic(t, func() { vf.New(vf.String(`a`).Type(), vf.String(`b`)) }, `cannot be assigned`)
	require.Panic(t, func() { vf.New(tf.Pattern(regexp.MustCompile(`a`)), vf.String(`d`)) }, `cannot be assigned`)
	require.Panic(t, func() { vf.New(tf.String(2, 3), vf.String(`d`)) }, `cannot be assigned`)
}

func TestDgoStringType(t *testing.T) {
	require.Assignable(t, typ.DgoString, typ.DgoString)
	require.Equal(t, typ.DgoString, typ.DgoString)
	require.NotEqual(t, typ.DgoString, typ.String)
	require.NotAssignable(t, typ.DgoString, tf.OneOf(typ.DgoString, typ.String))
	require.Instance(t, typ.DgoString.Type(), typ.DgoString)

	s := `dgo`
	require.Assignable(t, typ.DgoString, vf.String(s).Type())
	require.Instance(t, typ.DgoString, s)
	require.Instance(t, typ.DgoString, vf.String(s))

	s = `map[string]1..3`
	require.Instance(t, typ.DgoString, s)
	require.Instance(t, typ.DgoString, vf.String(s))
	require.Assignable(t, typ.DgoString, vf.String(s).Type())
	require.False(t, typ.DgoString.Unbounded())
	require.Equal(t, 1, typ.DgoString.Min())
	require.Equal(t, math.MaxInt64, typ.DgoString.Max())
	require.Equal(t, `dgo`, typ.DgoString.String())

	s = `hello`
	require.NotInstance(t, typ.DgoString, s)
	require.NotInstance(t, typ.DgoString, vf.String(s))
	require.NotAssignable(t, typ.DgoString, vf.String(s).Type())

	require.NotEqual(t, 0, typ.DgoString.HashCode())
	require.Equal(t, typ.String.ReflectType(), typ.DgoString.ReflectType())
}

func TestString(t *testing.T) {
	v := vf.String(`hello`)
	require.Equal(t, v, `hello`)
	require.Equal(t, v, vf.String(`hello`))
	require.NotEqual(t, v, `hi`)
	require.NotEqual(t, v, vf.String(`hi`))
	require.NotEqual(t, v, 3)

	c, ok := vf.String(`hello`).CompareTo(`hello`)
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.String(`hello`).CompareTo(v)
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.String(`hallo`).CompareTo(`hello`)
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.String(`hallo`).CompareTo(v)
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.String(`hi`).CompareTo(`hello`)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.String(`hi`).CompareTo(v)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = v.CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	_, ok = v.CompareTo(3)
	require.False(t, ok)

	require.True(t, `hello` == v.GoString())
}
