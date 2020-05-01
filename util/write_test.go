package util_test

import (
	"bytes"
	"errors"
	"strconv"
	"testing"

	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/util"
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

func TestWriteByte(t *testing.T) {
	assert.Panic(t, func() { util.WriteByte(badWriter(0), 'x') }, `bang`)
}

func TestWriteRune(t *testing.T) {
	assert.Equal(t, 2, util.WriteRune(nullWriter(0), 'ö'))
	assert.Panic(t, func() { util.WriteRune(badWriter(0), 'x') }, `bang`)
	assert.Panic(t, func() { util.WriteRune(badWriter(0), 'ö') }, `bang`)
}

func TestWriteString(t *testing.T) {
	assert.Equal(t, 4, util.WriteString(nullWriter(0), `röd`))
	assert.Panic(t, func() { util.WriteString(badWriter(0), `röd`) }, `bang`)
}
func TestWriteQuotedString(t *testing.T) {
	b := bytes.Buffer{}
	theString := "Quote \", BS \\, NewLine \n, Tab \t, VT \v, CR \r, \a, \b, \f, \x04, \u1234"
	util.WriteQuotedString(&b, theString)
	assert.Equal(t, strconv.Quote(theString), b.String())
}
