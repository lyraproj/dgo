package vf_test

import (
	"fmt"
	"regexp"

	"github.com/lyraproj/dgo/vf"
)

func ExampleRegexp() {
	rx := vf.Regexp(regexp.MustCompile(`[0-9a-fA-F]+`))
	fmt.Println(rx)
	// Output: [0-9a-fA-F]+
}
