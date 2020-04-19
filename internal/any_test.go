package internal_test

import (
	"regexp"
	"testing"

	"github.com/tada/dgo/dgo"
	require "github.com/tada/dgo/dgo_test"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
)

func TestAny(t *testing.T) {
	require.Equal(t, typ.Any, typ.Any)
	require.NotEqual(t, typ.Any, typ.Boolean)
	require.Instance(t, typ.Any, 3)
	require.Instance(t, typ.Any, `foo`)
	require.Assignable(t, typ.Any, typ.String)
	require.Assignable(t, typ.Any, tf.String(3, 3))
	require.Assignable(t, typ.Any, tf.Pattern(regexp.MustCompile(`f`)))
	require.Assignable(t, typ.Any, tf.Enum(`f`, `foo`, `foobar`))
	require.Assignable(t, typ.Any.Type(), typ.Any.Type())
	require.Instance(t, typ.Any.Type(), typ.Any)
	require.Instance(t, typ.Any.Type(), typ.Boolean)
	require.Equal(t, typ.Any.HashCode(), typ.Any.HashCode())
	require.NotEqual(t, 0, typ.Any.HashCode())

	// Yes, since the Not is more constrained
	require.Assignable(t, typ.Any, tf.Not(typ.Any))

	require.Equal(t, `any`, typ.Any.String())
	require.Equal(t, dgo.TiAny, typ.Any.TypeIdentifier())
}
