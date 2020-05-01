// +build test

// Package assert contains the test helper functions for easy testing of assignability, instance of, equality, and
// panic assertions. When an assertion is false, all methods here will fail the test but continue execution.
package assert

import (
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
	"github.com/tada/dgo/test/util"
)

func errorlog(t *testing.T, dflt string, args []interface{}) {
	t.Helper()
	if len(args) > 0 {
		t.Error(args...)
	} else {
		t.Error(dflt)
	}
}

// Assignable will fail unless a is assignable from b
func Assignable(t *testing.T, a, b dgo.Type) {
	if !a.Assignable(b) {
		t.Helper()
		t.Errorf(`%v is not a assignable from %v`, a, b)
	}
}

// NotAssignable will fail if a is assignable from b
func NotAssignable(t *testing.T, a, b dgo.Type) {
	if a.Assignable(b) {
		t.Helper()
		t.Errorf(`%v is assignable from %v`, a, b)
	}
}

// Instance will fail unless val is an instance of typ
func Instance(t *testing.T, typ dgo.Type, val interface{}) {
	if !typ.Instance(val) {
		t.Helper()
		t.Errorf(`%v is not an instance of %v`, internal.Value(val), typ)
	}
}

// NotInstance will fail if val is an instance of typ
func NotInstance(t *testing.T, typ dgo.Type, val interface{}) {
	if typ.Instance(val) {
		t.Helper()
		t.Errorf(`%v is an instance of %v`, internal.Value(val), typ)
	}
}

// Equal will fail unless a is equal to b
func Equal(t *testing.T, a, b interface{}) {
	va := internal.Value(a)
	if !va.Equals(b) {
		t.Helper()
		t.Errorf(`%v is not equal to %v`, va, internal.Value(b))
	}
}

// NotEqual will fail if a is equal to b
func NotEqual(t *testing.T, a, b interface{}) {
	va := internal.Value(a)
	if va.Equals(b) {
		t.Helper()
		t.Errorf(`%v is equal to %v`, va, internal.Value(b))
	}
}

// Same will fail unless a and b are the same valuess
func Same(t *testing.T, a, b interface{}) {
	if a != b {
		t.Helper()
		t.Error(`not same instance`)
	}
}

// NotSame will fail if a and b are the same values
func NotSame(t *testing.T, a, b interface{}) {
	if a == b {
		t.Helper()
		t.Error(`same instance`)
	}
}

// Match will fail unless b matches the regexp a
func Match(t *testing.T, a, b interface{}) {
	if !util.IsMatch(a, b) {
		t.Helper()
		t.Errorf(`'%v' does not match '%v'`, internal.Value(b), internal.Value(a))
	}
}

// NoMatch will fail if b matches the regexp a
func NoMatch(t *testing.T, a, b interface{}) {
	if util.IsMatch(a, b) {
		t.Helper()
		t.Errorf(`'%v' matches '%v'`, internal.Value(b), internal.Value(a))
	}
}

// False will fail unless v is false
func False(t *testing.T, v bool, args ...interface{}) {
	if v {
		t.Helper()
		errorlog(t, `not false`, args)
	}
}

// True will fail unless v is true
func True(t *testing.T, v bool, args ...interface{}) {
	if !v {
		t.Helper()
		errorlog(t, `not true`, args)
	}
}

// Error will fail unless error matches the given pattern
func Error(t *testing.T, msg interface{}, err error) {
	t.Helper()
	if err == nil {
		t.Errorf(`Expected error matching '%v' did not occur`, internal.Value(msg))
	} else if !util.IsMatch(msg, err.Error()) {
		t.Errorf(`Expected error matching '%v', got '%s'`, internal.Value(msg), err.Error())
	}
}

// NoError will fail unless error is nil
func NoError(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Errorf(`Expected no error, got '%v'`, err.Error())
	}
}

// Nil will fail unless v is nil
func Nil(t *testing.T, v interface{}) {
	if !(internal.Nil == v || v == nil) {
		t.Helper()
		t.Errorf(`%v is not nil`, internal.Value(v))
	}
}

// NotNil will fail if v is nil
func NotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Helper()
		t.Errorf(`%v is nil`, internal.Value(v))
	}
}

// Panic will fail unless a call to f results in a panic and the recovered value matches v
func Panic(t *testing.T, f func(), v string) {
	t.Helper()
	util.CheckPanic(t, f, v, false)
}
