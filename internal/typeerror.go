package internal

import (
	"fmt"

	"github.com/lyraproj/dgo/dgo"
	"github.com/tada/catch"
)

// IllegalAssignment returns the error that represents an assignment type constraint mismatch
func IllegalAssignment(t dgo.Type, v dgo.Value) error {
	var what string
	switch actual := v.Type().(type) {
	case *hstring:
		what = fmt.Sprintf(`the string %q`, actual.string)
	default:
		if dgo.IsExact(actual) {
			what = fmt.Sprintf(`the value %v`, actual)
		} else {
			what = fmt.Sprintf(`a value of type %s`, TypeString(actual))
		}
	}
	return catch.Error("%v cannot be assigned to a variable of type %s", what, TypeString(t))
}

// IllegalMapKey returns the error that represents an assignment map key constraint mismatch
func IllegalMapKey(t dgo.Type, v dgo.Value) error {
	return catch.Error("key %q cannot be added to type %s", v, TypeString(t))
}

// IllegalSize returns the error that represents an size constraint mismatch
func IllegalSize(t dgo.Type, sz int) error {
	return catch.Error("size constraint violation on type %s when attempting resize to %d", TypeString(t), sz)
}
