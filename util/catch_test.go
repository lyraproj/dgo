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

func TestCatch_noError(t *testing.T) {
	defer func() {
		r, ok := recover().(string)
		if !(ok && r == `not an error`) {
			t.Fatal(`did not recover expected "not an error" string`)
		}
	}()
	_ = Catch(func() {
		panic(`not an error`)
	})
}
