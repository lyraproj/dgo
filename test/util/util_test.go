package util

import "testing"

func TestCheckFailed(t *testing.T) {
	CheckFailed(t, func(t2 *testing.T) {
		CheckFailed(t2, func(tt *testing.T) {
			tt.Error("no panic")
		}, true)
	}, false)
}

func TestCheckFailed_errorMsg_notFatal(t *testing.T) {
	CheckFailed(t, func(tt *testing.T) {
		tt.Error("no panic")
	}, false, "no panic")

	CheckFailed(t, func(tt *testing.T) {
		CheckFailed(tt, func(t2 *testing.T) {
			t2.Error("no panic")
		}, false, "don't worry")
	}, false)
}

func TestCheckFailed_errorMsg_fatal(t *testing.T) {
	CheckFailed(t, func(tt *testing.T) {
		tt.Fatal("no panic")
	}, true, "no panic")

	CheckFailed(t, func(tt *testing.T) {
		CheckFailed(tt, func(t2 *testing.T) {
			t2.Fatal("no panic")
		}, true, "don't worry")
	}, false)
}

type withStringer int

func (withStringer) String() string {
	return "the string"
}

func TestCheckPanic_stringer(t *testing.T) {
	CheckFailed(t, func(tt *testing.T) {
		CheckPanic(tt, func() {
			panic(withStringer(0))
		}, "not this", false)
	}, false, `recovered "the string" does not match "not this"`)
}

type notStringable int

func TestCheckPanic_not_stringable(t *testing.T) {
	CheckFailed(t, func(tt *testing.T) {
		CheckPanic(tt, func() {
			panic(notStringable(0))
		}, "not this", false)
	}, false, `recovered "util.notStringable 0" does not match "not this"`)
}
