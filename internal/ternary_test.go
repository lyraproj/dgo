package internal_test

import (
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestAllOfType(t *testing.T) {
	tp := newtype.AllOf(newtype.Enum(`a`, `b`, `c`), newtype.Enum(`b`, `c`, `d`))
	require.Instance(t, tp, `b`)
	require.Instance(t, tp, `c`)
	require.NotInstance(t, tp, `a`)
	require.Assignable(t, typ.AllOf, tp)
	require.NotAssignable(t, tp, typ.AllOf)
	require.Assignable(t, tp, newtype.AllOf(newtype.Enum(`c`, `b`, `a`), newtype.Enum(`b`, `d`, `c`)))
	require.Assignable(t, tp, newtype.AllOf(newtype.Enum(`a`, `b`, `c`), newtype.Enum(`b`, `c`, `d`), newtype.Enum(`c`, `d`, `e`)))
	require.NotAssignable(t, tp, newtype.AllOf(newtype.Enum(`a`, `b`, `c`), newtype.Enum(`e`, `c`, `d`)))

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, tp, newtype.AllOf(newtype.Enum(`b`, `c`, `d`), newtype.Enum(`a`, `b`, `c`)))
	require.NotEqual(t, tp, newtype.AllOf(newtype.Enum(`b`, `c`), newtype.Enum(`a`, `b`, `c`)))
	require.NotEqual(t, tp, newtype.Enum(`a`, `b`, `c`))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Same(t, newtype.AllOf(), typ.Any)
	require.Same(t, newtype.AllOf(typ.String), typ.String)

	require.Equal(t, tp.(dgo.TernaryType).Operands(), vf.Values(newtype.Enum(`a`, `b`, `c`), newtype.Enum(`b`, `c`, `d`)))
	require.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpAnd)
	require.Equal(t, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType().(dgo.TernaryType).Operator(), dgo.OpAnd)

	require.Equal(t, `("a"|"b"|"c")&("b"|"c"|"d")`, tp.String())

	tp = newtype.AllOf(newtype.Pattern(regexp.MustCompile(`a`)), newtype.Pattern(regexp.MustCompile(`b`)), newtype.Pattern(regexp.MustCompile(`c`)))
	require.Instance(t, tp, `abc`)
	require.NotInstance(t, tp, `c`)

	require.Assignable(t, tp, newtype.AllOf(newtype.Pattern(regexp.MustCompile(`a`)), newtype.Pattern(regexp.MustCompile(`b`)), newtype.Pattern(regexp.MustCompile(`c`))))
	require.NotAssignable(t, tp, newtype.AllOf(newtype.Pattern(regexp.MustCompile(`a`)), newtype.Pattern(regexp.MustCompile(`b`))))
	require.Assignable(t, tp, newtype.AllOf(newtype.Pattern(regexp.MustCompile(`a`)), vf.Value(`abc`).Type()))
	require.NotAssignable(t, tp, newtype.AllOf(newtype.Pattern(regexp.MustCompile(`a`)), vf.Value(`abc`).Type(), typ.String))

	require.Assignable(t, tp, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType())
	require.NotAssignable(t, vf.Values("abc", "bca").Type().(dgo.ArrayType).ElementType(), tp)
	require.Assignable(t, tp, newtype.Parse(`"abc"^/a/&/b/&/c/`))
	require.Assignable(t, newtype.Parse(`/a/|/b/|/c/`), tp)
}

