package internal_test

import (
	"reflect"
	"regexp"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
)

func TestFromReflected(t *testing.T) {
	v := newtype.FromReflected(reflect.ValueOf(regexp.MustCompile(`a`)).Type())
	require.Assignable(t, typ.Regexp, v)

	v = newtype.FromReflected(reflect.ValueOf([]string{`a`}).Type())
	require.Assignable(t, newtype.Array(typ.String), v)

	v = newtype.FromReflected(reflect.ValueOf(map[int]string{1: `a`}).Type())
	require.Assignable(t, newtype.Map(typ.Integer, typ.String), v)

	v = newtype.FromReflected(reflect.ValueOf(&map[int]string{1: `a`}).Type())
	require.Assignable(t, v, newtype.Map(typ.Integer, typ.String))
	require.Assignable(t, v, typ.Nil)

	v = newtype.FromReflected(reflect.ValueOf(struct{ A int }{3}).Type())
	require.Assignable(t, typ.Native, v)
}
