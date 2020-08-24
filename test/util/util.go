// +build test

// Package util contains utility function for the assert and require packages.
package util

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"unsafe"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// CheckFailed checks that the given test function fails with an Error or Fatal
func CheckFailed(t *testing.T, f func(t *testing.T), expectFatal bool, substrings ...string) {
	tt := testing.T{}
	x := make(chan bool, 1)
	fatal := true
	go func() {
		defer func() { x <- true }() // GoExit runs all deferred calls
		f(&tt)
		fatal = false
	}()
	<-x
	if !tt.Failed() || expectFatal != fatal {
		t.Fail()
	}
	if len(substrings) > 0 {
		// Pick the output bytes from the testing.T using an unsafe.Pointer.
		rs := reflect.ValueOf(&tt).Elem()
		rf := rs.FieldByName("common").FieldByName("output")
		// #nosec
		rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
		le := string(rf.Interface().([]byte))
		for _, ss := range substrings {
			if !strings.Contains(le, ss) {
				t.Errorf("string %q does not contain %q", le, ss)
			}
		}
	}
}

// IsMatch is a helper method that checks if a, which can be a regexp, a string, or
// a dgo.String() matches b, which must be a string or a dgo.String.
func IsMatch(a, b interface{}) bool {
	if sv, ok := internal.Value(b).(dgo.String); ok {
		var rx *regexp.Regexp
		switch a := a.(type) {
		case *regexp.Regexp:
			rx = a
		case string:
			rx = regexp.MustCompile(a)
		case dgo.String:
			rx = regexp.MustCompile(a.GoString())
		default:
			return false
		}
		return rx.MatchString(sv.GoString())
	}
	return false
}

// CheckPanic will fail unless a call to f results in a panic and the recovered value matches v
func CheckPanic(t *testing.T, f func(), v string, failNow bool) {
	t.Helper()
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				switch r := r.(type) {
				case error:
					err = r
				case string:
					err = errors.New(r)
				case fmt.Stringer:
					err = errors.New(r.String())
				default:
					err = fmt.Errorf(`%T %v`, r, r)
				}
			}
		}()
		f()
	}()

	if err != nil && regexp.MustCompile(v).MatchString(err.Error()) {
		return
	}

	if err == nil {
		t.Logf(`expect panic matching '%v' did not occur`, v)
	} else {
		t.Logf(`recovered %q does not match %q`, err.Error(), v)
	}
	if failNow {
		t.FailNow()
	} else {
		t.Fail()
	}
}
