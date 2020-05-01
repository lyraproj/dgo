package util_test

import (
	"errors"
	"testing"

	"github.com/lyraproj/dgo/test/assert"
	"github.com/lyraproj/dgo/util"
)

type badWriter int
type nullWriter int

func (b badWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New(`bang`)
}

func (b nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestFprintf(t *testing.T) {
	assert.Panic(t, func() { util.Fprintf(badWriter(0), `x`) }, `bang`)
}

func TestFprintln(t *testing.T) {
	assert.Equal(t, 5, util.Fprintln(nullWriter(0), `röd`))
	assert.Panic(t, func() { util.Fprintln(badWriter(0), `x`) }, `bang`)
}
