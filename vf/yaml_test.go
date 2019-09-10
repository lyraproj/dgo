package vf

import (
	"fmt"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
)

func ExampleUnmarshalYAML() {
	v, err := UnmarshalYAML([]byte(`
- hello
- true
- 1
- 3.14
- null
- a: 1`))
	if err == nil {
		fmt.Println(v)
	}
	// Output: ["hello",true,1,3.14,null,{"a":1}]
}

func TestUnmarshalYAML_fail(t *testing.T) {
	_, err := UnmarshalYAML([]byte(`: yaml`))
	require.NotNil(t, err)
}
