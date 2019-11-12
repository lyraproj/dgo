package tf_test

import (
	"fmt"

	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func ExampleMap_min() {
	tp := tf.Map(1)
	fmt.Println(tp.Instance(vf.Map(`a`, 42)))
	fmt.Println(tp.Instance(vf.Map()))
	// Output:
	// true
	// false
}

func ExampleMap_min_max() {
	tp := tf.Map(1, 2)
	fmt.Println(tp.Instance(vf.Map(`a`, 42, `b`, 84)))
	fmt.Println(tp.Instance(vf.Map(`a`, 42, `b`, 84, `c`, 126)))
	// Output:
	// true
	// false
}

func ExampleMap_type_min() {
	tp := tf.Map(typ.String, typ.String, 2)
	fmt.Println(tp.Instance(vf.Map(`hello`, `word`, `hi`, `earth`)))
	fmt.Println(tp.Instance(vf.Map(`hello`, `word`)))
	// Output:
	// true
	// false
}

func ExampleMap_type_min_max() {
	tp := tf.Map(typ.String, typ.String, 2, 3)
	fmt.Println(tp.Instance(vf.Map(`hello`, `word`, `hi`, `earth`, `good day`, `sunshine`)))
	fmt.Println(tp.Instance(vf.Map(`hello`, `word`, `hi`, `earth`, `good day`, `sunshine`, `howdy`, `galaxy`)))
	// Output:
	// true
	// false
}
