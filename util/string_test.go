package util_test

import (
	"fmt"

	"github.com/lyraproj/dgo/util"
)

func ExampleFtoa() {
	fmt.Println(util.Ftoa(3))
	// Output: 3.0
}

func ExampleContainsString() {
	fmt.Println(util.ContainsString([]string{`foo`, `fee`, `fum`}, `fee`))
	fmt.Println(util.ContainsString([]string{`no`, `such`, `text`}, `fee`))
	// Output:
	// true
	// false
}
