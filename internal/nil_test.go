package internal_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func TestNil_AppendTo(t *testing.T) {
	assert.Equal(t, `nil`, util.ToStringERP(vf.Nil))
}

func TestNil_CompareTo(t *testing.T) {
	c, ok := vf.Nil.CompareTo(vf.Nil)
	assert.True(t, ok)
	assert.Equal(t, 0, c)

	c, ok = vf.Nil.CompareTo(vf.Int64(1))
	assert.True(t, ok)
	assert.Equal(t, -1, c)
	assert.Nil(t, vf.Nil.GoNil())
}

func TestNil_ReflectTo(t *testing.T) {
	err := errors.New(`some error`)
	vf.Nil.ReflectTo(reflect.ValueOf(&err).Elem())
	assert.True(t, err == nil)

	err = errors.New(`some error`)
	ep := &err
	vf.Nil.ReflectTo(reflect.ValueOf(&ep).Elem())
	assert.True(t, ep == nil)

	var mi interface{} = err
	mip := &mi
	vf.Nil.ReflectTo(reflect.ValueOf(mip).Elem())
	assert.True(t, mi == nil)
}

func TestNilType(t *testing.T) {
	nt := typ.Nil
	assert.Assignable(t, nt, nt)
	assert.Instance(t, nt, vf.Nil)
	assert.Assignable(t, typ.Any, nt)
	assert.NotAssignable(t, nt, typ.Any)
	assert.Equal(t, nt, nt)
	assert.NotEqual(t, nt, typ.Any)
	assert.Instance(t, nt.Type(), nt)
	assert.NotInstance(t, nt.Type(), typ.Any)
	assert.Equal(t, `nil`, nt.String())
	assert.NotEqual(t, 0, nt.HashCode())
	assert.Equal(t, nt.HashCode(), nt.HashCode())
	assert.Equal(t, nt.ReflectType(), typ.Any.ReflectType())
}
