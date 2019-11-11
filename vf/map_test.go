package vf_test

import (
	"fmt"

	"github.com/lyraproj/dgo/util"
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
	m.SetType(`map[string]0..0x7f`)
	m.Put(`a`, 32)
	fmt.Println(m)
	// Output: {"a":32}
}

func ExampleMutableMap_illegalAssignment() {
	m := vf.MutableMap()
	m.SetType(`map[string]0..0x7f`)
	if err := util.Catch(func() { m.Put(`c`, 132) }); err != nil {
		fmt.Println(err)
	}
	// Output: the value 132 cannot be assigned to a variable of type 0..127
}
