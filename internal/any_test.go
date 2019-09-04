package internal_test

import (
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/typ"
)

func TestAny(t *testing.T) {
	require.Equal(t, typ.Any, typ.Any)
	require.NotEqual(t, typ.Any, typ.Boolean)
	require.Instance(t, typ.Any, 3)
	require.Instance(t, typ.Any, `foo`)
	require.Assignable(t, typ.Any, typ.String)
	require.Assignable(t, typ.Any, newtype.String(3, 3))
	require.Assignable(t, typ.Any, newtype.Pattern(regexp.MustCompile(`f`)))
	require.Assignable(t, typ.Any, newtype.Enum(`f`, `foo`, `foobar`))
	require.Assignable(t, typ.Any.Type(), typ.Any.Type())
	require.Instance(t, typ.Any.Type(), typ.Any)
	require.NotInstance(t, typ.Any.Type(), typ.Boolean)
	require.Equal(t, typ.Any.HashCode(), typ.Any.HashCode())
	require.NotEqual(t, 0, typ.Any.HashCode())

	// Yes, since the Not is more constrained
	require.Assignable(t, typ.Any, newtype.Not(typ.Any))

	require.Equal(t, `any`, typ.Any.String())
	require.Equal(t, dgo.TiAny, typ.Any.TypeIdentifier())
}
