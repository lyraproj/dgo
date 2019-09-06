package internal_test

import (
	"regexp"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestMap_StructType(t *testing.T) {
	tp := newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, true),
		newtype.StructEntry(`b`, typ.String, false))

	require.False(t, tp.Additional())

	m := vf.Map(map[string]interface{}{`a`: 3, `b`: `yes`})
	require.Assignable(t, tp, tp)
	require.Assignable(t, tp, m.Type())
	require.Instance(t, tp, m)

	m = vf.Map(map[string]interface{}{`a`: 3})
	require.Assignable(t, tp, m.Type())
	require.Instance(t, tp, m)

	m = vf.Map(map[string]interface{}{`b`: `yes`})
	require.NotAssignable(t, tp, m.Type())
	require.NotInstance(t, tp, m)

	m = vf.Map(map[string]interface{}{`a`: 3, `b`: 4})
	require.NotAssignable(t, tp, m.Type())
	require.NotInstance(t, tp, m)

	require.NotInstance(t, tp, vf.Values(`a`, `b`))

	require.Instance(t, tp.Type(), tp)

	tps := newtype.Struct(false,
		newtype.StructEntry(`a`, newtype.IntegerRange(0, 10), true),
		newtype.StructEntry(`b`, newtype.String(20), false))
	require.Assignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, true),
		newtype.StructEntry(`b`, typ.String, true))
	require.Assignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, false),
		newtype.StructEntry(`b`, typ.String, false))
	require.NotAssignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`a`, typ.Integer, true))
	require.Assignable(t, tp, tps)

	tps = newtype.Struct(true,
		newtype.StructEntry(`a`, typ.Integer, true))
	require.NotAssignable(t, tp, tps)

	tps = newtype.Struct(false,
		newtype.StructEntry(`b`, typ.String, false))
	require.NotAssignable(t, tp, tps)

	require.NotEqual(t, 0, tp.HashCode())
	require.NotEqual(t, tp.HashCode(), tps.HashCode())

	require.Panic(t, func() {
		newtype.Struct(false,
			internal.StructEntry2(newtype.Pattern(regexp.MustCompile(`a*`)), typ.Integer, true))
	}, `non exact key types`)
}
