package util

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// An Indenter helps building strings where all newlines are supposed to be followed by
// a sequence of zero or many spaces that reflect an indent level.
type Indenter struct {
	b *bytes.Buffer
	i int
	s string
}

// An Indentable can create build a string representation of itself using an Indenter
type Indentable interface {
	// AppendTo appends a string representation of the Node to the Indenter
	AppendTo(w *Indenter)
}

// ToString will produce an unindented string from an Indentable
func ToString(ia Indentable) string {
	i := NewIndenter(``)
	ia.AppendTo(i)
	return i.String()
}

// ToIndentedString will produce a string from an Indentable using an Indenter initialized
// with a two space indentation.
func ToIndentedString(ia Indentable) string {
	i := NewIndenter(` `)
	ia.AppendTo(i)
	return i.String()
}

// NewIndenter creates a new Indenter for indent level zero using the given string to perform
// one level of indentation. An empty string will yield unindented output
func NewIndenter(indent string) *Indenter {
	return &Indenter{b: &bytes.Buffer{}, i: 0, s: indent}
}

// Len returns the current number of bytes that has been appended to the indenter
func (i *Indenter) Len() int {
	return i.b.Len()
}

// Level returns the indent level for the indenter
func (i *Indenter) Level() int {
	return i.i
}

// Reset resets the internal buffer. It does not reset the indent
func (i *Indenter) Reset() {
	i.b.Reset()
}

// String returns the current string that has been built using the indenter. Trailing whitespaces
// are deleted from all lines.
func (i *Indenter) String() string {
	n := bytes.NewBuffer(make([]byte, 0, i.b.Len()))
	wb := &bytes.Buffer{}
	for {
		r, _, err := i.b.ReadRune()
		if err == io.EOF {
			break
		}
		if r == ' ' || r == '\t' {
			// Defer whitespace output
			wb.WriteByte(byte(r))
			continue
		}
		if r == '\n' {
			// Truncate trailing space
			wb.Reset()
		} else {
			if wb.Len() > 0 {
				n.Write(wb.Bytes())
				wb.Reset()
			}
		}
		n.WriteRune(r)
	}
	return n.String()
}

// WriteString appends a string to the internal buffer without checking for newlines
func (i *Indenter) WriteString(s string) (n int, err error) {
	return i.b.WriteString(s)
}

// Write appends a slice of bytes to the internal buffer without checking for newlines
func (i *Indenter) Write(p []byte) (n int, err error) {
	return i.b.Write(p)
}

// AppendRune appends a rune to the internal buffer without checking for newlines
func (i *Indenter) AppendRune(r rune) {
	i.b.WriteRune(r)
}

// Append appends a string to the internal buffer without checking for newlines
func (i *Indenter) Append(s string) {
	WriteString(i.b, s)
}

// AppendValue appends the string form of the given value to the internal buffer. Indentation
// will be recursive if the value implements Indentable
func (i *Indenter) AppendValue(v interface{}) {
	if vi, ok := v.(Indentable); ok {
		vi.AppendTo(i)
	} else if vs, ok := v.(fmt.Stringer); ok {
		WriteString(i.b, vs.String())
	} else {
		Fprintf(i.b, "%v", v)
	}
}

// AppendIndented is like Append but replaces all occurrences of newline with an indented newline
func (i *Indenter) AppendIndented(s string) {
	for ni := strings.IndexByte(s, '\n'); ni >= 0; ni = strings.IndexByte(s, '\n') {
		if ni > 0 {
			WriteString(i.b, s[:ni])
		}
		i.NewLine()
		ni++
		if ni >= len(s) {
			return
		}
		s = s[ni:]
	}
	if len(s) > 0 {
		WriteString(i.b, s)
	}
}

// AppendBool writes the string "true" or "false" to the internal buffer
func (i *Indenter) AppendBool(b bool) {
	var s string
	if b {
		s = `true`
	} else {
		s = `false`
	}
	WriteString(i.b, s)
}

// AppendInt writes the result of calling strconf.Itoa() in the given argument
func (i *Indenter) AppendInt(b int) {
	WriteString(i.b, strconv.Itoa(b))
}

// Indent returns a new Indenter instance that shares the same buffer but has an
// indent level that is increased by one.
func (i *Indenter) Indent() *Indenter {
	c := *i
	c.i++
	return &c
}

// Indenting returns true if this indenter has an indent string with a length > 0
func (i *Indenter) Indenting() bool {
	return len(i.s) > 0
}

// Printf formats according to a format specifier and writes to the internal buffer.
func (i *Indenter) Printf(s string, args ...interface{}) {
	Fprintf(i.b, s, args...)
}

// NewLine writes a newline followed by the current indent after trimming trailing whitespaces
func (i *Indenter) NewLine() {
	if len(i.s) > 0 {
		i.b.WriteByte('\n')
		for n := 0; n < i.i; n++ {
			WriteString(i.b, i.s)
		}
	}
}
