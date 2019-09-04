package internal_test

import (
	"bytes"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"

	"github.com/lyraproj/dgo/vf"
)

func TestBinary(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.Binary(bs)
	require.Equal(t, `AQID`, v.String())

	// Test immutability
	bs[1] = 4
	require.Equal(t, `AQID`, v.String())
	require.False(t, bytes.Equal(bs, v.GoBytes()))
}

func TestBinaryString(t *testing.T) {
	bs := []byte{1, 2, 3}
	v := vf.BinaryFromString(`AQID`)
	require.True(t, bytes.Equal(bs, v.GoBytes()))
}
