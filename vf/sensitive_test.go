package vf_test

import (
	"fmt"

	"github.com/tada/dgo/vf"
)

func ExampleSensitive() {
	s := vf.Sensitive("don't reveal this in logs")
	fmt.Println(s)
	// Output: sensitive [value redacted]
}
