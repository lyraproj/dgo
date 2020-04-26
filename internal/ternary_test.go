package internal_test

import (
	"regexp"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestAllOfType(t *testing.T) {
	tp := tf.AllOf(tf.Enum(`a`, `b`, `c`), tf.Enum(`b`, `c`, `d`))
	assert.Instance(t, tp, `b`)
	assert.Instance(t, tp, `c`)
	assert.NotInstance(t, tp, `a`)
	assert.Assignable(t, typ.AllOf, tp)
	assert.NotAssignable(t, tp, typ.AllOf)
	assert.Assignable(t, tp, tf.AllOf(tf.Enum(`c`, `b`, `a`), tf.Enum(`b`, `d`, `c`)))
	assert.Assignable(t, tp, tf.AllOf(tf.Enum(`a`, `b`, `c`), tf.Enum(`b`, `c`, `d`), tf.Enum(`c`, `d`, `e`)))
	assert.NotAssignable(t, tp, tf.AllOf(tf.Enum(`a`, `b`, `c`), tf.Enum(`e`, `c`, `d`)))
	assert.Instance(t, tp.Type(), tp)
	assert.Equal(t, tp, tf.AllOf(tf.Enum(`b`, `c`, `d`), tf.Enum(`a`, `b`, `c`)))
	assert.NotEqual(t, tp, tf.AllOf(tf.Enum(`b`, `c`), tf.Enum(`a`, `b`, `c`)))
	assert.NotEqual(t, tp, tf.Enum(`a`, `b`, `c`))
	assert.NotEqual(t, tf.AllOf(`b`, `c`), tf.AnyOf(`b`, `c`))
	assert.NotEqual(t, tf.AllOf(`b`, `c`), tf.OneOf(`b`, `c`))
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Same(t, tf.AllOf(), typ.Any)
	assert.Same(t, tf.AllOf(typ.String), typ.String)
	assert.Equal(t, tp.(dgo.TernaryType).Operands(), vf.Values(tf.Enum(`a`, `b`, `c`), tf.Enum(`b`, `c`, `d`)))
	assert.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpAnd)
	assert.Equal(t, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType().(dgo.TernaryType).Operator(), dgo.OpAnd)
	assert.Equal(t, `("a"|"b"|"c")&("b"|"c"|"d")`, tp.String())
	assert.Equal(t, typ.String.ReflectType(), tp.ReflectType())

	tp = tf.AllOf(
		tf.Pattern(regexp.MustCompile(`a`)),
		tf.Pattern(regexp.MustCompile(`b`)),
		tf.Pattern(regexp.MustCompile(`c`)))
	assert.Instance(t, tp, `abc`)
	assert.NotInstance(t, tp, `c`)
	assert.Assignable(t, tp, tf.AllOf(
		tf.Pattern(regexp.MustCompile(`a`)),
		tf.Pattern(regexp.MustCompile(`b`)), tf.Pattern(regexp.MustCompile(`c`))))
	assert.NotAssignable(t, tp, tf.AllOf(tf.Pattern(regexp.MustCompile(`a`)), tf.Pattern(regexp.MustCompile(`b`))))
	assert.Assignable(t, tp, tf.AllOf(tf.Pattern(regexp.MustCompile(`a`)), vf.Value(`abc`).Type()))
	assert.NotAssignable(t, tp, tf.AllOf(tf.Pattern(regexp.MustCompile(`a`)), vf.Value(`abc`).Type(), typ.String))
	assert.Assignable(t, tp, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType())
	assert.NotAssignable(t, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType(), tp)
	assert.Assignable(t, tp, tf.ParseType(`"abc"^/a/&/b/&/c/`))
	assert.Assignable(t, tf.ParseType(`/a/|/b/|/c/`), tp)

	a := vf.Values(`a`, `b`, nil)
	assert.Same(t, a.Type().(dgo.ArrayType).ElementType().(dgo.ExactType).ExactValue(), a)
	assert.Same(t, typ.Any, typ.Generic(a.Type()).(dgo.ArrayType).ElementType())
}

func TestAllOfValueType(t *testing.T) {
	tp := vf.Values(`a`, `b`).Type().(dgo.ArrayType).ElementType()
	assert.Equal(t, tf.AllOf(`a`, `b`), tp)
	assert.Equal(t, tp, tf.AllOf(`a`, `b`))
	assert.NotEqual(t, tp, tf.AllOf(`a`, `b`, `c`))
	assert.NotEqual(t, tp, vf.Values(`a`, `b`))
}

func TestAnyOfType(t *testing.T) {
	tp := tf.AnyOf(typ.Integer, typ.String)
	assert.Instance(t, tp, `hello`)
	assert.Instance(t, tp, 3)
	assert.NotInstance(t, tp, true)
	assert.Assignable(t, typ.AnyOf, tp)
	assert.NotAssignable(t, tp, typ.AnyOf)
	assert.Assignable(t, tp, tf.AnyOf(typ.String, typ.Integer))
	assert.Assignable(t, tf.AnyOf(typ.String, typ.Integer, typ.Boolean), tp)
	assert.NotAssignable(t, tp, tf.AnyOf(typ.String, typ.Integer, typ.Boolean))
	assert.NotAssignable(t, tf.AnyOf(typ.String, typ.Boolean), tp)

	assert.Instance(t, tp.Type(), tp)

	assert.Equal(t, tp, tf.AnyOf(typ.Integer, typ.String))
	assert.NotEqual(t, tp, tf.AnyOf(typ.Integer, typ.Boolean))
	assert.NotEqual(t, tp, typ.Integer)

	assert.NotEqual(t, tf.AnyOf(`b`, `c`), tf.AllOf(`b`, `c`))
	assert.NotEqual(t, tf.AnyOf(`b`, `c`), tf.OneOf(`b`, `c`))

	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())

	assert.Equal(t, tf.AnyOf(), tf.Not(typ.Any))
	assert.Same(t, tf.AnyOf(typ.String), typ.String)

	assert.Equal(t, tp.(dgo.TernaryType).Operands(), vf.Values(typ.Integer, typ.String))
	assert.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpOr)

	assert.Equal(t, `int|string`, tp.String())

	assert.Equal(t, typ.Any.ReflectType(), tp.ReflectType())
}

