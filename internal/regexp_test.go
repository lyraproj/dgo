package internal_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestRegexpDefault(t *testing.T) {
	tp := typ.Regexp
	r := regexp.MustCompile(`[a-z]+`)
	assert.Instance(t, tp, r)
	assert.True(t, tp.IsInstance(r))
	assert.NotInstance(t, tp, `r`)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.String)
	assert.Equal(t, tp, tp)
	assert.NotEqual(t, tp, typ.String)
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Instance(t, tp.Type(), tp)
	assert.Equal(t, `regexp`, tp.String())
	assert.True(t, reflect.ValueOf(r).Type().AssignableTo(tp.ReflectType()))
}

func TestRegexpExact(t *testing.T) {
	rx := regexp.MustCompile(`[a-z]+`)
	r := vf.Value(rx)
	tp := r.Type().(dgo.RegexpType)
	assert.Instance(t, tp, r)
	assert.True(t, tp.IsInstance(rx))
	assert.NotInstance(t, tp, regexp.MustCompile(`[a-z]*`))
	assert.NotInstance(t, tp, `[a-z]*`)
	assert.Assignable(t, typ.Regexp, tp)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.Regexp)
	assert.Equal(t, tp, tp)
	assert.NotEqual(t, tp, typ.Regexp)
	assert.NotEqual(t, tp, typ.String)
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Instance(t, tp.Type(), tp)
	assert.Same(t, typ.Regexp, typ.Generic(tp))
	assert.Equal(t, `regexp "[a-z]+"`, tp.String())
	assert.Equal(t, typ.Regexp.ReflectType(), tp.ReflectType())
}

func TestRegexp(t *testing.T) {
	rx := regexp.MustCompile(`[a-z]+`)
	assert.Equal(t, rx, regexp.MustCompile(`[a-z]+`))
	assert.NotEqual(t, rx, regexp.MustCompile(`[a-z]*`))
	assert.Equal(t, rx, vf.Value(regexp.MustCompile(`[a-z]+`)))
	assert.NotEqual(t, rx, vf.Value(regexp.MustCompile(`[a-z]*`)))
	assert.NotEqual(t, rx, `[a-z]*`)
	assert.Same(t, rx, vf.Value(rx).(dgo.Regexp).GoRegexp())
}

func TestRegexp_ReflectTo(t *testing.T) {
	var ex *regexp.Regexp
	v := vf.Value(regexp.MustCompile(`[a-z]+`))
	vf.ReflectTo(v, reflect.ValueOf(&ex).Elem())
	assert.NotNil(t, ex)
	assert.Equal(t, v, ex)

	var ev regexp.Regexp
	vf.ReflectTo(v, reflect.ValueOf(&ev).Elem())
	assert.Equal(t, v, &ev)

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	ec, ok := mi.(*regexp.Regexp)
	require.True(t, ok)
	assert.Same(t, ex, ec)
}
