package vf_test

import (
	"fmt"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/typ"

	"github.com/lyraproj/dgo/vf"
)

func ExampleValues() {
	fmt.Printf("%#v\n", vf.Values(`hello`, 1, 2.3, true, []byte{1, 2, 3}))
	// Output: []any{"hello", 1, 2.3, true, []byte{0x1, 0x2, 0x3}}
}

func ExampleStrings() {
	fmt.Printf("%#v\n", vf.Strings(`one`, `two`))
	// Output: []string{"one", "two"}
}

func ExampleMutableValues() {
	a := vf.MutableValues()
	a.Add(32)
	fmt.Printf("%#v\n", a)
	// Output: []int{32}
}

func ExampleArguments() {
	args := vf.Arguments(`first`, 2)
	a0 := args.Arg(`myFunc`, 0, typ.String).(dgo.String)
	fmt.Println(a0)
	// Output: first
}
