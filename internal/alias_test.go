package internal_test

import (
	"sync"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/parser"
	"github.com/lyraproj/dgo/stringer"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestAlias_Freeze(t *testing.T) {
	alias := parser.NewAlias(vf.String(`hello`)).(dgo.Freezable)
	assert.False(t, alias.Frozen())
	assert.Same(t, alias, alias.ThawedCopy())
	assert.Panic(t, func() {
		alias.Freeze()
	}, `attempt to freeze unresolved alias`)
	assert.Panic(t, func() {
		alias.FrozenCopy()
	}, `attempt to freeze unresolved alias`)
}

func TestAliasMap_Get(t *testing.T) {
	am := tf.BuiltInAliases().Collect(func(a dgo.AliasAdder) {
		tf.ParseFile(a, `example.dgo`, `a=string[1]`)
	})
	assert.Equal(t, `a`, am.GetName(tf.String(1)))
	assert.Nil(t, am.GetName(typ.String))
}

func TestAliasMap_exactStromg(t *testing.T) {
	am := tf.BuiltInAliases().Collect(func(a dgo.AliasAdder) {
		tf.ParseFile(a, `example.dgo`, `a="the exact string"`)
	})
	assert.Equal(t, `a`, stringer.TypeStringWithAliasMap(vf.String("the exact string").Type(), am))
}

func TestDefaultAliasMap_Get(t *testing.T) {
	rd := vf.String(`richData`)
	assert.Equal(t, tf.DefaultAliases().GetType(rd), tf.BuiltInAliases().GetType(rd))
}

func TestAddDefaultAliases(t *testing.T) {
	rd := vf.String(`richData`)
	assert.Equal(t, tf.DefaultAliases().GetType(rd), tf.BuiltInAliases().GetType(rd))
}

func TestAddAliases(t *testing.T) {
	lock := sync.Mutex{}
	bi := tf.BuiltInAliases()
	aliases := bi
	tf.AddAliases(&aliases, &lock, func(aa dgo.AliasAdder) {
	})
	assert.Same(t, bi, aliases)

	tf.AddAliases(&aliases, &lock, func(aa dgo.AliasAdder) {
		aa.Add(tf.String(10, 12), vf.String(`pnr`))
	})
	assert.NotSame(t, bi, aliases)
	assert.Equal(t, aliases.GetType(vf.String(`pnr`)), tf.String(10, 12))
	assert.Nil(t, bi.GetType(vf.String(`pnr`)))
}

func TestNewCall(t *testing.T) {
	nc := internal.NewCall(typ.Binary, vf.Arguments(vf.Values(0x01, 0x02, 0x03)))
	tf.AddDefaultAliases(func(aa dgo.AliasAdder) {
		assert.Equal(t, vf.Binary([]byte{0x01, 0x02, 0x03}, true), aa.Replace(nc))
	})
}
