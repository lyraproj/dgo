package vf_test

import (
	"fmt"

	"github.com/lyraproj/dgo/vf"
)

func ExampleBinary_frozen() {
	bs := []byte{1, 2, 3}
	b := vf.Binary(bs, true)
	fmt.Println(b)
	bs[0] = 3
	fmt.Println(b)
	// Output:
	// AQID
	// AQID
}

func ExampleBinary_mutable() {
	bs := []byte{1, 2, 3}
	b := vf.Binary(bs, false)
	fmt.Println(b)
	bs[0] = 3
	fmt.Println(b)
	// Output:
	// AQID
	// AwID
}

func ExampleBinaryFromString() {
	b := vf.BinaryFromString(`AQID`)
	fmt.Println(b.GoBytes())
	// Output: [1 2 3]
}