func TestOneOfType(t *testing.T) {
	tp := tf.OneOf(typ.Integer, tf.Pattern(regexp.MustCompile(`a`)), tf.Pattern(regexp.MustCompile(`b`)))
	assert.Instance(t, tp, `a`)
	assert.Instance(t, tp, 3)
	assert.NotInstance(t, tp, `ab`)
	assert.NotInstance(t, tp, true)
	assert.Assignable(t, typ.OneOf, tp)
	assert.NotAssignable(t, tp, typ.OneOf)
	assert.NotAssignable(t, tp, tf.OneOf(typ.String, typ.Integer))
	assert.Assignable(t, tp, tf.OneOf(vf.String(`a`).Type(), vf.String(`b`).Type(), typ.Integer))
	assert.NotAssignable(t, tp, tf.OneOf(vf.String(`ab`).Type(), typ.Integer))
	assert.Assignable(t, tf.OneOf(typ.String, typ.Integer), tp)
	assert.Assignable(t, tf.AnyOf(typ.String, typ.Integer, typ.Boolean), tp)
	assert.NotAssignable(t, tp, tf.AnyOf(typ.String, typ.Integer, typ.Boolean))
	assert.NotAssignable(t, tf.AnyOf(typ.String, typ.Boolean), tp)

	assert.Instance(t, tp.Type(), tp)

	assert.Equal(t, tp, tf.OneOf(typ.Integer, tf.Pattern(regexp.MustCompile(`a`)), tf.Pattern(regexp.MustCompile(`b`))))
	assert.NotEqual(t, tp, tf.OneOf(typ.Integer, typ.Boolean))
	assert.NotEqual(t, tp, typ.Integer)

	assert.NotEqual(t, tf.OneOf(`b`, `c`), tf.AllOf(`b`, `c`))
	assert.NotEqual(t, tf.OneOf(`b`, `c`), tf.AnyOf(`b`, `c`))

	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())

	assert.Equal(t, tf.OneOf(), tf.Not(typ.Any))
	assert.Same(t, tf.OneOf(typ.String), typ.String)

	assert.Equal(t,
		tp.(dgo.TernaryType).Operands(),
		vf.Values(typ.Integer, tf.Pattern(regexp.MustCompile(`a`)), tf.Pattern(regexp.MustCompile(`b`))))
	assert.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpOne)

	assert.Equal(t, `int^/a/^/b/`, tp.String())
	assert.Equal(t, typ.Any.ReflectType(), tp.ReflectType())
}

func TestEnum(t *testing.T) {
	tp := tf.Enum()
	assert.Equal(t, tp, tf.Not(typ.Any))

	tp = tf.Enum(`f`, `foo`, `foobar`)
	assert.Instance(t, tp, `foo`)
	assert.NotInstance(t, tp, `oops`)
	assert.Assignable(t, typ.String, tp)
	assert.NotAssignable(t, tp, tf.CiEnum(`f`, `foo`, `foobar`))
	assert.Assignable(t, tf.String(1, 6), tp)
	assert.NotAssignable(t, tf.String(3, 3), tp)
	assert.Assignable(t, tf.Pattern(regexp.MustCompile(`f`)), tp)
	assert.NotAssignable(t, tf.Pattern(regexp.MustCompile(`o`)), tp)
	assert.Equal(t, tp, tf.Enum(`foo`, `f`, `foobar`))
	assert.NotEqual(t, tp, tf.Enum(`foo`, `f`, `foobar`, `x`))
	assert.NotEqual(t, tp, tf.Enum(`foo`, `foobar`))
}

func TestCiEnum(t *testing.T) {
	tp := tf.CiEnum()
	assert.Equal(t, tp, tf.Not(typ.Any))

	tp = tf.CiEnum(`f`, `foo`, `foobar`)
	assert.Instance(t, tp, `FOO`)
	assert.NotInstance(t, tp, `Oops`)
	assert.Assignable(t, typ.String, tp)
	assert.NotAssignable(t, tf.String(3, 3), tp)
	assert.Assignable(t, tf.String(1, 6), tp)
	assert.Assignable(t, tp, tf.CiEnum(`foo`))
	assert.NotAssignable(t, tf.Pattern(regexp.MustCompile(`f`)), tp)
	assert.Assignable(t, tp, tf.Enum(`f`, `foo`, `foobar`))

	assert.Equal(t, `~"f"|~"foo"|~"foobar"`, tp.String())
}

func TestIntEnum(t *testing.T) {
	tp := tf.IntEnum()
	assert.Equal(t, tp, tf.Not(typ.Any))

	tp = tf.IntEnum(2, 4, 8)
	assert.Instance(t, tp, 4)
	assert.NotInstance(t, tp, 5)
	assert.Assignable(t, typ.Integer, tp)
	assert.NotAssignable(t, tf.Integer64(3, 3, true), tp)
	assert.Assignable(t, tf.Integer64(2, 8, true), tp)
	assert.Assignable(t, tp, tf.IntEnum(8))
	assert.Assignable(t, tp, tf.IntEnum(2, 4))
	assert.Assignable(t, tp, tf.IntEnum(8, 2, 4))

	assert.Equal(t, `2|4|8`, tp.String())
}
