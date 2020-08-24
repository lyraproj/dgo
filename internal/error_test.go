package internal_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestErrorType(t *testing.T) {
	tp := typ.Error
	require.Same(t, typ.Error, typ.Generic(tp))
	require.Instance(t, tp.Type(), tp)
	require.Equal(t, tp, tp)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())
	require.True(t, reflect.TypeOf((*error)(nil)).Elem().AssignableTo(tp.ReflectType()))
	require.Equal(t, `error`, tp.String())

	er := errors.New(`some error`)
	require.True(t, tp.IsInstance(er))

	v := vf.Value(er)
	tp = v.Type().(dgo.ErrorType)
	require.Same(t, typ.Error, typ.Generic(tp))
	require.NotSame(t, tp, typ.Generic(tp))
	require.Instance(t, tp, v)
	require.Instance(t, tp, er)
	require.Instance(t, tp.Type(), tp)
	require.Assignable(t, typ.Error, tp)
	require.NotAssignable(t, tp, typ.Error)
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.String)
	require.Equal(t, tp, tp)
	require.NotEqual(t, typ.Error, tp)
	require.Instance(t, tp.Type(), tp)
	require.True(t, tp.IsInstance(er))
	require.False(t, tp.IsInstance(errors.New(`some other error`)))

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `error "some error"`, tp.String())

	require.True(t, reflect.TypeOf(er).AssignableTo(tp.ReflectType()))
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

	require.Equal(t, v, ev)
	require.Equal(t, v, errors.New(`some error`))
	require.Equal(t, v, vf.Value(errors.New(`some error`)))
	require.NotEqual(t, v, `some error`)

	require.NotEqual(t, 0, v.HashCode())
	require.Equal(t, v.HashCode(), vf.Value(errors.New(`some error`)).HashCode())
	require.NotEqual(t, v.HashCode(), vf.Value(errors.New(`some other error`)).HashCode())

	require.Equal(t, `some error`, ve.Error())

	type uvt interface {
		Unwrap() error
	}

	u, ok := v.(uvt)
	require.True(t, ok)
	require.Nil(t, u.Unwrap())

	uv := vf.Value(&testUnwrapError{ev})
	u, ok = uv.(uvt)
	require.True(t, ok)
	require.Equal(t, u.Unwrap(), ev)
}

func TestError_ReflectTo(t *testing.T) {
	var err error
	v := vf.Value(errors.New(`some error`))
	vf.ReflectTo(v, reflect.ValueOf(&err).Elem())
	require.NotNil(t, err)
	require.Equal(t, `some error`, err.Error())

	var bp *error
	vf.ReflectTo(v, reflect.ValueOf(&bp).Elem())
	require.NotNil(t, bp)
	require.NotNil(t, *bp)
	require.Equal(t, `some error`, (*bp).Error())

	var mi interface{}
	mip := &mi
	vf.ReflectTo(v, reflect.ValueOf(mip).Elem())
	ec, ok := mi.(error)
	require.True(t, ok)
	require.Same(t, err, ec)
}

func TestError_TypeString(t *testing.T) {
	require.Equal(t, `error "some error"`, internal.TypeString(vf.Value(errors.New(`some error`)).Type()))
}
