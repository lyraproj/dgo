package streamer_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/streamer"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestDgoDialect_names(t *testing.T) {
	for k, v := range map[string]func(streamer.Dialect) dgo.String{
		`alias`:     streamer.Dialect.AliasTypeName,
		`bigfloat`:  streamer.Dialect.BigFloatTypeName,
		`bigint`:    streamer.Dialect.BigIntTypeName,
		`binary`:    streamer.Dialect.BinaryTypeName,
		`map`:       streamer.Dialect.MapTypeName,
		`sensitive`: streamer.Dialect.SensitiveTypeName,
		`time`:      streamer.Dialect.TimeTypeName,
		`__ref`:     streamer.Dialect.RefKey,
		`__type`:    streamer.Dialect.TypeKey,
		`__value`:   streamer.Dialect.ValueKey,
	} {
		assert.Equal(t, k, v(streamer.DgoDialect()))
	}
}

func TestDgoDialect_ParseType(t *testing.T) {
	assert.Equal(t, tf.Array(typ.String, 3, 8), streamer.DgoDialect().ParseType(nil, vf.String(`[3,8]string`)))
}
