package newtype_test

import (
	"fmt"

	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func ExampleMap_min() {
	tp := newtype.Map(1)
	fmt.Println(tp.Instance(vf.Map(map[string]int{`a`: 42})))
	fmt.Println(tp.Instance(vf.Map(map[string]int{})))
	// Output:
	// true
	// false
}

func ExampleMap_min_max() {
	tp := newtype.Map(1, 2)
	fmt.Println(tp.Instance(vf.Map(map[string]int{`a`: 42, `b`: 84})))
	fmt.Println(tp.Instance(vf.Map(map[string]int{`a`: 42, `b`: 84, `c`: 126})))
	// Output:
	// true
	// false
}

func ExampleMap_type_min() {
	tp := newtype.Map(typ.String, typ.String, 2)
	fmt.Println(tp.Instance(vf.Map(map[string]string{`hello`: `word`, `hi`: `earth`})))
	fmt.Println(tp.Instance(vf.Map(map[string]string{`hello`: `word`})))
	// Output:
	// true
	// false
}

func ExampleMap_type_min_max() {
	tp := newtype.Map(typ.String, typ.String, 2, 3)
	fmt.Println(tp.Instance(vf.Map(map[string]string{`hello`: `word`, `hi`: `earth`, `god day`: `sunshine`})))
	fmt.Println(tp.Instance(vf.Map(map[string]string{`hello`: `word`, `hi`: `earth`, `god day`: `sunshine`, `howdy`: `galaxy`})))
	// Output:
	// true
	// false
}
