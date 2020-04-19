package internal_test

import (
	"regexp"
	"testing"

	"github.com/tada/dgo/dgo"
	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestAllOfType(t *testing.T) {
	tp := tf.AllOf(tf.Enum(`a`, `b`, `c`), tf.Enum(`b`, `c`, `d`))
	require.Instance(t, tp, `b`)
	require.Instance(t, tp, `c`)
	require.NotInstance(t, tp, `a`)
	require.Assignable(t, typ.AllOf, tp)
	require.NotAssignable(t, tp, typ.AllOf)
	require.Assignable(t, tp, tf.AllOf(tf.Enum(`c`, `b`, `a`), tf.Enum(`b`, `d`, `c`)))
	require.Assignable(t, tp, tf.AllOf(tf.Enum(`a`, `b`, `c`), tf.Enum(`b`, `c`, `d`), tf.Enum(`c`, `d`, `e`)))
	require.NotAssignable(t, tp, tf.AllOf(tf.Enum(`a`, `b`, `c`), tf.Enum(`e`, `c`, `d`)))

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, tp, tf.AllOf(tf.Enum(`b`, `c`, `d`), tf.Enum(`a`, `b`, `c`)))
	require.NotEqual(t, tp, tf.AllOf(tf.Enum(`b`, `c`), tf.Enum(`a`, `b`, `c`)))
	require.NotEqual(t, tp, tf.Enum(`a`, `b`, `c`))

	require.NotEqual(t, tf.AllOf(`b`, `c`), tf.AnyOf(`b`, `c`))
	require.NotEqual(t, tf.AllOf(`b`, `c`), tf.OneOf(`b`, `c`))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Same(t, tf.AllOf(), typ.Any)
	require.Same(t, tf.AllOf(typ.String), typ.String)

	require.Equal(t, tp.(dgo.TernaryType).Operands(), vf.Values(tf.Enum(`a`, `b`, `c`), tf.Enum(`b`, `c`, `d`)))
	require.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpAnd)
	require.Equal(t, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType().(dgo.TernaryType).Operator(), dgo.OpAnd)

	require.Equal(t, `("a"|"b"|"c")&("b"|"c"|"d")`, tp.String())

	require.Equal(t, typ.String.ReflectType(), tp.ReflectType())

	tp = tf.AllOf(
		tf.Pattern(regexp.MustCompile(`a`)),
		tf.Pattern(regexp.MustCompile(`b`)),
		tf.Pattern(regexp.MustCompile(`c`)))
	require.Instance(t, tp, `abc`)
	require.NotInstance(t, tp, `c`)

	require.Assignable(t, tp, tf.AllOf(
		tf.Pattern(regexp.MustCompile(`a`)),
		tf.Pattern(regexp.MustCompile(`b`)), tf.Pattern(regexp.MustCompile(`c`))))
	require.NotAssignable(t, tp, tf.AllOf(tf.Pattern(regexp.MustCompile(`a`)), tf.Pattern(regexp.MustCompile(`b`))))
	require.Assignable(t, tp, tf.AllOf(tf.Pattern(regexp.MustCompile(`a`)), vf.Value(`abc`).Type()))
	require.NotAssignable(t, tp, tf.AllOf(tf.Pattern(regexp.MustCompile(`a`)), vf.Value(`abc`).Type(), typ.String))

	require.Assignable(t, tp, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType())
	require.NotAssignable(t, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType(), tp)
	require.Assignable(t, tp, tf.ParseType(`"abc"^/a/&/b/&/c/`))
	require.Assignable(t, tf.ParseType(`/a/|/b/|/c/`), tp)

	a := vf.Values(`a`, `b`, nil)
	require.Same(t, a.Type().(dgo.ArrayType).ElementType().(dgo.ExactType).ExactValue(), a)
	require.Same(t, typ.Any, typ.Generic(a.Type()).(dgo.ArrayType).ElementType())
}

func TestAllOfValueType(t *testing.T) {
	tp := vf.Values(`a`, `b`).Type().(dgo.ArrayType).ElementType()
	require.Equal(t, tf.AllOf(`a`, `b`), tp)
	require.Equal(t, tp, tf.AllOf(`a`, `b`))
	require.NotEqual(t, tp, tf.AllOf(`a`, `b`, `c`))
	require.NotEqual(t, tp, vf.Values(`a`, `b`))
}

func TestAnyOfType(t *testing.T) {
	tp := tf.AnyOf(typ.Integer, typ.String)
	require.Instance(t, tp, `hello`)
	require.Instance(t, tp, 3)
	require.NotInstance(t, tp, true)
	require.Assignable(t, typ.AnyOf, tp)
	require.NotAssignable(t, tp, typ.AnyOf)
	require.Assignable(t, tp, tf.AnyOf(typ.String, typ.Integer))
	require.Assignable(t, tf.AnyOf(typ.String, typ.Integer, typ.Boolean), tp)
	require.NotAssignable(t, tp, tf.AnyOf(typ.String, typ.Integer, typ.Boolean))
	require.NotAssignable(t, tf.AnyOf(typ.String, typ.Boolean), tp)

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, tp, tf.AnyOf(typ.Integer, typ.String))
	require.NotEqual(t, tp, tf.AnyOf(typ.Integer, typ.Boolean))
	require.NotEqual(t, tp, typ.Integer)

	require.NotEqual(t, tf.AnyOf(`b`, `c`), tf.AllOf(`b`, `c`))
	require.NotEqual(t, tf.AnyOf(`b`, `c`), tf.OneOf(`b`, `c`))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, tf.AnyOf(), tf.Not(typ.Any))
	require.Same(t, tf.AnyOf(typ.String), typ.String)

	require.Equal(t, tp.(dgo.TernaryType).Operands(), vf.Values(typ.Integer, typ.String))
	require.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpOr)

	require.Equal(t, `int|string`, tp.String())

	require.Equal(t, typ.Any.ReflectType(), tp.ReflectType())
}

