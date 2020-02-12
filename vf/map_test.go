package vf_test

import (
	"fmt"

	"github.com/lyraproj/dgo/vf"
)

func ExampleMap() {
	// when created with a list of key, value pairs, the order is retained within the map
	fmt.Println(vf.Map(`hello`, `word`, `hi`, `earth`, `good day`, `sunshine`))
	// Output: {"hello":"word","hi":"earth","good day":"sunshine"}
}

func ExampleMap_goMap() {
	// when created from a go map (always unordered), the map is sorted using natural
	// order of the map keys
	fmt.Println(vf.Map(map[string]string{`hello`: `word`, `hi`: `earth`, `good day`: `sunshine`}))
	// Output: {"good day":"sunshine","hello":"word","hi":"earth"}
}

func ExampleMutableMap() {
	m := vf.MutableMap()
	m.Put(`a`, 32)
	fmt.Println(m)
	// Output: {"a":32}
}
