package util_test

import (
	"fmt"
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgo/util"
	"github.com/tada/dgo/vf"
)

func ExampleToString() {
	v := vf.Values(`a`, map[string]int{`b`: 3}).(dgo.Array)
	fmt.Println(util.ToString(v))
	// Output: {"a",{"b":3}}
}

func ExampleToIndentedString() {
	v := vf.Values(`a`, map[string]int{`b`: 3}).(dgo.Array)
	fmt.Println(util.ToIndentedString(v))
	// Output:
	// {
	//   "a",
	//   {
	//     "b": 3
	//   }
	// }
}

func TestIndenter_AppendBool(t *testing.T) {
	i := util.NewIndenter(``)
	i.AppendBool(true)
	i.AppendRune(',')
	i.AppendBool(false)
	assert.Equal(t, `true,false`, i.String())
}

func TestIndenter_AppendInt(t *testing.T) {
	i := util.NewIndenter(``)
	i.AppendInt(-31)
	assert.Equal(t, `-31`, i.String())
}

func TestIndenter_AppendIndented_onlyNewline(t *testing.T) {
	i := util.NewIndenter(`  `).Indent()
	i.AppendIndented("\n")
	assert.Equal(t, "\n", i.String())
}

func TestIndenter_AppendIndented_twoNewline(t *testing.T) {
	i := util.NewIndenter(`  `).Indent()
	i.AppendIndented("\n\n")
	assert.Equal(t, "\n\n", i.String())
}

func TestIndenter_AppendIndented_twoLines(t *testing.T) {
	i := util.NewIndenter(`  `).Indent()
	i.AppendIndented("\na\n")
	assert.Equal(t, "\n  a\n", i.String())

	i.Reset()
	i.AppendIndented("\na\nb")
	assert.Equal(t, "\n  a\n  b", i.String())
}

func TestIndenter_Len(t *testing.T) {
	i := util.NewIndenter(``)
	i.Append(`one`)
	assert.Equal(t, 3, i.Len())
}

func TestIndenter_Level(t *testing.T) {
	i := util.NewIndenter(``)
	assert.Equal(t, 0, i.Level())
	assert.Equal(t, 1, i.Indent().Level())
}

func TestIndenter_Printf(t *testing.T) {
	i := util.NewIndenter(``)
	i.Printf(`the %s %d`, `digit`, 12)
	assert.Equal(t, "the digit 12", i.String())
}

func TestIndenter_Reset(t *testing.T) {
	i := util.NewIndenter(``)
	i.Append(`one`)
	assert.Equal(t, 3, i.Len())
	i.Reset()
	assert.Equal(t, 0, i.Len())
}

func TestIndenter_WriteString(t *testing.T) {
	i := util.NewIndenter(``)
	n, err := i.WriteString(`one`)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
}

func TestIndenter_Write(t *testing.T) {
	i := util.NewIndenter(``)
	n, err := i.Write([]byte(`one`))
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
}
