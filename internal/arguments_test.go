package internal_test

import (
	"testing"

	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestArguments_Arg(t *testing.T) {
	args := vf.Arguments(`first`, 2)
	assert.Equal(t, `first`, args.Arg(`myfunc`, 0, typ.String))
	assert.Panic(t, func() { args.Arg(`myfunc`, 1, typ.String) }, `illegal argument 2 for myfunc`)
}

func TestArgumentsFromArray(t *testing.T) {
	assert.Equal(t, vf.Arguments(`first`, 2), vf.ArgumentsFromArray(vf.Values(`first`, 2)))
}
