package streamer_test

import (
	"testing"

	"github.com/lyraproj/dgo/vf"

	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"

	"github.com/lyraproj/dgo/streamer"

	require "github.com/lyraproj/dgo/dgo_test"

	"github.com/lyraproj/dgo/dgo"
)

func TestDgoDialect_names(t *testing.T) {
	for k, v := range map[string]func(streamer.Dialect) dgo.String{
		`alias`:     streamer.Dialect.AliasTypeName,
		`binary`:    streamer.Dialect.BinaryTypeName,
		`map`:       streamer.Dialect.MapTypeName,
		`sensitive`: streamer.Dialect.SensitiveTypeName,
		`time`:      streamer.Dialect.TimeTypeName,
		`__ref`:     streamer.Dialect.RefKey,
		`__type`:    streamer.Dialect.TypeKey,
		`__value`:   streamer.Dialect.ValueKey,
	} {
		require.Equal(t, k, v(streamer.DgoDialect()))
	}
}

func TestDgoDialect_ParseType(t *testing.T) {
	require.Equal(t, newtype.Array(typ.String, 3, 8), streamer.DgoDialect().ParseType(vf.String(`[3,8]string`)))
}
