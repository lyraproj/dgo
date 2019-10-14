package util

import (
	"errors"
	"fmt"
	"testing"
)

func ExampleCatch() {
	err := Catch(func() {
		panic(errors.New(`overheated`))
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Output: overheated
}

func ExampleCatch_string() {
	err := Catch(func() {
		panic(`overheated`)
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Output: overheated
}

func TestCatch_noPanic(t *testing.T) {
	err := Catch(func() {
		// nothing
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCatch_notError(t *testing.T) {
	defer func() {
		r, ok := recover().(int)
		if !(ok && r == 32) {
			t.Fatal(`did not recover expected int 32`)
		}
	}()
	_ = Catch(func() {
		panic(32)
	})
}
