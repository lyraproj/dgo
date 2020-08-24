package internal_test

import (
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestIllegalAssignment(t *testing.T) {
	v := tf.IllegalAssignment(typ.String, vf.Integer(3))
	assert.Equal(t, v, tf.IllegalAssignment(typ.String, vf.Integer(3)))
	assert.NotEqual(t, v, tf.IllegalAssignment(typ.String, vf.Integer(4)))
	assert.NotEqual(t, v, `oops`)
	assert.Instance(t, v.Type(), v)
	assert.NotEqual(t, 0, v.HashCode())
	assert.Equal(t, v.HashCode(), v.HashCode())
	assert.Equal(t, `the value 3 cannot be assigned to a variable of type string`, v.String())
}

func TestIllegalSize(t *testing.T) {
	v := tf.IllegalSize(tf.String(1, 10), 12)
	assert.Equal(t, v, tf.IllegalSize(tf.String(1, 10), 12))
	assert.NotEqual(t, v, tf.IllegalSize(tf.String(1, 10), 11))
	assert.NotEqual(t, v, `oops`)
	assert.Instance(t, v.Type(), v)
	assert.NotEqual(t, 0, v.HashCode())
	assert.Equal(t, v.HashCode(), v.HashCode())
	assert.Equal(t, `size constraint violation on type string[1,10] when attempting resize to 12`, v.String())
}

func TestIllegalMapKey(t *testing.T) {
	tp := tf.ParseType(`{a:string}`).(dgo.StructMapType)
	v := tf.IllegalMapKey(tp, vf.String(`b`))
	assert.Equal(t, v, tf.IllegalMapKey(tp, vf.String(`b`)))
	assert.NotEqual(t, v, tf.IllegalMapKey(tp, vf.String(`c`)))
	assert.NotEqual(t, v, `oops`)
	assert.Instance(t, v.Type(), v)
	assert.NotEqual(t, 0, v.HashCode())
	assert.Equal(t, v.HashCode(), v.HashCode())
	assert.Equal(t, `key "b" cannot be added to type {"a":string}`, v.String())
}
