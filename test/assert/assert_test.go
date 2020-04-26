// +build test

package assert

import (
	"errors"
	"regexp"
	"testing"

	"github.com/tada/dgo/test/util"
	"github.com/tada/dgo/typ"
	"github.com/tada/dgo/vf"
)

func failed(t *testing.T, f func(t *testing.T)) {
	util.CheckFailed(t, f, false)
}

// Specific tests that causes failures in the tester to test that those failures
// are indeed handled.
func TestTheTester(t *testing.T) {
	failed(t, func(ft *testing.T) {
		Assignable(ft, typ.String, typ.Integer)
	})
	failed(t, func(ft *testing.T) {
		NotAssignable(ft, typ.String, typ.String)
	})
	failed(t, func(ft *testing.T) {
		Instance(ft, typ.String, 2)
	})
	failed(t, func(ft *testing.T) {
		NotInstance(ft, typ.String, `a`)
	})
	failed(t, func(ft *testing.T) {
		Equal(ft, `a`, `b`)
	})
	failed(t, func(ft *testing.T) {
		NotEqual(ft, `a`, `a`)
	})
	failed(t, func(ft *testing.T) {
		Match(ft, regexp.MustCompile(`foo`), `bar`)
	})
	failed(t, func(ft *testing.T) {
		Match(ft, 23, `bar`)
	})
	failed(t, func(ft *testing.T) {
		Match(ft, `bar`, 23)
	})
	Match(t, `xyz`, `has xyz in it`)
	Match(t, vf.String(`xyz`), `has xyz in it`)
	failed(t, func(ft *testing.T) {
		NoMatch(ft, `bar`, `bar`)
	})
	NoMatch(t, `abc`, `has xyz in it`)
	failed(t, func(ft *testing.T) {
		NoError(ft, errors.New(`nope`))
	})
	failed(t, func(ft *testing.T) {
		Error(ft, `yep`, nil)
	})
	failed(t, func(ft *testing.T) {
		Error(ft, `yep`, errors.New(`nope`))
	})
	failed(t, func(ft *testing.T) {
		Same(ft, `a`, `b`)
	})
	failed(t, func(ft *testing.T) {
		NotSame(ft, `a`, `a`)
	})
	failed(t, func(ft *testing.T) {
		False(ft, true)
	})
	failed(t, func(ft *testing.T) {
		True(ft, false, `not`, `true`)
	})
	failed(t, func(ft *testing.T) {
		Nil(ft, true)
	})
	failed(t, func(ft *testing.T) {
		NotNil(ft, nil)
	})
	failed(t, func(ft *testing.T) {
		Panic(ft, func() {}, `this`)
	})
	failed(t, func(ft *testing.T) {
		Panic(ft, func() { panic(`this`) }, `that`)
	})
	failed(t, func(ft *testing.T) {
		Panic(ft, func() { panic(errors.New(`this`)) }, `that`)
	})
}
