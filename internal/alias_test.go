package internal_test

import (
	"sync"
	"testing"

	"github.com/tada/dgo/parser"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
	"github.com/tada/dgo/stringer"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestAlias_Freeze(t *testing.T) {
	alias := parser.NewAlias(vf.String(`hello`)).(dgo.Mutability)
	assert.False(t, alias.Frozen())
	assert.Same(t, alias, alias.ThawedCopy())
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

	mt := vf.MutableValues(`a`, `b`).Type()
	tf.AddAliases(&aliases, &lock, func(aa dgo.AliasAdder) {
		aa.Add(tf.String(10, 12), vf.String(`pnr`))
		aa.Add(mt, vf.String(`mt`))
	})
	assert.NotSame(t, bi, aliases)
	assert.Equal(t, aliases.GetType(vf.String(`pnr`)), tf.String(10, 12))
	assert.False(t, mt.(dgo.Mutability).Frozen())
	assert.True(t, aliases.GetType(vf.String(`mt`)).(dgo.Mutability).Frozen())
	assert.Nil(t, bi.GetType(vf.String(`pnr`)))
}

func TestNewCall(t *testing.T) {
	nc := internal.NewCall(typ.Binary, vf.Arguments(vf.Values(0x01, 0x02, 0x03)))
	tf.AddDefaultAliases(func(aa dgo.AliasAdder) {
		assert.Equal(t, vf.Binary([]byte{0x01, 0x02, 0x03}, true), aa.Replace(nc))
	})
}
