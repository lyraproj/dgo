package internal_test

import (
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"

	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestArguments_Arg(t *testing.T) {
	args := vf.Arguments(`first`, 2)
	require.Equal(t, `first`, args.Arg(`myfunc`, 0, typ.String))
	require.Panic(t, func() { args.Arg(`myfunc`, 1, typ.String) }, `illegal argument 2 for myfunc`)
}
