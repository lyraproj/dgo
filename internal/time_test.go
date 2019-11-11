package internal_test

import (
	"reflect"
	"testing"
	"time"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestTimeDefault(t *testing.T) {
	tp := typ.Time
	ts := time.Now()
	require.Instance(t, tp, ts)
	require.NotInstance(t, tp, `r`)
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.String)

	require.Equal(t, tp, tp)
	require.NotEqual(t, tp, typ.String)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, `time`, tp.String())

	require.True(t, reflect.ValueOf(ts).Type().AssignableTo(tp.ReflectType()))
}

func TestTimeType_New(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-11-11T17:57:00-00:00`)
	tv := vf.Time(ts)
	require.Same(t, tv, vf.New(typ.Time, tv))
	require.Same(t, tv, vf.New(tv.Type(), tv))
	require.Same(t, tv, vf.New(tv.Type(), vf.Arguments(tv)))
	require.Equal(t, tv, vf.New(typ.Time, vf.Integer(tv.GoTime().Unix())))
	require.Equal(t, tv, vf.New(typ.Time, vf.Float(tv.SecondsWithFraction())))

	require.Equal(t, tv, vf.New(typ.Time, vf.String(`2019-11-11T17:57:00-00:00`)))

	require.Panic(t, func() { vf.New(tv.Type(), vf.String(`2019-10-06T07:15:00-07:00`)) }, `cannot be assigned`)
	require.Panic(t, func() { vf.New(typ.Time, vf.Sensitive(5)) }, `illegal argument`)
}

func TestTimeExact(t *testing.T) {
	now := time.Now()
	ts := vf.Value(now)
	tp := ts.Type()
	require.Instance(t, tp, ts)
	require.NotInstance(t, tp, now.Add(1))
	require.NotInstance(t, tp, now.String())
	require.Assignable(t, typ.Time, tp)
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.Time)

	require.Equal(t, tp, tp)
	require.NotEqual(t, tp, typ.Time)
	require.NotEqual(t, tp, typ.String)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Instance(t, tp.Type(), tp)

	require.Same(t, typ.Time, typ.Generic(tp))

	require.Equal(t, `time["`+now.Format(time.RFC3339Nano)+`"]`, tp.String())
	require.Equal(t, typ.Time.ReflectType(), tp.ReflectType())
}

func TestTime(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	zt, _ := time.Parse(time.RFC3339, `2019-10-06T16:15:00+02:00`)
	ot, _ := time.Parse(time.RFC3339, `2019-10-06T16:13:00+00:00`)
	met, _ := time.LoadLocation(`Europe/Zurich`)

	v := vf.Time(ts)
	require.Equal(t, v, v)
	require.Equal(t, v, vf.Time(zt.In(met)))
	require.NotEqual(t, v, ot)
	require.NotEqual(t, v, `2019-10-06:16:15:00-01:00`)
	require.Equal(t, v, v.GoTime())

	ts, _ = time.Parse(time.RFC3339, `1628-08-10T17:15:00.123456-07:00`)
	ot = ts.AddDate(400, 0, 0)
	v = vf.Time(ts)
	v2 := vf.Time(ot)
	dv := ot.Unix() - ts.Unix()
	require.Equal(t, float64(dv), v2.SecondsWithFraction()-v.SecondsWithFraction())
}

func TestTimeFromString(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	v := vf.TimeFromString(`2019-10-06T07:15:00-07:00`)
	require.Equal(t, v, ts)
}

func TestTimeFromString_bad(t *testing.T) {
	require.Panic(t, func() { vf.TimeFromString(`2019-13-06T07:15:00-07:00`) }, `month out of range`)
}

func TestTime_ReflectTo(t *testing.T) {
	var ex *time.Time
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T16:15:00-01:00`)
	v := vf.Time(ts)
	vf.ReflectTo(v, reflect.ValueOf(&ex).Elem())
	require.NotNil(t, ex)
	require.Equal(t, v, ex)

	var ev time.Time
	vf.ReflectTo(v, reflect.ValueOf(&ev).Elem())
	require.Equal(t, v, &ev)

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	ec, ok := mi.(*time.Time)
	require.True(t, ok)
	require.Same(t, ex, ec)
}
