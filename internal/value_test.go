package internal_test

import (
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/test/require"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestValue(t *testing.T) {
	s := vf.String(`a`)
	assert.Same(t, s, vf.Value(s))
	assert.True(t, vf.True == vf.Value(true))
	assert.True(t, vf.False == vf.Value(false))
	assert.True(t, vf.Value([]dgo.Value{s}).(dgo.Array).Frozen())
	assert.True(t, vf.Value([]string{`a`}).(dgo.Array).Frozen())
	assert.Equal(t, vf.Value([]dgo.Value{s}), vf.Value([]string{`a`}))
	assert.True(t, vf.Value([]int{1}).(dgo.Array).Frozen())

	v := vf.Value(regexp.MustCompile(`.*`))
	_, ok := v.(dgo.Regexp)
	require.True(t, ok)

	v = vf.Value(int8(42))
	i, ok := v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(int16(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(int32(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(int64(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(uint8(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(uint16(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(uint32(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(uint(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(uint64(42))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(uint(math.MaxUint64))
	i, ok = v.(dgo.BigInt)
	require.True(t, ok)
	assert.Equal(t, v, new(big.Int).SetUint64(math.MaxUint64))

	v = vf.Value(uint64(math.MaxUint64))
	i, ok = v.(dgo.BigInt)
	require.True(t, ok)
	assert.Equal(t, v, new(big.Int).SetUint64(math.MaxUint64))

	v = vf.Value(big.NewInt(123))
	bi, ok := v.(dgo.BigInt)
	require.True(t, ok)
	assert.True(t, bi.GoBigInt().Cmp(big.NewInt(123)) == 0)

	v = vf.Value(float32(3.14))
	f, ok := v.(dgo.Float)
	require.True(t, ok)
	assert.True(t, float32(3.14) == float32(f.GoFloat()))

	v = vf.Value(3.14)
	f, ok = v.(dgo.Float)
	require.True(t, ok)
	assert.True(t, f.GoFloat() == 3.14)

	v = vf.Value(big.NewFloat(3.14))
	bf, ok := v.(dgo.BigFloat)
	require.True(t, ok)
	assert.True(t, bf.GoBigFloat().Cmp(big.NewFloat(3.14)) == 0)

	v = vf.Value(struct{ A int }{10})
	assert.Equal(t, struct{ A int }{10}, v)

	assert.Equal(t, 42, vf.Value(json.Number(`42`)))
	assert.Equal(t, 3.14, vf.Value(json.Number(`3.14`)))
	assert.Panic(t, func() { vf.Value(json.Number(`not a float`)) }, `invalid`)
}

func TestValue_reflected(t *testing.T) {
	s := vf.String(`a`)
	assert.True(t, vf.Nil == vf.Value(reflect.ValueOf(nil)))
	assert.True(t, vf.Nil == vf.Value(reflect.ValueOf(([]string)(nil))))
	assert.True(t, vf.Nil == vf.Value(reflect.ValueOf((map[string]string)(nil))))
	assert.True(t, vf.Nil == vf.Value(reflect.ValueOf((*string)(nil))))

	assert.True(t, vf.True == vf.Value(reflect.ValueOf(true)))
	assert.True(t, vf.False == vf.Value(reflect.ValueOf(false)))
	assert.Same(t, s, vf.Value(reflect.ValueOf(s)))
	assert.True(t, vf.Value(reflect.ValueOf([]dgo.Value{s})).(dgo.Array).Frozen())
	assert.True(t, vf.Value(reflect.ValueOf([]string{`a`})).(dgo.Array).Frozen())
	assert.Equal(t, vf.Value(reflect.ValueOf([]dgo.Value{s})), vf.Value([]string{`a`}))
	assert.True(t, vf.Value(reflect.ValueOf([]int{1})).(dgo.Array).Frozen())

	v := vf.Value(reflect.ValueOf(regexp.MustCompile(`.*`)))
	_, ok := v.(dgo.Regexp)
	require.True(t, ok)

	v = vf.Value(reflect.ValueOf(int8(42)))
	i, ok := v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(int16(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(int32(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(int64(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint8(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint16(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint32(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint64(42)))
	i, ok = v.(dgo.Integer)
	require.True(t, ok)
	assert.True(t, i.GoInt() == 42)

	v = vf.Value(reflect.ValueOf(uint(math.MaxUint64)))
	i, ok = v.(dgo.BigInt)
	require.True(t, ok)
	assert.Equal(t, v, new(big.Int).SetUint64(math.MaxUint64))

	v = vf.Value(reflect.ValueOf(uint64(math.MaxUint64)))
	i, ok = v.(dgo.BigInt)
	require.True(t, ok)
	assert.Equal(t, v, new(big.Int).SetUint64(math.MaxUint64))

	v = vf.Value(reflect.ValueOf(float32(3.14)))
	f, ok := v.(dgo.Float)
	require.True(t, ok)
	assert.True(t, float32(3.14) == float32(f.GoFloat()))

	v = vf.Value(reflect.ValueOf(3.14))
	f, ok = v.(dgo.Float)
	require.True(t, ok)
	assert.True(t, f.GoFloat() == 3.14)

	v = vf.Value(reflect.ValueOf(reflect.ValueOf))
	_, ok = v.(dgo.Function)
	require.True(t, ok)

	v = vf.Value([]interface{}{map[string]interface{}{`a`: 1}})
	assert.Equal(t, []map[string]int{{`a`: 1}}, v)
}

func TestFromValue(t *testing.T) {
	v := vf.Integer(32)
	var vc int
	vf.FromValue(v, &vc)
	assert.Equal(t, v, vc)
}

func TestFromValue_notPointer(t *testing.T) {
	v := vf.Integer(32)
	var vc int
	assert.Panic(t, func() { vf.FromValue(v, vc) }, `not a pointer`)
}

func TestNew(t *testing.T) {
	assert.Same(t, vf.Nil, vf.New(typ.Any, vf.Nil))
	assert.Same(t, vf.Nil, vf.New(typ.Any, vf.Arguments(vf.Nil)))
	assert.Panic(t, func() { vf.New(typ.Any, vf.Arguments(vf.Nil, vf.Nil)) }, `unable to create`)
	assert.Panic(t, func() { vf.New(tf.Not(typ.Nil), vf.Nil) }, `unable to create`)
}
