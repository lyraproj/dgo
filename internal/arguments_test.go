package internal_test

import (
	"testing"

	require "github.com/tada/dgo/dgo_test"

	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestArguments_Arg(t *testing.T) {
	args := vf.Arguments(`first`, 2)
	require.Equal(t, `first`, args.Arg(`myfunc`, 0, typ.String))
	require.Panic(t, func() { args.Arg(`myfunc`, 1, typ.String) }, `illegal argument 2 for myfunc`)
}

func TestArgumentsFromArray(t *testing.T) {
	require.Equal(t, vf.Arguments(`first`, 2), vf.ArgumentsFromArray(vf.Values(`first`, 2)))
}