func TestAnyOfType(t *testing.T) {
	tp := newtype.AnyOf(typ.Integer, typ.String)
	require.Instance(t, tp, `hello`)
	require.Instance(t, tp, 3)
	require.NotInstance(t, tp, true)
	require.Assignable(t, typ.AnyOf, tp)
	require.NotAssignable(t, tp, typ.AnyOf)
	require.Assignable(t, tp, newtype.AnyOf(typ.String, typ.Integer))
	require.Assignable(t, newtype.AnyOf(typ.String, typ.Integer, typ.Boolean), tp)
	require.NotAssignable(t, tp, newtype.AnyOf(typ.String, typ.Integer, typ.Boolean))
	require.NotAssignable(t, newtype.AnyOf(typ.String, typ.Boolean), tp)

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, tp, newtype.AnyOf(typ.Integer, typ.String))
	require.NotEqual(t, tp, newtype.AnyOf(typ.Integer, typ.Boolean))
	require.NotEqual(t, tp, typ.Integer)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, newtype.AnyOf(), newtype.Not(typ.Any))
	require.Same(t, newtype.AnyOf(typ.String), typ.String)

	require.Equal(t, tp.(dgo.TernaryType).Operands(), vf.Values(typ.Integer, typ.String))
	require.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpOr)

	require.Equal(t, `int|string`, tp.String())
}

func TestOneOfType(t *testing.T) {
	tp := newtype.OneOf(typ.Integer, newtype.Pattern(regexp.MustCompile(`a`)), newtype.Pattern(regexp.MustCompile(`b`)))
	require.Instance(t, tp, `a`)
	require.Instance(t, tp, 3)
	require.NotInstance(t, tp, `ab`)
	require.NotInstance(t, tp, true)
	require.Assignable(t, typ.OneOf, tp)
	require.NotAssignable(t, tp, typ.OneOf)
	require.NotAssignable(t, tp, newtype.OneOf(typ.String, typ.Integer))
	require.Assignable(t, tp, newtype.OneOf(vf.String(`a`).Type(), vf.String(`b`).Type(), typ.Integer))
	require.NotAssignable(t, tp, newtype.OneOf(vf.String(`ab`).Type(), typ.Integer))
	require.Assignable(t, newtype.OneOf(typ.String, typ.Integer), tp)
	require.Assignable(t, newtype.AnyOf(typ.String, typ.Integer, typ.Boolean), tp)
	require.NotAssignable(t, tp, newtype.AnyOf(typ.String, typ.Integer, typ.Boolean))
	require.NotAssignable(t, newtype.AnyOf(typ.String, typ.Boolean), tp)

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, tp, newtype.OneOf(typ.Integer, newtype.Pattern(regexp.MustCompile(`a`)), newtype.Pattern(regexp.MustCompile(`b`))))
	require.NotEqual(t, tp, newtype.OneOf(typ.Integer, typ.Boolean))
	require.NotEqual(t, tp, typ.Integer)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, newtype.OneOf(), newtype.Not(typ.Any))
	require.Same(t, newtype.OneOf(typ.String), typ.String)

	require.Equal(t, tp.(dgo.TernaryType).Operands(), vf.Values(typ.Integer, newtype.Pattern(regexp.MustCompile(`a`)), newtype.Pattern(regexp.MustCompile(`b`))))
	require.Equal(t, tp.(dgo.TernaryType).Operator(), dgo.OpOne)

	require.Equal(t, `int^/a/^/b/`, tp.String())
}

func TestEnum(t *testing.T) {
	tp := newtype.Enum()
	require.Equal(t, tp, newtype.Not(typ.Any))

	tp = newtype.Enum(`f`, `foo`, `foobar`)
	require.Instance(t, tp, `foo`)
	require.NotInstance(t, tp, `oops`)
	require.Assignable(t, typ.String, tp)
	require.NotAssignable(t, newtype.String(3, 3), tp)
	require.Assignable(t, newtype.String(1, 6), tp)
	require.Assignable(t, newtype.Pattern(regexp.MustCompile(`f`)), tp)
	require.NotAssignable(t, newtype.Pattern(regexp.MustCompile(`o`)), tp)
	require.Equal(t, tp, newtype.Enum(`foo`, `f`, `foobar`))
	require.NotEqual(t, tp, newtype.Enum(`foo`, `f`, `foobar`, `x`))
	require.NotEqual(t, tp, newtype.Enum(`foo`, `foobar`))
}
