package streamer_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/tf"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestDataDecoder(t *testing.T) {
	a := vf.Values(
		typ.String,
		vf.Binary([]byte{1, 2, 3}, false),
		vf.Sensitive(`secret`),
		vf.TimeFromString(`2019-10-06T16:15:00.123-01:00`))

	// Transform rich data to plain data
	c := streamer.DataCollector()
	streamer.New(nil, nil).Stream(a, c)

	// Transform back to plain data
	d := streamer.DataDecoder(nil, nil)
	streamer.New(nil, nil).Stream(c.Value(), d)

	require.Equal(t, a, d.Value())
}

func TestDataDecoder_bad_time(t *testing.T) {
	dv := vf.Map(`__type`, `time`, `__value`, `2019-13-06T16:15:00.123-01:00`)

	// Transform back to plain data
	d := streamer.DataDecoder(nil, nil)
	require.Panic(t, func() { streamer.New(nil, nil).Stream(dv, d) }, `month out of range`)
}

func TestDataDecoder_bad_type(t *testing.T) {
	dv := vf.Map(`__type`, `unrecognized`, `__value`, `foo`)

	// Transform back to plain data
	d := streamer.DataDecoder(nil, nil)
	require.Panic(t, func() { streamer.New(nil, nil).Stream(dv, d) }, `unable to decode __type: unrecognized`)
}

func TestDataDecoder_selfref(t *testing.T) {
	m := vf.MutableMap()
	m.Put(typ.Integer, vf.Values(`a`, `b`))
	m.Put(typ.String, vf.MutableValues(vf.Sensitive(m)))

	// Transform rich data to plain data
	c := streamer.DataCollector()
	streamer.New(nil, nil).Stream(m, c)

	// Transform back to plain data
	d := streamer.DataDecoder(nil, nil)
	streamer.New(nil, nil).Stream(c.Value(), d)

	dv := d.Value()
	require.Equal(t, m, dv)

	// Check that self reference is restored
	require.Same(t, dv, dv.(dgo.Map).Get(typ.String).(dgo.Array).Get(0).(dgo.Sensitive).Unwrap())
}

func TestDataDecoder_alias(t *testing.T) {
	aliasMap := tf.NewAliasMap()
	tf.ParseFile(aliasMap, ``, `ne=string[1]`)
	a := vf.Values(tf.String(1))

	// Transform rich data to plain data
	c := streamer.DataCollector()
	streamer.New(aliasMap, nil).Stream(a, c)

	// Transform back to plain data
	aliasMap = tf.NewAliasMap()
	d := streamer.DataDecoder(aliasMap, nil)
	streamer.New(nil, nil).Stream(c.Value(), d)

	require.Equal(t, `ne`, aliasMap.GetName(tf.String(1)))
}
