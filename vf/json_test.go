package vf

import (
	"fmt"
)

func ExampleUnmarshalJSON() {
	v, err := UnmarshalJSON([]byte(`["hello",true,1,3.14,null,{"a":1}]`))
	if err == nil {
		fmt.Println(v.Equals(Values(`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1})))
	}
	// Output: true
}

func ExampleMarshalJSON_slice() {
	v, err := MarshalJSON(Values(
		`hello`, true, 1, 3.14, nil, map[string]interface{}{"a": 1}))
	if err == nil {
		fmt.Println(string(v))
	}
	// Output: ["hello",true,1,3.14,null,{"a":1}]
}

func ExampleMarshalJSON_string() {
	v, err := MarshalJSON("hello")
	if err == nil {
		fmt.Println(string(v))
	}
	// Output: "hello"
}
