package internal_test

import (
	"encoding/json"
	"math"
	"reflect"
	"regexp"
	"testing"

	"github.com/tada/dgo/dgo"
	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestValue(t *testing.T) {
	s := vf.String(`a`)
	require.Same(t, s, vf.Value(s))
	require.True(t, vf.True == vf.Value(true))
	require.True(t, vf.False == vf.Value(false))
	require.True(t, vf.Value([]dgo.Value{s}).(dgo.Array).Frozen())
	require.True(t, vf.Value([]string{`a`}).(dgo.Array).Frozen())
	require.Equal(t, vf.Value([]dgo.Value{s}), vf.Value([]string{`a`}))
	require.True(t, vf.Value([]int{1}).(dgo.Array).Frozen())

	v := vf.Value(regexp.MustCompile(`.*`))
	_, ok := v.(dgo.Regexp)
	require.True(t, ok)

	v = vf.Value(int8(42))
	i, ok := v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(int16(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(int32(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(int64(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(uint8(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(uint16(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(uint32(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(uint(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(uint64(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	require.Panic(t, func() { vf.Value(uint(math.MaxUint64)) }, `overflows`)
	require.Panic(t, func() { vf.Value(uint64(math.MaxUint64)) }, `overflows`)

	v = vf.Value(float32(3.14))
	f, ok := v.(dgo.Float)
	require.True(t, ok)
	require.True(t, float32(3.14) == float32(f.GoFloat()))

	v = vf.Value(3.14)
	f, ok = v.(dgo.Float)
	require.True(t, ok)
	require.True(t, f.GoFloat() == 3.14)

	v = vf.Value(struct{ A int }{10})
	require.Equal(t, struct{ A int }{10}, v)

	require.Equal(t, 42, vf.Value(json.Number(`42`)))
	require.Equal(t, 3.14, vf.Value(json.Number(`3.14`)))
	require.Panic(t, func() { vf.Value(json.Number(`not a float`)) }, `invalid`)
}

func TestValue_reflected(t *testing.T) {
	s := vf.String(`a`)
	require.True(t, vf.Nil == vf.Value(reflect.ValueOf(nil)))
	require.True(t, vf.Nil == vf.Value(reflect.ValueOf(([]string)(nil))))
	require.True(t, vf.Nil == vf.Value(reflect.ValueOf((map[string]string)(nil))))
	require.True(t, vf.Nil == vf.Value(reflect.ValueOf((*string)(nil))))

	require.True(t, vf.True == vf.Value(reflect.ValueOf(true)))
	require.True(t, vf.False == vf.Value(reflect.ValueOf(false)))
	require.Same(t, s, vf.Value(reflect.ValueOf(s)))
	require.True(t, vf.Value(reflect.ValueOf([]dgo.Value{s})).(dgo.Array).Frozen())
	require.True(t, vf.Value(reflect.ValueOf([]string{`a`})).(dgo.Array).Frozen())
	require.Equal(t, vf.Value(reflect.ValueOf([]dgo.Value{s})), vf.Value([]string{`a`}))
	require.True(t, vf.Value(reflect.ValueOf([]int{1})).(dgo.Array).Frozen())

	v := vf.Value(reflect.ValueOf(regexp.MustCompile(`.*`)))
	_, ok := v.(dgo.Regexp)
	require.True(t, ok)

	v = vf.Value(reflect.ValueOf(int8(42)))
	i, ok := v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(int16(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(int32(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(int64(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint8(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint16(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint32(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint64(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	require.True(t, i.GoInt() == 42)

	require.Panic(t, func() { vf.Value(reflect.ValueOf(uint(math.MaxUint64))) }, `overflows`)
	require.Panic(t, func() { vf.Value(reflect.ValueOf(uint64(math.MaxUint64))) }, `overflows`)

	v = vf.Value(reflect.ValueOf(float32(3.14)))
	f, ok := v.(dgo.Float)
	require.True(t, ok)
	require.True(t, float32(3.14) == float32(f.GoFloat()))

	v = vf.Value(reflect.ValueOf(3.14))
	f, ok = v.(dgo.Float)
	require.True(t, ok)
	require.True(t, f.GoFloat() == 3.14)

	v = vf.Value(reflect.ValueOf(reflect.ValueOf))
	_, ok = v.(dgo.Function)
	require.True(t, ok)

	v = vf.Value([]interface{}{map[string]interface{}{`a`: 1}})
	require.Equal(t, []map[string]int{{`a`: 1}}, v)
}

func TestFromValue(t *testing.T) {
	v := vf.Integer(32)
	var vc int
	vf.FromValue(v, &vc)
	require.Equal(t, v, vc)
}

func TestFromValue_notPointer(t *testing.T) {
	v := vf.Integer(32)
	var vc int
	require.Panic(t, func() { vf.FromValue(v, vc) }, `not a pointer`)
}

func TestNew(t *testing.T) {
	require.Same(t, vf.Nil, vf.New(typ.Any, vf.Nil))
	require.Same(t, vf.Nil, vf.New(typ.Any, vf.Arguments(vf.Nil)))
	require.Panic(t, func() { vf.New(typ.Any, vf.Arguments(vf.Nil, vf.Nil)) }, `unable to create`)
	require.Panic(t, func() { vf.New(tf.Not(typ.Nil), vf.Nil) }, `unable to create`)
}
