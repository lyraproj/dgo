package tf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// IllegalAssignment returns the error that represents an assignment type constraint mismatch
func IllegalAssignment(expected dgo.Type, actual dgo.Value) error {
	return internal.IllegalAssignment(expected, actual)
}

// IllegalSize returns the error that represents an size constraint mismatch
func IllegalSize(expected dgo.Type, size int) error {
	return internal.IllegalSize(expected, size)
}

// IllegalMapKey returns the error that represents an assignment map key constraint mismatch
func IllegalMapKey(t dgo.Type, v dgo.Value) error {
	return internal.IllegalMapKey(t, v)
}
