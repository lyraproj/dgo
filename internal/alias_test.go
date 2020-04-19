package internal_test

import (
	"sync"
	"testing"

	"github.com/tada/dgo/internal"

	"github.com/tada/dgo/parser"

	"github.com/tada/dgo/dgo"

	"github.com/tada/dgo/vf"

	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/typ"

	"github.com/tada/dgo/tf"
)

func TestAlias_Freeze(t *testing.T) {
	alias := parser.NewAlias(vf.String(`hello`)).(dgo.Freezable)
	require.False(t, alias.Frozen())
	require.Same(t, alias, alias.ThawedCopy())
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

func TestAddAliases(t *testing.T) {
	lock := sync.Mutex{}
	bi := tf.BuiltInAliases()
	aliases := bi
	tf.AddAliases(&aliases, &lock, func(aa dgo.AliasAdder) {
	})
	require.Same(t, bi, aliases)

	tf.AddAliases(&aliases, &lock, func(aa dgo.AliasAdder) {
		aa.Add(tf.String(10, 12), vf.String(`pnr`))
	})
	require.NotSame(t, bi, aliases)
	require.Equal(t, aliases.GetType(vf.String(`pnr`)), tf.String(10, 12))
	require.Nil(t, bi.GetType(vf.String(`pnr`)))
}

func TestNewCall(t *testing.T) {
	nc := internal.NewCall(typ.Binary, vf.Arguments(vf.Values(0x01, 0x02, 0x03)))
	tf.AddDefaultAliases(func(aa dgo.AliasAdder) {
		require.Equal(t, vf.Binary([]byte{0x01, 0x02, 0x03}, true), aa.Replace(nc))
	})
}
