package util_test

import (
	"testing"

	require "github.com/tada/dgo/dgo_test"

	"github.com/tada/dgo/util"
)

func TestReader_Next(t *testing.T) {
	v := util.NewStringReader(`röd`)
	require.Equal(t, 0, v.Pos())
	require.Equal(t, 'r', v.Next())
	require.Equal(t, 1, v.Pos())
	require.Equal(t, 'ö', v.Next())
	require.Equal(t, 3, v.Pos())
	require.Equal(t, 'd', v.Next())
	require.Equal(t, 4, v.Pos())
	require.Equal(t, 0, v.Next())
	require.Equal(t, 5, v.Pos())
	require.Equal(t, 0, v.Next())
	require.Equal(t, 5, v.Pos())
}

func TestReader_Rewind(t *testing.T) {
	v := util.NewStringReader(`röd`)
	require.Equal(t, 'r', v.Next())
	require.Equal(t, 'ö', v.Next())
	require.Equal(t, 'd', v.Next())
	v.Rewind()
	require.Equal(t, 0, v.Pos())
	require.Equal(t, 'r', v.Next())
}

func TestReader_Peek(t *testing.T) {
	v := util.NewStringReader(`röd`)
	require.Equal(t, 'r', v.Next())
	require.Equal(t, 'ö', v.Peek())
	require.Equal(t, 'ö', v.Next())
	require.Equal(t, 'd', v.Next())
	require.Equal(t, 0, v.Peek())

	v = util.NewStringReader(string([]byte{'r', 0x82, 0xff}))
	require.Equal(t, 'r', v.Next())
	require.Panic(t, func() { v.Peek() }, `unicode error`)
}

func TestReader_Peek2(t *testing.T) {
	v := util.NewStringReader(`röd`)
	require.Equal(t, 'ö', v.Peek2())
	require.Equal(t, 'r', v.Next())
	require.Equal(t, 'ö', v.Next())
	require.Equal(t, 0, v.Peek2())
	require.Equal(t, 'd', v.Next())
	require.Equal(t, 0, v.Peek2())

	v = util.NewStringReader(string([]byte{0xc3, 0xb6, 0x82, 0xff}))
	require.Panic(t, func() { v.Peek2() }, `unicode error`)

	v = util.NewStringReader(string([]byte{0x82, 0xff, 0xc3, 0xb6}))
	require.Panic(t, func() { v.Peek2() }, `unicode error`)
}
