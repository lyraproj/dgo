package internal_test

import (
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"

	"github.com/lyraproj/dgo/newtype"
)

func TestAliasMap_Get(t *testing.T) {
	am := newtype.NewAliasMap()
	newtype.ParseFile(am, `example.dgo`, `a=string[1]`)
	require.Equal(t, `a`, am.GetName(newtype.String(1)))
	require.Nil(t, am.GetName(typ.String))
}
