package vf_test

import (
	"fmt"
	"reflect"

	"github.com/lyraproj/dgo/vf"
)

func ExampleValue() {
	v := vf.Value([]string{`a`, `b`})
	fmt.Println(v.Equals([]string{`a`, `b`}))
	// Output: true
}

func ExampleValueFromReflected() {
	v := vf.ValueFromReflected(reflect.ValueOf([]string{`a`, `b`}))
	fmt.Println(v.Equals([]string{`a`, `b`}))
	// Output: true
}

func ExampleFromValue_map_string_string() {
	v := vf.Map(
		`first`, `one`,
		`second`, `two`,
		`third`, `three`)

	var gv map[string]string

	vf.FromValue(v, &gv)
	fmt.Println(gv) // As of go 1.12 maps are printed in key-sorted order
	// Output: map[first:one second:two third:three]
}

func ExampleFromValue_map_string_interface() {
	v := vf.Map(
		`first`, 1,
		`second`, 2.1,
		`third`, `three`)

	var gim map[string]interface{}
	vf.FromValue(v, &gim)
	fmt.Println(gim)
	// Output: map[first:1 second:2.1 third:three]
}

func ExampleFromValue_slice() {
	v := vf.Strings(`first`, `second`, `third`)

	var gs []string
	vf.FromValue(v, &gs)
	fmt.Println(gs)
	// Output: [first second third]
}

func ExampleFromValue_slice_ints() {
	v := vf.Integers(10, 20, 30)

	var gs []int16
	vf.FromValue(v, &gs)
	fmt.Println(gs)

	var gu []uint64
	vf.FromValue(v, &gu)
	fmt.Println(gu)
	// Output:
	// [10 20 30]
	// [10 20 30]
}
