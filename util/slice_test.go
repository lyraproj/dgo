package util_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func TestSliceCopy(t *testing.T) {
	vs := []dgo.Value{vf.String(`a`), vf.Integer(32)}
	vc := util.SliceCopy(vs)
	require.Equal(t, vs[0], vc[0])
	vs[0] = vf.String(`b`)
	require.NotEqual(t, vs[0], vc[0])
}