func TestOneOfType(t *testing.T) {
	tp := tf.OneOf(typ.Integer, tf.Pattern(regexp.MustCompile(`a`)), tf.Pattern(regexp.MustCompile(`b`)))
	require.Instance(t, tp, `a`)
	require.Instance(t, tp, 3)
	require.NotInstance(t, tp, `ab`)
	require.NotInstance(t, tp, true)
	require.Assignable(t, typ.OneOf, tp)
	require.NotAssignable(t, tp, typ.OneOf)
	require.NotAssignable(t, tp, tf.OneOf(typ.String, typ.Integer))
	require.Assignable(t, tp, tf.OneOf(vf.String(`a`).Type(), vf.String(`b`).Type(), typ.Integer))
	require.NotAssignable(t, tp, tf.OneOf(vf.String(`ab`).Type(), typ.Integer))
	require.Assignable(t, tf.OneOf(typ.String, typ.Integer), tp)
	require.Assignable(t, tf.AnyOf(typ.String, typ.Integer, typ.Boolean), tp)
	require.NotAssignable(t, tp, tf.AnyOf(typ.String, typ.Integer, typ.Boolean))
	require.NotAssignable(t, tf.AnyOf(typ.String, typ.Boolean), tp)

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, tp, tf.OneOf(typ.Integer, tf.Pattern(regexp.MustCompile(`a`)), tf.Pattern(regexp.MustCompile(`b`))))
	require.NotEqual(t, tp, tf.OneOf(typ.Integer, typ.Boolean))
	require.NotEqual(t, tp, typ.Integer)

	require.NotEqual(t, tf.OneOf(`b`, `c`), tf.AllOf(`b`, `c`))
	require.NotEqual(t, tf.OneOf(`b`, `c`), tf.AnyOf(`b`, `c`))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, tf.OneOf(), tf.Not(typ.Any))
	require.Same(t, tf.OneOf(typ.String), typ.String)

	require.Equal(t,
		tp.(dgo.TernaryType).Operands(),
		vf.Values(typ.Integer, tf.Pattern(regexp.MustCompile(`a`)), tf.Pattern(regexp.MustCompile(`b`))))
	require.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpOne)

	require.Equal(t, `int^/a/^/b/`, tp.String())
	require.Equal(t, typ.Any.ReflectType(), tp.ReflectType())
}

func TestEnum(t *testing.T) {
	tp := tf.Enum()
	require.Equal(t, tp, tf.Not(typ.Any))

	tp = tf.Enum(`f`, `foo`, `foobar`)
	require.Instance(t, tp, `foo`)
	require.NotInstance(t, tp, `oops`)
	require.Assignable(t, typ.String, tp)
	require.NotAssignable(t, tp, tf.CiEnum(`f`, `foo`, `foobar`))
	require.Assignable(t, tf.String(1, 6), tp)
	require.NotAssignable(t, tf.String(3, 3), tp)
	require.Assignable(t, tf.Pattern(regexp.MustCompile(`f`)), tp)
	require.NotAssignable(t, tf.Pattern(regexp.MustCompile(`o`)), tp)
	require.Equal(t, tp, tf.Enum(`foo`, `f`, `foobar`))
	require.NotEqual(t, tp, tf.Enum(`foo`, `f`, `foobar`, `x`))
	require.NotEqual(t, tp, tf.Enum(`foo`, `foobar`))
}

func TestCiEnum(t *testing.T) {
	tp := tf.CiEnum()
	require.Equal(t, tp, tf.Not(typ.Any))

	tp = tf.CiEnum(`f`, `foo`, `foobar`)
	require.Instance(t, tp, `FOO`)
	require.NotInstance(t, tp, `Oops`)
	require.Assignable(t, typ.String, tp)
	require.NotAssignable(t, tf.String(3, 3), tp)
	require.Assignable(t, tf.String(1, 6), tp)
	require.Assignable(t, tp, tf.CiEnum(`foo`))
	require.NotAssignable(t, tf.Pattern(regexp.MustCompile(`f`)), tp)
	require.Assignable(t, tp, tf.Enum(`f`, `foo`, `foobar`))

	require.Equal(t, `~"f"|~"foo"|~"foobar"`, tp.String())
}

func TestIntEnum(t *testing.T) {
	tp := tf.IntEnum()
	require.Equal(t, tp, tf.Not(typ.Any))

	tp = tf.IntEnum(2, 4, 8)
	require.Instance(t, tp, 4)
	require.NotInstance(t, tp, 5)
	require.Assignable(t, typ.Integer, tp)
	require.NotAssignable(t, tf.Integer(3, 3, true), tp)
	require.Assignable(t, tf.Integer(2, 8, true), tp)
	require.Assignable(t, tp, tf.IntEnum(8))
	require.Assignable(t, tp, tf.IntEnum(2, 4))
	require.Assignable(t, tp, tf.IntEnum(8, 2, 4))

	require.Equal(t, `2|4|8`, tp.String())
}
