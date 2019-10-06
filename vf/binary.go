package vf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Binary creates a new Binary based on the given slice. If frozen is true, the
// binary will be immutable and contain a copy of the slice, otherwise the slice
// is simply wrapped and modifications to its elements will also modify the binary.
func Binary(bs []byte, frozen bool) dgo.Binary {
	return internal.Binary(bs, frozen)
}

// BinaryFromString creates a new Binary from the base64 encoded string
func BinaryFromString(base64 string) dgo.Binary {
	return internal.BinaryFromString(base64)
}
