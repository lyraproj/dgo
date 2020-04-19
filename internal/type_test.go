package internal_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/tada/dgo/vf"

	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
)

func TestFromReflected(t *testing.T) {
	v := tf.FromReflected(reflect.ValueOf(regexp.MustCompile(`a`)).Type())
	require.Assignable(t, typ.Regexp, v)

	v = tf.FromReflected(reflect.ValueOf([]string{`a`}).Type())
	require.Assignable(t, tf.Array(typ.String), v)

	v = tf.FromReflected(reflect.ValueOf(map[int]string{1: `a`}).Type())
	require.Assignable(t, tf.Map(typ.Integer, typ.String), v)

	v = tf.FromReflected(reflect.ValueOf(map[int]interface{}{1: `a`}).Type())
	require.Assignable(t, tf.Map(typ.Integer, typ.Any), v)

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
