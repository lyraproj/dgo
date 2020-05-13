package streamer_test

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestNew(t *testing.T) {
	s := streamer.New(nil, nil)
	c := streamer.NewCollector()
	m := vf.Map(`a`, 1, `b`, 2)
	s.Stream(m, c)
	assert.Equal(t, m, c.Value())
}

func TestEncode_dedup(t *testing.T) {
	v := vf.Strings(`a`, `b`)
	a := vf.Values(v, v)
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(a, c)
	ac := c.Value().(dgo.Array)
	assert.Same(t, ac.Get(0), ac.Get(1))
}

func TestEncode_demoteDedup(t *testing.T) {
	a := vf.Values(vf.Map(`the key`, `a`), vf.Map(`the key`, `b`))
	b := bytes.Buffer{}
	o := streamer.DefaultOptions()
	o.DedupLevel = streamer.MaxDedup
	streamer.New(nil, o).Stream(a, streamer.JSON(&b))
	assert.Equal(t, `[{"the key":"a"},{"the key":"b"}]`, b.String())
}

func TestEncode_sensitive(t *testing.T) {
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(vf.Sensitive(`secret`), c)
	assert.Equal(t, vf.Map(`__type`, `sensitive`, `__value`, `secret`), c.Value())
}

func TestEncode_unknown(t *testing.T) {
	c := streamer.NewCollector()
	assert.Panic(t,
		func() { streamer.New(nil, nil).Stream(vf.Value(struct{ A string }{A: `hello`}), c) },
		`unable to serialize value`)
}

func TestEncode_type(t *testing.T) {
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(tf.String(1), c)
	assert.Equal(t, vf.Map(`__type`, `string[1]`), c.Value())
}

func TestEncode_bigfloat(t *testing.T) {
	v := vf.New(typ.BigFloat, vf.Arguments(`12345678901234567890123.3412`))
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(v, c)
	assert.Same(t, v, c.Value())
}

func TestEncode_bigfloat_data_128_bits(t *testing.T) {
	v := vf.New(typ.BigFloat, vf.Arguments(`1234567890123456789012345678912345.3412`, 128))
	c := streamer.DataCollector()
	streamer.New(nil, nil).Stream(v, c)
	m, ok := c.Value().(dgo.Map)
	require.True(t, ok)

	assert.Equal(t, vf.Map(`__type`, `bigfloat`, `__value`, `0x.f379bffe5ccb7a097341fa5b6d655d64p+110`), m)
	assert.Equal(t, v, vf.New(typ.BigFloat, vf.Arguments(m.Get(`__value`), 128)))
}

func TestEncode_bigint(t *testing.T) {
	v := vf.New(typ.BigInt, vf.Arguments(`12345678901234567890123`))
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(v, c)
	assert.Same(t, v, c.Value())

	v = vf.New(typ.BigInt, vf.Arguments(`123`))
	c = streamer.NewCollector()
	streamer.New(nil, nil).Stream(v, c)
	v2 := c.Value()
	assert.NotSame(t, v, v2)
	assert.Equal(t, 123, v2)
}

func TestEncode_binary(t *testing.T) {
	v := vf.BinaryFromString(`AQID`)
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(v, c)
	assert.Same(t, v, c.Value())
}

func TestEncode_time(t *testing.T) {
	v := vf.Time(time.Now())
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(v, c)
	assert.Same(t, v, c.Value())
}

func TestEncode_native(t *testing.T) {
	v := vf.Value(struct{ A string }{A: `a`})
	c := streamer.NewCollector()
	o := streamer.DefaultOptions()
	o.RichData = false
	assert.Panic(t, func() { streamer.New(nil, o).Stream(v, c) }, `unable to serialize`)
}

func TestEncode_alias(t *testing.T) {
	var tp dgo.Value
	am := tf.BuiltInAliases().Collect(func(aa dgo.AliasAdder) {
		tp = tf.ParseFile(aa, ``, `ne=string[1]`)
	})
	c := streamer.NewCollector()
	streamer.New(am, nil).Stream(tp, c)
	assert.Equal(t, vf.Map(`__type`, `alias`, `__value`, vf.Strings(`ne`, `string[1]`)), c.Value())
}

func TestEncode_noDedup(t *testing.T) {
	v := vf.Strings(`a`, `b`)
	a := vf.Values(v, v)
	c := streamer.NewCollector()
	o := streamer.DefaultOptions()
	o.DedupLevel = streamer.NoDedup
	streamer.New(nil, o).Stream(a, c)
	ac := c.Value().(dgo.Array)
	assert.NotSame(t, ac.Get(0), ac.Get(1))
	assert.Equal(t, ac.Get(0), ac.Get(1))
}

