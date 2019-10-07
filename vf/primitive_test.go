package vf_test

import (
	"fmt"
	"regexp"
	"time"

	"github.com/lyraproj/dgo/vf"
)

func ExampleRegexp() {
	rx := vf.Regexp(regexp.MustCompile(`[0-9a-fA-F]+`))
	fmt.Println(rx)
	// Output: [0-9a-fA-F]+
}

func ExampleTimeFromString() {
	t := vf.TimeFromString(`2019-10-06T07:15:00-07:00`)
	fmt.Println(t.GoTime().Format(time.RFC850))
	// Output: Sunday, 06-Oct-19 07:15:00 -0700
}
