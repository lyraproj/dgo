package internal_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/test/require"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestErrorType(t *testing.T) {
	tp := typ.Error
	assert.Same(t, typ.Error, typ.Generic(tp))
	assert.Instance(t, tp, errors.New(`an error`))
	assert.Instance(t, tp.Type(), tp)
	assert.Equal(t, tp, tp)
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.True(t, reflect.TypeOf((*error)(nil)).Elem().AssignableTo(tp.ReflectType()))
	assert.Equal(t, `error`, tp.String())

	er := errors.New(`some error`)
	assert.True(t, tp.IsInstance(er))

	v := vf.Value(er)
	tp = v.Type().(dgo.ErrorType)
	assert.Same(t, typ.Error, typ.Generic(tp))
	assert.NotSame(t, tp, typ.Generic(tp))
	assert.Instance(t, tp, v)
	assert.Instance(t, tp, er)
	assert.Instance(t, tp.Type(), tp)
	assert.Assignable(t, typ.Error, tp)
	assert.NotAssignable(t, tp, typ.Error)
	assert.Assignable(t, tp, tp)
	assert.NotAssignable(t, tp, typ.String)
	assert.Equal(t, tp, tp)
	assert.NotEqual(t, typ.Error, tp)
	assert.Instance(t, tp.Type(), tp)
	assert.True(t, tp.IsInstance(er))
	assert.False(t, tp.IsInstance(errors.New(`some other error`)))
	assert.Equal(t, tp.HashCode(), tp.HashCode())
	assert.NotEqual(t, 0, tp.HashCode())
	assert.Equal(t, `error "some error"`, tp.String())
	assert.True(t, reflect.TypeOf(er).AssignableTo(tp.ReflectType()))
}

type testUnwrapError struct {
	err error
}

func (e *testUnwrapError) Error() string {
	return `unwrap me`
}

func (e *testUnwrapError) Unwrap() error {
	return e.err
}

func TestError(t *testing.T) {
	ev := errors.New(`some error`)
	v := vf.Value(ev)

	ve, ok := v.(error)
	require.True(t, ok)
	assert.Equal(t, v, ev)
	assert.Equal(t, v, errors.New(`some error`))
	assert.Equal(t, v, vf.Value(errors.New(`some error`)))
	assert.NotEqual(t, v, `some error`)
	assert.NotEqual(t, 0, v.HashCode())
	assert.Equal(t, v.HashCode(), vf.Value(errors.New(`some error`)).HashCode())
	assert.NotEqual(t, v.HashCode(), vf.Value(errors.New(`some other error`)).HashCode())

	assert.Equal(t, `some error`, ve.Error())

	type uvt interface {
		Unwrap() error
	}

	u, ok := v.(uvt)
	require.True(t, ok)
	assert.Nil(t, u.Unwrap())

	uv := vf.Value(&testUnwrapError{ev})
	u, ok = uv.(uvt)
	require.True(t, ok)
	assert.Equal(t, u.Unwrap(), ev)
}

func TestError_ReflectTo(t *testing.T) {
	var err error
	v := vf.Value(errors.New(`some error`))
	vf.ReflectTo(v, reflect.ValueOf(&err).Elem())
	assert.NotNil(t, err)
	assert.Equal(t, `some error`, err.Error())

	var bp *error
	vf.ReflectTo(v, reflect.ValueOf(&bp).Elem())
	assert.NotNil(t, bp)
	assert.NotNil(t, *bp)
	assert.Equal(t, `some error`, (*bp).Error())

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	ec, ok := mi.(error)
	require.True(t, ok)
	assert.Same(t, err, ec)
}

func TestError_TypeString(t *testing.T) {
	assert.Equal(t, `error "some error"`, internal.TypeString(vf.Value(errors.New(`some error`)).Type()))
}
