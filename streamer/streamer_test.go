package streamer_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/vf"
)

func TestNew(t *testing.T) {
	s := streamer.New(nil, nil)
	c := streamer.NewCollector()
	m := vf.Map(`a`, 1, `b`, 2)
	s.Stream(m, c)
	require.Equal(t, m, c.Value())
}

func TestEncode_dedup(t *testing.T) {
	v := vf.Strings(`a`, `b`)
	a := vf.Values(v, v)
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(a, c)
	ac := c.Value().(dgo.Array)
	require.Same(t, ac.Get(0), ac.Get(1))
}

func TestEncode_demoteDedup(t *testing.T) {
	a := vf.Values(vf.Map(`the key`, `a`), vf.Map(`the key`, `b`))
	b := bytes.Buffer{}
	o := streamer.DefaultOptions()
	o.DedupLevel = streamer.MaxDedup
	streamer.New(nil, o).Stream(a, streamer.JSON(&b))
	require.Equal(t, `[{"the key":"a"},{"the key":"b"}]`, b.String())
}

func TestEncode_sensitive(t *testing.T) {
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(vf.Sensitive(`secret`), c)
	require.Equal(t, vf.Map(`__type`, `sensitive`, `__value`, `secret`), c.Value())
}

func TestEncode_unknown(t *testing.T) {
	c := streamer.NewCollector()
	require.Panic(t, func() { streamer.New(nil, nil).Stream(vf.Value(struct{ A string }{A: `hello`}), c) }, `unable to serialize value`)
}

func TestEncode_type(t *testing.T) {
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(tf.String(1), c)
	require.Equal(t, vf.Map(`__type`, `string[1]`), c.Value())
}

func TestEncode_binary(t *testing.T) {
	v := vf.BinaryFromString(`AQID`)
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(v, c)
	require.Same(t, v, c.Value())
}

func TestEncode_time(t *testing.T) {
	v := vf.Time(time.Now())
	c := streamer.NewCollector()
	streamer.New(nil, nil).Stream(v, c)
	require.Same(t, v, c.Value())
}

func TestEncode_alias(t *testing.T) {
	am := tf.NewAliasMap()
	tp := tf.ParseFile(am, ``, `ne=string[1]`)
	c := streamer.NewCollector()
	streamer.New(am, nil).Stream(tp, c)
	require.Equal(t, vf.Map(`__type`, `alias`, `__value`, vf.Strings(`ne`, `string[1]`)), c.Value())
}

func TestEncode_noDedup(t *testing.T) {
	v := vf.Strings(`a`, `b`)
	a := vf.Values(v, v)
	c := streamer.NewCollector()
	o := streamer.DefaultOptions()
	o.DedupLevel = streamer.NoDedup
	streamer.New(nil, o).Stream(a, c)
	ac := c.Value().(dgo.Array)
	require.NotSame(t, ac.Get(0), ac.Get(1))
	require.Equal(t, ac.Get(0), ac.Get(1))
}

func TestEncode_binary_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	require.Panic(t, func() {
		streamer.New(nil, o).Stream(vf.BinaryFromString(`AQID`), streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_time_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	require.Panic(t, func() {
		streamer.New(nil, o).Stream(vf.Time(time.Now()), streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_sensitive_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	require.Panic(t, func() {
		streamer.New(nil, o).Stream(vf.Sensitive(`secret`), streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_not_string_key_not_rich(t *testing.T) {
	o := streamer.DefaultOptions()
	o.RichData = false
	require.Panic(t, func() {
		streamer.New(nil, o).Stream(vf.Map(1, "one"), streamer.JSON(&bytes.Buffer{}))
	}, `unable to serialize`)
}

func TestEncode_dedup_string_keys(t *testing.T) {
	a := vf.String(`the key`)
	b := vf.String(`the key`)
	require.NotSame(t, a, b)
	c := streamer.NewCollector()
	o := streamer.DefaultOptions()
	o.DedupLevel = streamer.MaxDedup
	streamer.New(nil, o).Stream(vf.Values(vf.Map(a, true), vf.Map(b, true)), c)
	ac := c.Value().(dgo.Array)
	var aes, bes []dgo.MapEntry
	ac.Get(0).(dgo.Map).EachEntry(func(e dgo.MapEntry) { aes = append(aes, e) })
	ac.Get(1).(dgo.Map).EachEntry(func(e dgo.MapEntry) { bes = append(bes, e) })
	require.Same(t, aes[0].Key(), bes[0].Key())
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

func (v *testNamed) HashCode() int {
	return int(reflect.ValueOf(v).Pointer())
}

func TestEncode_namedUsingMap(t *testing.T) {
	defer tf.RemoveNamed(`testNamed`)
	tp := tf.NewNamed(`testNamed`, func(arg dgo.Value) dgo.Value {
		sm := vf.Map(&testNamed{}).(dgo.Struct)
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
	require.Equal(t, arg.With(`__type`, `testNamed`), c.Value())

	c = streamer.DataDecoder(nil, nil)
	s.Stream(v, c)
	require.Equal(t, v, c.Value())
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
	require.Equal(t, vf.Map(`__type`, `testNamed`, `__value`, arg), c.Value())

	c = streamer.DataDecoder(nil, nil)
	s.Stream(v, c)
	require.Equal(t, v, c.Value())

	c = streamer.DataDecoder(nil, nil)
	s.Stream(vf.Map(`__type`, `testNamed`), c)
	require.Same(t, tp, c.Value())
}
