package util_test

import (
	"testing"

	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/util"
)

func TestReader_Next(t *testing.T) {
	v := util.NewStringReader(`röd`)
	assert.Equal(t, 0, v.Pos())
	assert.Equal(t, 'r', v.Next())
	assert.Equal(t, 1, v.Pos())
	assert.Equal(t, 'ö', v.Next())
	assert.Equal(t, 3, v.Pos())
	assert.Equal(t, 'd', v.Next())
	assert.Equal(t, 4, v.Pos())
	assert.Equal(t, 0, v.Next())
	assert.Equal(t, 5, v.Pos())
	assert.Equal(t, 0, v.Next())
	assert.Equal(t, 5, v.Pos())
}

func TestReader_Rewind(t *testing.T) {
	v := util.NewStringReader(`röd`)
	assert.Equal(t, 'r', v.Next())
	assert.Equal(t, 'ö', v.Next())
	assert.Equal(t, 'd', v.Next())
	v.Rewind()
	assert.Equal(t, 0, v.Pos())
	assert.Equal(t, 'r', v.Next())
}

func TestReader_Peek(t *testing.T) {
	v := util.NewStringReader(`röd`)
	assert.Equal(t, 'r', v.Next())
	assert.Equal(t, 'ö', v.Peek())
	assert.Equal(t, 'ö', v.Next())
	assert.Equal(t, 'd', v.Next())
	assert.Equal(t, 0, v.Peek())

	v = util.NewStringReader(string([]byte{'r', 0x82, 0xff}))
	assert.Equal(t, 'r', v.Next())
	assert.Panic(t, func() { v.Peek() }, `unicode error`)
}

func TestReader_Peek2(t *testing.T) {
	v := util.NewStringReader(`röd`)
	assert.Equal(t, 'ö', v.Peek2())
	assert.Equal(t, 'r', v.Next())
	assert.Equal(t, 'ö', v.Next())
	assert.Equal(t, 0, v.Peek2())
	assert.Equal(t, 'd', v.Next())
	assert.Equal(t, 0, v.Peek2())

	v = util.NewStringReader(string([]byte{0xc3, 0xb6, 0x82, 0xff}))
	assert.Panic(t, func() { v.Peek2() }, `unicode error`)

	v = util.NewStringReader(string([]byte{0x82, 0xff, 0xc3, 0xb6}))
	assert.Panic(t, func() { v.Peek2() }, `unicode error`)
}
