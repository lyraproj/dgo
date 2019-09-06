package require

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

func errorlog(t *testing.T, dflt string, args []interface{}) {
	t.Helper()
	if len(args) > 0 {
		t.Log(args...)
	} else {
		t.Log(dflt)
	}
	t.Fail()
}

// Assignable will fail unless a is assignable from b
func Assignable(t *testing.T, a, b dgo.Type) {
	t.Helper()
	if !a.Assignable(b) {
		t.Errorf(`%s is not a assignable from %s`, a, b)
	}
}

// NotAssignable will fail if a is assignable from b
func NotAssignable(t *testing.T, a, b dgo.Type) {
	t.Helper()
	if a.Assignable(b) {
		t.Errorf(`%s is assignable from %s`, a, b)
	}
}

// Instance will fail unless val is an instance of typ
func Instance(t *testing.T, typ dgo.Type, val interface{}) {
	t.Helper()
	if !typ.Instance(val) {
		t.Errorf(`%v is not an instance of %s`, val, typ)
	}
}

// NotInstance will fail if val is an instance of typ
func NotInstance(t *testing.T, typ dgo.Type, val interface{}) {
	t.Helper()
	if typ.Instance(val) {
		t.Errorf(`%v is an instance of %s`, val, typ)
	}
}

// Equal will fail unless a is equal to b
func Equal(t *testing.T, a, b interface{}) {
	t.Helper()
	if internal.Value(a).Equals(b) {
		return
	}
	t.Errorf(`%v is not equal to %v`, a, b)
}

// NotEqual will fail if a is equal to b
func NotEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if internal.Value(a).Equals(b) {
		t.Errorf(`%v is equal to %v`, a, b)
	}
}

// Same will fail unless a and b are the same valuess
func Same(t *testing.T, a, b interface{}) {
	t.Helper()
	if a != b {
		t.Error(`not same instance`)
	}
}

// NotSame will fail if a and b are the same values
func NotSame(t *testing.T, a, b interface{}) {
	t.Helper()
	if a == b {
		t.Error(`same instance`)
	}
}

// False will fail unless v is false
func False(t *testing.T, v bool, args ...interface{}) {
	t.Helper()
	if v {
		errorlog(t, `not false`, args)
	}
}

// True will fail unless v is true
func True(t *testing.T, v bool, args ...interface{}) {
	t.Helper()
	if !v {
		errorlog(t, `not true`, args)
	}
}

// Nil will fail unless v is nil
func Nil(t *testing.T, v interface{}) {
	t.Helper()
	if v != internal.Nil && v != nil {
		t.Errorf(`%v is not nil`, v)
	}
}

// NotNil will fail if v is nil
func NotNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		t.Errorf(`%v is nil`, v)
	}
}

// Panic will fail unless a call to f results in a panic and the recovered value matches v
func Panic(t *testing.T, f func(), v interface{}) {
	t.Helper()
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf(`%v`, r)
				}
			}
		}()
		f()
	}()

	if err == nil {
		t.Errorf(`expect panic matching '%v' did not occur`, v)
		return
	}

	switch v := v.(type) {
	case string:
		if regexp.MustCompile(v).MatchString(err.Error()) {
			return
		}
	case dgo.Value:
		if v.Equals(err) {
			return
		}
	}
	t.Errorf(`recovered "%s" does not match "%v"`, err.Error(), v)
}
