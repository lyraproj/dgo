package pcore

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestDgoDialect_names(t *testing.T) {
	for k, v := range map[string]func(pcoreDialect) dgo.String{
		`Alias`:     pcoreDialect.AliasTypeName,
		`Float`:     pcoreDialect.BigFloatTypeName,
		`Integer`:   pcoreDialect.BigIntTypeName,
		`Binary`:    pcoreDialect.BinaryTypeName,
		`Hash`:      pcoreDialect.MapTypeName,
		`Sensitive`: pcoreDialect.SensitiveTypeName,
		`Timespan`:  pcoreDialect.DurationTypeName,
		`Timestamp`: pcoreDialect.TimeTypeName,
		`__pref`:    pcoreDialect.RefKey,
		`__ptype`:   pcoreDialect.TypeKey,
		`__pvalue`:  pcoreDialect.ValueKey,
	} {
		assert.Equal(t, k, v(pcoreDialectSingleton))
	}
}

func TestPcoreDialect_ParseType(t *testing.T) {
	require.Equal(t, tf.Array(typ.String, 3, 8), Dialect().ParseType(nil, vf.String(`Array[String,3,8]`)))
}
