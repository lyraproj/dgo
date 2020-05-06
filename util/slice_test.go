package util_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func TestSliceCopy(t *testing.T) {
	vs := []dgo.Value{vf.String(`a`), vf.Int64(32)}
	vc := util.SliceCopy(vs)
	assert.Equal(t, vs[0], vc[0])
	vs[0] = vf.String(`b`)
	assert.NotEqual(t, vs[0], vc[0])
}
