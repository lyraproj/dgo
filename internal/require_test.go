package internal_test

import (
	"errors"
	"testing"

	require "github.com/lyraproj/got/dgo_test"
	"github.com/lyraproj/got/typ"
)

func ensureFailed(t *testing.T, f func(t *testing.T)) {
	tt := testing.T{}
	f(&tt)
	if !tt.Failed() {
		t.Fail()
	}
}

// Specific tests that causes failures in the tester to test that those failures
// are indeed handled.
func TestTheTester(t *testing.T) {
	ensureFailed(t, func(ft *testing.T) {
		require.Assignable(ft, typ.String, typ.Integer)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.NotAssignable(ft, typ.String, typ.String)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Instance(ft, typ.String, 2)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.NotInstance(ft, typ.String, `a`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Equal(ft, `a`, `b`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.NotEqual(ft, `a`, `a`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Same(ft, `a`, `b`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.NotSame(ft, `a`, `a`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.False(ft, true)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.True(ft, false, `not`, `true`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Nil(ft, true)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.NotNil(ft, nil)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Panic(ft, func() {}, `this`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Panic(ft, func() { panic(`this`) }, `that`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Panic(ft, func() { panic(errors.New(`this`)) }, `that`)
	})
}
