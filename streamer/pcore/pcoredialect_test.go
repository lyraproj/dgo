package pcore_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/streamer/pcore"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestPcoreDialect_names(t *testing.T) {
	for k, v := range map[string]func(streamer.Dialect) dgo.String{
		`Alias`:     streamer.Dialect.AliasTypeName,
		`Binary`:    streamer.Dialect.BinaryTypeName,
		`Hash`:      streamer.Dialect.MapTypeName,
		`Sensitive`: streamer.Dialect.SensitiveTypeName,
		`Timestamp`: streamer.Dialect.TimeTypeName,
		`__pref`:    streamer.Dialect.RefKey,
		`__ptype`:   streamer.Dialect.TypeKey,
		`__pvalue`:  streamer.Dialect.ValueKey,
	} {
		require.Equal(t, k, v(pcore.Dialect()))
	}
}

func TestPcoreDialect_ParseType(t *testing.T) {
	require.Equal(t, tf.Array(typ.String, 3, 8), pcore.Dialect().ParseType(nil, vf.String(`Array[String,3,8]`)))
}
