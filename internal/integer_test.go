package internal_test

import (
	"math"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestInteger(t *testing.T) {
	require.Instance(t, typ.Integer, 3)
	require.NotInstance(t, typ.Integer, true)
	require.Assignable(t, typ.Integer, typ.Integer)
	require.Assignable(t, typ.Integer, newtype.IntegerRange(3, 5, true))
	require.Assignable(t, typ.Integer, vf.Integer(4).Type())
	require.Equal(t, typ.Integer, typ.Integer)
	require.Equal(t, typ.Integer, newtype.IntegerRange(math.MinInt64, math.MaxInt64, true))
	require.Instance(t, typ.Integer.Type(), typ.Integer)
	require.True(t, typ.Integer.IsInstance(1234))
	require.Equal(t, typ.Integer.Min(), math.MinInt64)
	require.Equal(t, typ.Integer.Max(), math.MaxInt64)
	require.True(t, typ.Integer.Inclusive())

	require.Equal(t, `int`, typ.Integer.String())
}

func TestIntegerExact(t *testing.T) {
	tp := vf.Integer(3).Type().(dgo.IntegerRangeType)
	require.Instance(t, tp, 3)
	require.NotInstance(t, tp, 2)
	require.NotInstance(t, tp, true)
	require.Assignable(t, newtype.IntegerRange(3, 5, true), tp)
	require.Assignable(t, tp, newtype.IntegerRange(3, 3, true))
	require.NotAssignable(t, tp, typ.Integer)
	require.Equal(t, tp, newtype.IntegerRange(3, 3, true))
	require.NotEqual(t, tp, newtype.IntegerRange(2, 5, true))
	require.Equal(t, tp.Min(), 3)
	require.Equal(t, tp.Max(), 3)
	require.True(t, tp.Inclusive())
	require.True(t, tp.IsInstance(3))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `3`, tp.String())

	require.Instance(t, tp.Type(), tp)
}

func TestIntegerRange(t *testing.T) {
	tp := newtype.IntegerRange(3, 5, true)
	require.Instance(t, tp, 3)
	require.NotInstance(t, tp, 2)
	require.NotInstance(t, tp, true)
	require.Assignable(t, tp, newtype.IntegerRange(3, 5, true))
	require.Assignable(t, tp, newtype.IntegerRange(4, 4, true))
	require.Assignable(t, tp, vf.Integer(4).Type())
	require.NotAssignable(t, tp, newtype.IntegerRange(2, 5, true))
	require.NotAssignable(t, tp, newtype.IntegerRange(3, 6, true))
	require.NotAssignable(t, tp, vf.Integer(6).Type())
	require.Equal(t, tp, newtype.IntegerRange(5, 3, true))
	require.NotEqual(t, tp, newtype.IntegerRange(2, 5, true))
	require.NotEqual(t, tp, newtype.IntegerRange(3, 4, true))
	require.NotEqual(t, tp, typ.Integer)
	require.Equal(t, tp.Min(), 3)
	require.Equal(t, tp.Max(), 5)
	require.Equal(t, vf.Integer(4).Type(), newtype.IntegerRange(4, 4, true))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `3..5`, tp.String())

	require.Instance(t, tp.Type(), tp)

	tp = newtype.IntegerRange(3, 5, false)
	require.Instance(t, tp, 4)
	require.NotInstance(t, tp, 5)
	require.Assignable(t, tp, newtype.IntegerRange(3, 5, false))
	require.NotAssignable(t, tp, newtype.IntegerRange(3, 5, true))
	require.Assignable(t, newtype.IntegerRange(3, 5, true), tp)
	require.Assignable(t, tp, newtype.IntegerRange(3, 4, true))
	require.Assignable(t, newtype.IntegerRange(3, 4, true), tp)

	require.Panic(t, func() { newtype.IntegerRange(4, 4, false) }, `cannot have equal min and max`)
}

func TestInteger_CompareTo(t *testing.T) {
	c, ok := vf.Integer(3).CompareTo(vf.Integer(3))
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.Integer(3).CompareTo(vf.Integer(-3))
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Integer(-3).CompareTo(vf.Integer(3))
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(3))
	require.True(t, ok)
	require.Equal(t, 0, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(3.1))
	require.True(t, ok)
	require.Equal(t, -1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Float(2.9))
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	c, ok = vf.Integer(3).CompareTo(vf.Nil)
	require.True(t, ok)
	require.Equal(t, 1, c)

	_, ok = vf.Integer(3).CompareTo(vf.True)
	require.False(t, ok)
}

func TestInteger_String(t *testing.T) {
	require.Equal(t, `1234`, vf.Integer(1234).String())
	require.Equal(t, `-4321`, vf.Integer(-4321).String())
}
