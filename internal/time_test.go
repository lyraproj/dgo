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

func TestTimeDefault(t *testing.T) {
	tp := typ.Time
	ts := time.Now()
	assert.Instance(t, tp, ts)
	assert.NotInstance(t, tp, `r`)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.String)
	assert.Equal(t, tp, tp)
	assert.NotEqual(t, tp, typ.String)
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Instance(t, tp.Type(), tp)
	assert.Equal(t, `time`, tp.String())
	assert.True(t, reflect.ValueOf(ts).Type().AssignableTo(tp.ReflectType()))
}

func TestTimeType_New(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-11-11T17:57:00-00:00`)
	tv := vf.Time(ts)
	assert.Same(t, tv, vf.New(typ.Time, tv))
	assert.Same(t, tv, vf.New(tv.Type(), tv))
	assert.Same(t, tv, vf.New(tv.Type(), vf.Arguments(tv)))
	assert.Equal(t, tv, vf.New(typ.Time, vf.Integer(tv.GoTime().Unix())))
	assert.Equal(t, tv, vf.New(typ.Time, vf.Float(tv.SecondsWithFraction())))
	assert.Equal(t, tv, vf.New(typ.Time, vf.String(`2019-11-11T17:57:00-00:00`)))
	assert.Panic(t, func() { vf.New(tv.Type(), vf.String(`2019-10-06T07:15:00-07:00`)) }, `cannot be assigned`)
	assert.Panic(t, func() { vf.New(typ.Time, vf.Sensitive(5)) }, `illegal argument`)
}

func TestTimeExact(t *testing.T) {
	now := time.Now()
	ts := vf.Value(now).(dgo.Time)
	tp := ts.Type().(dgo.TimeType)
	assert.Instance(t, tp, ts)
	assert.True(t, tp.IsInstance(now))
	assert.NotInstance(t, tp, now.Add(1))
	assert.NotInstance(t, tp, now.String())
	assert.Assignable(t, typ.Time, tp)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.Time)
	assert.Equal(t, tp, tp)
	assert.NotEqual(t, tp, typ.Time)
	assert.NotEqual(t, tp, typ.String)
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Instance(t, tp.Type(), tp)
	assert.Same(t, typ.Time, typ.Generic(tp))
	assert.Equal(t, `time "`+now.Format(time.RFC3339Nano)+`"`, tp.String())
	assert.Equal(t, typ.Time.ReflectType(), tp.ReflectType())

	n, ok := ts.(dgo.Number)
	require.True(t, ok)
	i, ok := n.ToInt()
	require.True(t, ok)
	assert.Equal(t, now.Unix(), i)

	f, ok := n.ToFloat()
	require.True(t, ok)
	assert.Equal(t, ts.SecondsWithFraction(), f)

	assert.Equal(t, big.NewInt(now.Unix()), n.ToBigInt())
	assert.Equal(t, big.NewFloat(ts.SecondsWithFraction()), n.ToBigFloat())
}

func TestTime_wrapped(t *testing.T) {
	now := time.Now()
	ts1, ok := vf.Value(now).(dgo.Time)
	require.True(t, ok)
	ts2, ok := vf.Value(&now).(dgo.Time)
	require.True(t, ok)
	assert.Equal(t, ts1, ts2)
}

func TestTime(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	zt, _ := time.Parse(time.RFC3339, `2019-10-06T16:15:00+02:00`)
	ot, _ := time.Parse(time.RFC3339, `2019-10-06T16:13:00+00:00`)
	met, _ := time.LoadLocation(`Europe/Zurich`)

	v := vf.Time(ts)
	assert.Equal(t, v, v)
	assert.Equal(t, v, vf.Time(zt.In(met)))
	assert.NotEqual(t, v, ot)
	assert.NotEqual(t, v, `2019-10-06:16:15:00-01:00`)
	assert.Equal(t, v, v.GoTime())

	ts, _ = time.Parse(time.RFC3339, `1628-08-10T17:15:00.123456-07:00`)
	ot = ts.AddDate(400, 0, 0)
	v = vf.Time(ts)
	v2 := vf.Time(ot)
	dv := ot.Unix() - ts.Unix()
	assert.Equal(t, float64(dv), v2.SecondsWithFraction()-v.SecondsWithFraction())
}

func TestTimeFromString(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	v := vf.TimeFromString(`2019-10-06T07:15:00-07:00`)
	assert.Equal(t, v, ts)
}

func TestTimeFromString_bad(t *testing.T) {
	assert.Panic(t, func() { vf.TimeFromString(`2019-13-06T07:15:00-07:00`) }, `month out of range`)
}

func TestTime_ReflectTo(t *testing.T) {
	var ex *time.Time
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T16:15:00-01:00`)
	v := vf.Time(ts)
	vf.ReflectTo(v, reflect.ValueOf(&ex).Elem())
	assert.NotNil(t, ex)
	assert.Equal(t, v, ex)

	var ev time.Time
	vf.ReflectTo(v, reflect.ValueOf(&ev).Elem())
	assert.Equal(t, v, &ev)

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	ec, ok := mi.(*time.Time)
	require.True(t, ok)
	assert.Same(t, ex, ec)
}
