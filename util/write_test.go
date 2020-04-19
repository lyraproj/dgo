package util_test

import (
	"errors"
	"testing"

	require "github.com/tada/dgo/dgo_test"
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
	require.Panic(t, func() { util.Fprintf(badWriter(0), `x`) }, `bang`)
}

func TestFprintln(t *testing.T) {
	require.Equal(t, 5, util.Fprintln(nullWriter(0), `röd`))
	require.Panic(t, func() { util.Fprintln(badWriter(0), `x`) }, `bang`)
}

func TestWriteByte(t *testing.T) {
	require.Panic(t, func() { util.WriteByte(badWriter(0), 'x') }, `bang`)
}

func TestWriteRune(t *testing.T) {
	require.Equal(t, 2, util.WriteRune(nullWriter(0), 'ö'))
	require.Panic(t, func() { util.WriteRune(badWriter(0), 'x') }, `bang`)
	require.Panic(t, func() { util.WriteRune(badWriter(0), 'ö') }, `bang`)
}

func TestWriteString(t *testing.T) {
	require.Equal(t, 4, util.WriteString(nullWriter(0), `röd`))
	require.Panic(t, func() { util.WriteString(badWriter(0), `röd`) }, `bang`)
}
