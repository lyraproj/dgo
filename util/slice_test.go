package util_test

import (
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/util"
	"github.com/tada/dgo/vf"
)

func TestSliceCopy(t *testing.T) {
	vs := []dgo.Value{vf.String(`a`), vf.Integer(32)}
	vc := util.SliceCopy(vs)
	assert.Equal(t, vs[0], vc[0])
	vs[0] = vf.String(`b`)
	assert.NotEqual(t, vs[0], vc[0])
}
