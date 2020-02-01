package internal_test

import (
	"testing"

	"github.com/lyraproj/dgo/parser"

	"github.com/lyraproj/dgo/dgo"

	"github.com/lyraproj/dgo/vf"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"

	"github.com/lyraproj/dgo/tf"
)

func TestAlias_Freeze(t *testing.T) {
	alias := parser.NewAlias(vf.String(`hello`)).(dgo.Freezable)
	require.False(t, alias.Frozen())
	require.Panic(t, func() {
		alias.Freeze()
	}, `attempt to freeze unresolved alias`)
	require.Panic(t, func() {
		alias.FrozenCopy()
	}, `attempt to freeze unresolved alias`)
}

func TestAliasMap_Get(t *testing.T) {
	am := tf.BuiltInAliases().Collect(func(a dgo.AliasAdder) {
		tf.ParseFile(a, `example.dgo`, `a=string[1]`)
	})
	require.Equal(t, `a`, am.GetName(tf.String(1)))
	require.Nil(t, am.GetName(typ.String))
}

func TestDefaultAliasMap_Get(t *testing.T) {
	rd := vf.String(`richData`)
	require.Equal(t, tf.DefaultAliases().GetType(rd), tf.BuiltInAliases().GetType(rd))
}

func TestAddDefaultAliases(t *testing.T) {
	rd := vf.String(`richData`)
	require.Equal(t, tf.DefaultAliases().GetType(rd), tf.BuiltInAliases().GetType(rd))
}
