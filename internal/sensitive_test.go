package internal_test

import (
	"reflect"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func TestSensitiveType(t *testing.T) {
	var tp dgo.Type = typ.Sensitive
	s := vf.Sensitive(vf.Integer(0))
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.Any)
	assert.Assignable(t, tp, s.Type())
	assert.NotAssignable(t, s.Type(), tp)
	assert.Assignable(t, tf.Sensitive(typ.Integer), s.Type())
	assert.Instance(t, tp.Type(), s.Type())

	tp = s.Type()
	assert.Equal(t, tp, tp)
	assert.Equal(t, tp, vf.Sensitive(vf.Integer(0)).Type())
	assert.Equal(t, tp, vf.Sensitive(vf.Integer(1)).Type()) // type uses generic of wrapped
	assert.NotEqual(t, tp, typ.Any)
	assert.NotEqual(t, tp, tf.Array(typ.String))
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.Instance(t, tp, s)
	assert.NotInstance(t, tp, vf.Integer(0))
	assert.True(t, reflect.TypeOf(s).AssignableTo(tp.ReflectType()))
	assert.Equal(t, dgo.TiSensitive, tp.TypeIdentifier())
	assert.Equal(t, dgo.OpSensitive, tp.(dgo.UnaryType).Operator())
	assert.Equal(t, `sensitive[int]`, s.Type().String())
	assert.Equal(t, tf.Sensitive(), tf.ParseType(`sensitive`))
	assert.Equal(t, tf.Sensitive(typ.Integer), tf.ParseType(`sensitive[int]`))
	assert.Equal(t, vf.Sensitive(typ.Integer).Type(), tf.ParseType(`sensitive int`))
	assert.Equal(t, vf.Sensitive(34).Type(), tf.ParseType(`sensitive 34`))
	assert.Panic(t, func() { tf.ParseType(`sensitive[int, string]`) }, `illegal number of arguments`)
}

func TestSensitiveType_New(t *testing.T) {
	s := vf.Sensitive(`hide me`)
	assert.Equal(t, s, vf.New(typ.Sensitive, vf.Arguments(`hide me`)))
	assert.Equal(t, s, vf.New(typ.Sensitive, vf.String(`hide me`)))
	assert.Same(t, s, vf.New(typ.Sensitive, vf.Arguments(s)))
	assert.Same(t, s, vf.New(typ.Sensitive, s))
	assert.Panic(t, func() { vf.New(typ.Sensitive, vf.Arguments(`hide me`, `and me`)) }, `illegal number of arguments`)
}

func TestSensitive(t *testing.T) {
	s := vf.Sensitive(vf.String(`a`))
	assert.Equal(t, s, s)
	assert.Equal(t, s, vf.Sensitive(vf.String(`a`)))
	assert.NotEqual(t, s, vf.Sensitive(vf.String(`b`)))
	assert.NotEqual(t, s, vf.Strings(`a`))
	assert.True(t, s.Frozen())
	a := vf.MutableValues(`a`)
	s = vf.Sensitive(a)
	assert.False(t, s.Frozen())
	s.Freeze()
	assert.True(t, s.Frozen())
	assert.True(t, a.Frozen())
	assert.Same(t, s.Unwrap(), a)

	a = vf.MutableValues(`a`)
	s = vf.Sensitive(a)
	c := s.FrozenCopy().(dgo.Sensitive)
	assert.False(t, s.Frozen())
	assert.True(t, c.Frozen())
	assert.Equal(t, s.Unwrap(), c.Unwrap())
	assert.NotSame(t, s.Unwrap(), c.Unwrap())

	s = vf.Sensitive(vf.MutableValues(`a`))
	assert.False(t, s.Frozen())
	c = s.FrozenCopy().(dgo.Sensitive)
	assert.NotSame(t, s, c)
	assert.True(t, c.Frozen())
	s = c.FrozenCopy().(dgo.Sensitive)
	assert.Same(t, s, c)

	c = s.ThawedCopy().(dgo.Sensitive)
	assert.NotSame(t, s, c)
	assert.False(t, c.Frozen()) // c no longer frozen

	s = vf.Sensitive(vf.String(`a`))
	c = s.ThawedCopy().(dgo.Sensitive)
	assert.Same(t, s, c)

	assert.Equal(t, `sensitive [value redacted]`, s.String())

	assert.NotEqual(t, typ.Sensitive.HashCode(), s.HashCode())
	assert.Equal(t, s.HashCode(), s.HashCode())
}
