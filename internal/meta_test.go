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

func TestMeta(t *testing.T) {
	tp := typ.Any.Type()
	assert.Instance(t, tp, typ.Any)
	assert.Instance(t, tp, typ.Boolean)
	assert.NotInstance(t, tp, 3)
	assert.Same(t, typ.Any, tp.(dgo.Meta).Describes())
	assert.Equal(t, tf.Meta(typ.Any), tp)
	assert.Equal(t, typ.Any.Type(), tp)
	assert.NotEqual(t, typ.Boolean.Type(), tp)
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Assignable(t, tp, typ.Any.Type())
	assert.NotAssignable(t, tp, typ.Boolean.Type())
	assert.Equal(t, dgo.OpMeta, tp.(dgo.UnaryType).Operator())
	assert.Equal(t, `type`, tp.String())
	assert.Equal(t, `type[int]`, typ.Integer.Type().String())
	assert.True(t, reflect.ValueOf(typ.Any).Type().AssignableTo(tp.ReflectType()))
}

func TestMetaType_New(t *testing.T) {
	m := tf.String(1, 10)
	assert.Same(t, m, vf.New(typ.String.Type(), m))
	assert.Same(t, m, vf.New(m.Type(), m))
	assert.Same(t, m, vf.New(m.Type(), vf.Arguments(m)))
	assert.Equal(t, m, vf.New(typ.String.Type(), vf.String(`string[1,10]`)))
	assert.Equal(t, vf.String("x").Type(), vf.New(typ.String.Type(), vf.String(`"x"`)))
	assert.Panic(t, func() { vf.New(typ.Integer.Type(), vf.String(`string[1, 10]`)) }, `cannot be assigned`)
}

func TestMetaMeta(t *testing.T) {
	mt := typ.Any.Type().Type()
	assert.Equal(t, mt, typ.Boolean.Type().Type())
	assert.NotEqual(t, mt, typ.Any)
	assert.Equal(t, mt.Type(), mt)
	assert.Instance(t, mt, typ.Any.Type())
	assert.Instance(t, mt, typ.Boolean.Type())
	assert.NotInstance(t, mt, typ.Any)
	assert.Assignable(t, mt, mt)
	assert.NotAssignable(t, mt, typ.Any.Type())
	assert.Equal(t, dgo.OpMeta, mt.(dgo.UnaryType).Operator())
	assert.Equal(t, `type[type]`, mt.String())
}
