package internal_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestRegexpDefault(t *testing.T) {
	tp := typ.Regexp
	r := regexp.MustCompile(`[a-z]+`)
	require.Instance(t, tp, r)
	require.NotInstance(t, tp, `r`)
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.String)

	require.Equal(t, tp, tp)
	require.NotEqual(t, tp, typ.String)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, `regexp`, tp.String())

	require.True(t, reflect.ValueOf(r).Type().AssignableTo(tp.ReflectType()))
}

func TestRegexpExact(t *testing.T) {
	r := vf.Value(regexp.MustCompile(`[a-z]+`))
	tp := r.Type()
	require.Instance(t, tp, r)
	require.NotInstance(t, tp, regexp.MustCompile(`[a-z]*`))
	require.NotInstance(t, tp, `[a-z]*`)
	require.Assignable(t, typ.Regexp, tp)
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.Regexp)

	require.Equal(t, tp, tp)
	require.NotEqual(t, tp, typ.Regexp)
	require.NotEqual(t, tp, typ.String)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Instance(t, tp.Type(), tp)

	require.Equal(t, `regexp["[a-z]+"]`, tp.String())
	require.Equal(t, typ.Regexp.ReflectType(), tp.ReflectType())
}

func TestRegexp(t *testing.T) {
	rx := regexp.MustCompile(`[a-z]+`)
	require.Equal(t, rx, regexp.MustCompile(`[a-z]+`))
	require.NotEqual(t, rx, regexp.MustCompile(`[a-z]*`))
	require.Equal(t, rx, vf.Value(regexp.MustCompile(`[a-z]+`)))
	require.NotEqual(t, rx, vf.Value(regexp.MustCompile(`[a-z]*`)))
	require.NotEqual(t, rx, `[a-z]*`)
	require.Same(t, rx, vf.Value(rx).(dgo.Regexp).GoRegexp())
	require.NotEqual(t, vf.Float(3.14).ToFloat(), vf.Integer(3).ToFloat())
}

func TestRegexp_ReflectTo(t *testing.T) {
	var ex *regexp.Regexp
	v := vf.Value(regexp.MustCompile(`[a-z]+`))
	vf.ReflectTo(v, reflect.ValueOf(&ex).Elem())
	require.NotNil(t, ex)
	require.Equal(t, v, ex)

	var ev regexp.Regexp
	vf.ReflectTo(v, reflect.ValueOf(&ev).Elem())
	require.Equal(t, v, &ev)

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	ec, ok := mi.(*regexp.Regexp)
	require.True(t, ok)
	require.Same(t, ex, ec)
}
