package internal_test

import (
	"errors"
	"regexp"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
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
		require.Match(ft, regexp.MustCompile(`foo`), `bar`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Match(ft, 23, `bar`)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.Match(ft, `bar`, 23)
	})
	require.Match(t, `xyz`, `has xyz in it`)
	require.Match(t, vf.String(`xyz`), `has xyz in it`)
	ensureFailed(t, func(ft *testing.T) {
		require.NoMatch(ft, `bar`, `bar`)
	})
	require.NoMatch(t, `abc`, `has xyz in it`)
	ensureFailed(t, func(ft *testing.T) {
		require.Ok(ft, errors.New(`nope`))
	})
	ensureFailed(t, func(ft *testing.T) {
		require.NotOk(ft, `yep`, nil)
	})
	ensureFailed(t, func(ft *testing.T) {
		require.NotOk(ft, `yep`, errors.New(`nope`))
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
