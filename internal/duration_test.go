package internal_test

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestDurationDefault(t *testing.T) {
	tp := typ.Duration
	dv := time.Second
	assert.Instance(t, tp, dv)
	assert.NotInstance(t, tp, `r`)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.String)

	assert.Equal(t, tp, tp)
	assert.NotEqual(t, tp, typ.String)

	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())

	assert.Instance(t, tp.Type(), tp)

	assert.Equal(t, `duration`, tp.String())

	assert.True(t, reflect.ValueOf(dv).Type().AssignableTo(tp.ReflectType()))
}

func TestDurationType_New(t *testing.T) {
	dv, _ := time.ParseDuration(`2h45m10s`)
	tv := vf.Duration(dv)
	assert.Same(t, tv, vf.New(typ.Duration, tv))
	assert.Same(t, tv, vf.New(tv.Type(), tv))
	assert.Same(t, tv, vf.New(tv.Type(), vf.Arguments(tv)))
	assert.Equal(t, tv, vf.New(typ.Duration, vf.Int64(int64(tv.GoDuration()))))
	assert.Equal(t, tv, vf.New(typ.Duration, vf.Float(tv.SecondsWithFraction())))

	assert.Equal(t, tv, vf.New(typ.Duration, vf.String(`2h45m10s`)))

	assert.Panic(t, func() { vf.New(tv.Type(), vf.String(`3h`)) }, `cannot be assigned`)
	assert.Panic(t, func() { vf.New(typ.Duration, vf.Sensitive(5)) }, `illegal argument`)
}

func TestDurationExact(t *testing.T) {
	dv, _ := time.ParseDuration(`2h45m10s`)
	ts := vf.Value(dv).(dgo.Duration)
	tp := ts.Type()
	assert.Instance(t, tp, ts)
	assert.NotInstance(t, tp, dv+1)
	assert.NotInstance(t, tp, dv.String())
	assert.Assignable(t, typ.Duration, tp)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.Time)

	assert.Equal(t, tp, tp)
	assert.NotEqual(t, tp, typ.Time)
	assert.NotEqual(t, tp, typ.String)

	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())

	assert.Instance(t, tp.Type(), tp)

	assert.Same(t, typ.Duration, typ.Generic(tp))

	assert.Equal(t, `duration "`+dv.String()+`"`, tp.String())
	assert.Equal(t, typ.Duration.ReflectType(), tp.ReflectType())

	n, ok := ts.(dgo.Number)
	require.True(t, ok)
	i, ok := n.ToInt()
	require.True(t, ok)
	assert.Equal(t, int64(dv), i)
	u, ok := n.ToUint()
	require.True(t, ok)
	assert.Equal(t, i, u)
	neg := vf.Value(time.Duration(-1)).(dgo.Number)
	_, ok = neg.ToInt()
	require.True(t, ok)
	_, ok = neg.ToUint()
	require.False(t, ok)

	f, ok := n.ToFloat()
	require.True(t, ok)
	assert.Equal(t, ts.SecondsWithFraction(), f)
	assert.Equal(t, n.Float(), f)

	assert.Equal(t, big.NewInt(n.Integer().GoInt()), n.ToBigInt())
	assert.Equal(t, big.NewFloat(ts.SecondsWithFraction()), n.ToBigFloat())
}

func TestDurationFromString_bad(t *testing.T) {
	assert.Panic(t, func() { vf.DurationFromString(`4h32r`) }, `unknown unit "r" in duration "4h32r"`)
}

func TestDuration_ReflectTo(t *testing.T) {
	var ex time.Duration
	dv, _ := time.ParseDuration(`2h45m10s`)
	v := vf.Duration(dv)
	vf.ReflectTo(v, reflect.ValueOf(&ex).Elem())
	assert.NotNil(t, ex)
	assert.Equal(t, v, ex)

	var ev time.Duration
	vf.ReflectTo(v, reflect.ValueOf(&ev).Elem())
	assert.Equal(t, v, ev)

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	ec, ok := mi.(time.Duration)
	require.True(t, ok)
	assert.Same(t, ex, ec)
}
