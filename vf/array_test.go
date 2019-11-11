package vf_test

import (
	"fmt"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/typ"

	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func ExampleValues() {
	fmt.Println(vf.Values(`hello`, 1, 2.3, true, []byte{1, 2, 3}))
	// Output: {"hello",1,2.3,true,AQID}
}

func ExampleStrings() {
	fmt.Println(vf.Strings(`one`, `two`))
	// Output: {"one","two"}
}

func ExampleMutableValues() {
	a := vf.MutableValues()
	a.SetType(`[]0..0x7f`)
	a.Add(32)
	fmt.Println(a)
	// Output: {32}
}

func ExampleMutableValues_illegalAssignment() {
	a := vf.MutableValues()
	a.SetType(`[]0..0x7f`)
	if err := util.Catch(func() { a.Add(132) }); err != nil {
		fmt.Println(err)
	}
	// Output: the value 132 cannot be assigned to a variable of type 0..127
}

func ExampleArguments() {
	args := vf.Arguments(`first`, 2)
	a0 := args.Arg(`myFunc`, 0, typ.String).(dgo.String)
	fmt.Println(a0)
	// Output: first
}
