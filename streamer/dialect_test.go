package streamer_test

import (
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/streamer"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
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
		assert.Equal(t, k, v(streamer.DgoDialect()))
	}
}

func TestDgoDialect_ParseType(t *testing.T) {
	assert.Equal(t, tf.Array(typ.String, 3, 8), streamer.DgoDialect().ParseType(nil, vf.String(`[3,8]string`)))
}
