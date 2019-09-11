package util_test

import (
	"fmt"
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
	"github.com/lyraproj/dgo/vf"
)

func ExampleToString() {
	v := vf.Values(`a`, map[string]int{`b`: 3}).(dgo.Array)
	fmt.Println(util.ToString(v))
	// Output: ["a",{"b":3}]
}

func ExampleToIndentedString() {
	v := vf.Values(`a`, map[string]int{`b`: 3}).(dgo.Array)
	fmt.Println(util.ToIndentedString(v))
	// Output:
	// [
	//  "a",
	//  {
	//   "b": 3
	//  }
	// ]
}

func TestIndenter_AppendBool(t *testing.T) {
	i := util.NewIndenter(``)
	i.AppendBool(true)
	i.AppendRune(',')
	i.AppendBool(false)
	require.Equal(t, `true,false`, i.String())
}

func TestIndenter_AppendInt(t *testing.T) {
	i := util.NewIndenter(``)
	i.AppendInt(-31)
	require.Equal(t, `-31`, i.String())
}

func TestIndenter_AppendIndented_onlyNewline(t *testing.T) {
	i := util.NewIndenter(`  `).Indent()
	i.AppendIndented("\n")
	require.Equal(t, "\n", i.String())
}

func TestIndenter_AppendIndented_twoNewline(t *testing.T) {
	i := util.NewIndenter(`  `).Indent()
	i.AppendIndented("\n\n")
	require.Equal(t, "\n\n", i.String())
}

func TestIndenter_AppendIndented_twoLines(t *testing.T) {
	i := util.NewIndenter(`  `).Indent()
	i.AppendIndented("\na\n")
	require.Equal(t, "\n  a\n", i.String())

	i.Reset()
	i.AppendIndented("\na\nb")
	require.Equal(t, "\n  a\n  b", i.String())
}

func TestIndenter_Len(t *testing.T) {
	i := util.NewIndenter(``)
	i.Append(`one`)
	require.Equal(t, 3, i.Len())
}

func TestIndenter_Level(t *testing.T) {
	i := util.NewIndenter(``)
	require.Equal(t, 0, i.Level())
	require.Equal(t, 1, i.Indent().Level())
}

func TestIndenter_Reset(t *testing.T) {
	i := util.NewIndenter(``)
	i.Append(`one`)
	require.Equal(t, 3, i.Len())
	i.Reset()
	require.Equal(t, 0, i.Len())
}

func TestIndenter_WriteString(t *testing.T) {
	i := util.NewIndenter(``)
	n, err := i.WriteString(`one`)
	require.Ok(t, err)
	require.Equal(t, 3, n)
}

func TestIndenter_Write(t *testing.T) {
	i := util.NewIndenter(``)
	n, err := i.Write([]byte(`one`))
	require.Ok(t, err)
	require.Equal(t, 3, n)
}
