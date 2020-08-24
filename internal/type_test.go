package internal_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestFromReflected(t *testing.T) {
	v := tf.FromReflected(reflect.ValueOf(regexp.MustCompile(`a`)).Type())
	assert.Assignable(t, typ.Regexp, v)

	v = tf.FromReflected(reflect.ValueOf([]string{`a`}).Type())
	assert.Assignable(t, tf.Array(typ.String), v)

	v = tf.FromReflected(reflect.ValueOf(map[int]string{1: `a`}).Type())
	assert.Assignable(t, tf.Map(typ.Integer, typ.String), v)

	v = tf.FromReflected(reflect.ValueOf(map[int]interface{}{1: `a`}).Type())
	assert.Assignable(t, tf.Map(typ.Integer, typ.Any), v)

	v = tf.FromReflected(reflect.ValueOf(&map[int]string{1: `a`}).Type())
	assert.Assignable(t, v, tf.Map(typ.Integer, typ.String))
	assert.Assignable(t, v, typ.Nil)

	v = tf.FromReflected(reflect.ValueOf(struct{ A int }{3}).Type())
	assert.Assignable(t, typ.Native, v)
}

func TestGeneric(t *testing.T) {
	assert.Same(t, typ.Generic(typ.String), typ.String)
	assert.NotEqual(t, typ.Generic(typ.String), tf.String(10))
	assert.Same(t, typ.String, typ.Generic(vf.String(`hello`).Type()))
}
