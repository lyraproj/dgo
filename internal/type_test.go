package internal_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/vf"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
)

func TestFromReflected(t *testing.T) {
	v := tf.FromReflected(reflect.ValueOf(regexp.MustCompile(`a`)).Type())
	require.Assignable(t, typ.Regexp, v)

	v = tf.FromReflected(reflect.ValueOf([]string{`a`}).Type())
	require.Assignable(t, tf.Array(typ.String), v)

	v = tf.FromReflected(reflect.ValueOf(map[int]string{1: `a`}).Type())
	require.Assignable(t, tf.Map(typ.Integer, typ.String), v)

	v = tf.FromReflected(reflect.ValueOf(&map[int]string{1: `a`}).Type())
	require.Assignable(t, v, tf.Map(typ.Integer, typ.String))
	require.Assignable(t, v, typ.Nil)

	v = tf.FromReflected(reflect.ValueOf(struct{ A int }{3}).Type())
	require.Assignable(t, typ.Native, v)
}

func TestGeneric(t *testing.T) {
	require.Same(t, typ.Generic(typ.String), typ.String)
	require.NotEqual(t, typ.Generic(typ.String), tf.String(10))
	require.Same(t, typ.String, typ.Generic(vf.String(`hello`).Type()))
}
