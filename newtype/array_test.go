package newtype_test

import (
	"fmt"

	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func ExampleArray() {
	tp := newtype.Array()
	fmt.Println(tp.Equals(typ.Array))
	fmt.Println(tp.Instance(vf.Values(`hello`, 42)))
	fmt.Println(tp.Instance(42))
	// Output:
	// true
	// true
	// false
}

func ExampleArray_min() {
	tp := newtype.Array(2)
	fmt.Println(tp.Instance(vf.Values(`hello`, 42)))
	fmt.Println(tp.Instance(vf.Values(`hello`)))
	// Output:
	// true
	// false
}

func ExampleArray_type() {
	tp := newtype.Array(typ.String)
	fmt.Println(tp.Instance(vf.Values(`hello`)))
	fmt.Println(tp.Instance(vf.Values(42)))
	// Output:
	// true
	// false
}

func ExampleArray_min_max() {
	tp := newtype.Array(1, 2)
	fmt.Println(tp.Instance(vf.Values(`hello`, 42)))
	fmt.Println(tp.Instance(vf.Values(`hello`, 42, `word`)))
	// Output:
	// true
	// false
}

func ExampleArray_type_min() {
	// Create a new array type with a minimum size of 2
	tp := newtype.Array(typ.String, 2)
	fmt.Println(tp.Instance(vf.Values(`hello`, `word`)))
	fmt.Println(tp.Instance(vf.Values(`hello`)))
	// Output:
	// true
	// false
}

func ExampleArray_type_min_max() {
	tp := newtype.Array(typ.String, 2, 3)
	fmt.Println(tp.Instance(vf.Values(`hello`, `word`)))
	// Output: true
}
