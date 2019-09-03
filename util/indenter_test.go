package util_test

import (
	"fmt"

	"github.com/lyraproj/got/dgo"

	"github.com/lyraproj/got/util"
	"github.com/lyraproj/got/vf"
)

func ExampleToString() {
	v := vf.Values(`a`, map[string]int{`b`: 3}).(dgo.Array)
	fmt.Println(util.ToString(v))
	// Output: ["a",{"b":3}]
}

func ExampleToIndentedString() {
	v := vf.Values(`a`, map[string]int{`b`: 3}).(dgo.Array)
	fmt.Println(util.ToIndentedString(v))
	// Output:
	// [
	//  "a",
	//  {
	//   "b": 3
	//  }
	// ]
}