func TestEncode_bigfloat_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	assert.Panic(t, func() {
		streamer.New(nil, o).Stream(
			vf.New(typ.BigFloat, vf.Value(`123456789012345678901234567890123.3412`)),
			streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_bigfloat_small_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	b := bytes.Buffer{}
	streamer.New(nil, o).Stream(
		vf.Value(big.NewFloat(3.14)),
		streamer.JSON(&b))
	assert.Equal(t, `3.14`, b.String())
}

func TestEncode_bigint_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	assert.Panic(t, func() {
		streamer.New(nil, o).Stream(
			vf.New(typ.Integer, vf.Value(`12345678901234567890123`)),
			streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_binary_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	assert.Panic(t, func() {
		streamer.New(nil, o).Stream(vf.BinaryFromString(`AQID`), streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_time_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	assert.Panic(t, func() {
		streamer.New(nil, o).Stream(vf.Time(time.Now()), streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_sensitive_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	assert.Panic(t, func() {
		streamer.New(nil, o).Stream(vf.Sensitive(`secret`), streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_not_string_key_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	assert.Panic(t, func() {
		streamer.New(nil, o).Stream(vf.Map(1, "one"), streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_dedup_string_keys(t *testing.T) {
	a := vf.String(`the key`)
	b := vf.String(`the key`)
	assert.NotSame(t, a, b)
	c := streamer.NewCollector()
	o := streamer.DefaultOptions()
	o.DedupLevel = streamer.MaxDedup
	streamer.New(nil, o).Stream(vf.Values(vf.Map(a, true), vf.Map(b, true)), c)
	ac := c.Value().(dgo.Array)
	var aes, bes []dgo.MapEntry
	ac.Get(0).(dgo.Map).EachEntry(func(e dgo.MapEntry) { aes = append(aes, e) })
	ac.Get(1).(dgo.Map).EachEntry(func(e dgo.MapEntry) { bes = append(bes, e) })
	assert.Same(t, aes[0].Key(), bes[0].Key())
}

func TestEncode_existingRef(t *testing.T) {
	v := vf.Map(`x`, vf.Map(`z`, vf.Values(`a`, `b`)), `y`, vf.Map(streamer.DgoDialect().RefKey(), 6))
	vc := streamer.DataCollector()
	streamer.New(nil, nil).Stream(v, vc)
	assert.Equal(t, vf.Map(`x`, vf.Map(`z`, vf.Values("a", `b`)), `y`, `b`), vc.Value())
}

type testNamed struct {
	A string
	B int64
}

func (v *testNamed) String() string {
	return v.Type().(dgo.NamedType).ValueString(v)
}

func (v *testNamed) Type() dgo.Type {
	return tf.ExactNamed(tf.Named(`testNamed`), v)
}

func (v *testNamed) Equals(other interface{}) bool {
	if ov, ok := other.(*testNamed); ok {
		return *v == *ov
	}
	return false
}

func (v *testNamed) HashCode() dgo.Hash {
	return dgo.Hash(reflect.ValueOf(v).Pointer())
}

func TestEncode_namedUsingMap(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		sm := vf.MutableMap(&testNamed{}).(dgo.Struct)
		sm.PutAll(arg.(dgo.Map))
		return sm.GoStruct().(dgo.Value)
	}, func(value dgo.Value) dgo.Value {
		return vf.Map(value)
	}, reflect.TypeOf(&testNamed{}), nil, nil)

	arg := vf.Map(`A`, `hello`, `B`, 32)
	s := streamer.New(nil, nil)

	v := tp.New(arg)
	c := streamer.NewCollector()
	s.Stream(v, c)
	assert.Equal(t, arg.With(`__type`, `testNamed`), c.Value())

	c = streamer.DataDecoder(nil, nil)
	s.Stream(v, c)
	assert.Equal(t, v, c.Value())
}

func TestEncode_namedUsingValue(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		a := arg.(dgo.Array)
		return &testNamed{A: a.Get(0).(dgo.String).GoString(), B: a.Get(1).(dgo.Integer).GoInt()}
	}, func(value dgo.Value) dgo.Value {
		v := value.(*testNamed)
		return vf.Values(v.A, v.B)
	}, reflect.TypeOf(&testNamed{}), nil, nil)

	arg := vf.Values(`hello`, 32)
	s := streamer.New(nil, nil)

	v := tp.New(arg)
	c := streamer.NewCollector()
	s.Stream(v, c)
	assert.Equal(t, vf.Map(`__type`, `testNamed`, `__value`, arg), c.Value())

	c = streamer.DataDecoder(nil, nil)
	s.Stream(v, c)
	assert.Equal(t, v, c.Value())

	c = streamer.DataDecoder(nil, nil)
	s.Stream(vf.Map(`__type`, `testNamed`), c)
	assert.Same(t, tp, c.Value())
}

func TestEncode_named_notRich(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		sm := vf.MutableMap(&testNamed{}).(dgo.Struct)
		sm.PutAll(arg.(dgo.Map))
		return sm.GoStruct().(dgo.Value)
	}, func(value dgo.Value) dgo.Value {
		return vf.Map(value)
	}, reflect.TypeOf(&testNamed{}), nil, nil)

	arg := vf.Map(`A`, `hello`, `B`, 32)
	o := streamer.DefaultOptions()
	o.RichData = false
	s := streamer.New(nil, o)

	v := tp.New(arg)
	c := streamer.NewCollector()
	assert.Panic(t, func() { s.Stream(v, c) }, `unable to serialize`)
}
