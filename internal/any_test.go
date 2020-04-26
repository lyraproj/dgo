package internal_test

import (
	"regexp"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
)

func TestAny(t *testing.T) {
	assert.Equal(t, typ.Any, typ.Any)
	assert.NotEqual(t, typ.Any, typ.Boolean)
	assert.Instance(t, typ.Any, 3)
	assert.Instance(t, typ.Any, `foo`)
	assert.Assignable(t, typ.Any, typ.String)
	assert.Assignable(t, typ.Any, tf.String(3, 3))
	assert.Assignable(t, typ.Any, tf.Pattern(regexp.MustCompile(`f`)))
	assert.Assignable(t, typ.Any, tf.Enum(`f`, `foo`, `foobar`))
	assert.Assignable(t, typ.Any.Type(), typ.Any.Type())
	assert.Instance(t, typ.Any.Type(), typ.Any)
	assert.Instance(t, typ.Any.Type(), typ.Boolean)
	assert.Equal(t, typ.Any.HashCode(), typ.Any.HashCode())
	assert.NotEqual(t, 0, typ.Any.HashCode())

	// Yes, since the Not is more constrained
	assert.Assignable(t, typ.Any, tf.Not(typ.Any))
	assert.Equal(t, `any`, typ.Any.String())
	assert.Equal(t, dgo.TiAny, typ.Any.TypeIdentifier())
}
