package internal_test

import (
	"errors"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func TestErrorType(t *testing.T) {
	er := errors.New(`some error`)
	v := vf.Value(er)
	tp := v.Type()
	require.Same(t, typ.Error, tp)
	require.Instance(t, tp, v)
	require.Instance(t, tp, er)
	require.Assignable(t, tp, tp)
	require.NotAssignable(t, tp, typ.String)
	require.Equal(t, tp, tp)
	require.Instance(t, tp.Type(), tp)

	require.Equal(t, tp.HashCode(), tp.HashCode())
	require.NotEqual(t, 0, tp.HashCode())

	require.Equal(t, `error`, tp.String())
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

	require.Equal(t, ve.Error(), v.String())

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
