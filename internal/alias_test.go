package internal_test

import (
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"

	"github.com/lyraproj/dgo/tf"
)

func TestAliasMap_Get(t *testing.T) {
	am := tf.NewAliasMap()
	tf.ParseFile(am, `example.dgo`, `a=string[1]`)
	require.Equal(t, `a`, am.GetName(tf.String(1)))
	require.Nil(t, am.GetName(typ.String))
}
