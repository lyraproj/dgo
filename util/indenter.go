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
type Indenter interface {
	// Append appends a string to the internal buffer without checking for newlines
	Append(string)

	// AppendRune appends a rune to the internal buffer without checking for newlines
	AppendRune(rune)

	// AppendValue appends the string form of the given value to the internal buffer. Indentation
	// will be recursive if the value implements Indentable
	AppendValue(interface{})

	// AppendIndented is like Append but replaces all occurrences of newline with an indented newline
	AppendIndented(string)

	// AppendBool writes the string "true" or "false" to the internal buffer
	AppendBool(bool)

	// AppendInt writes the result of calling strconf.Itoa() in the given argument
	AppendInt(int)

	// Indent returns a new indenter instance that shares the same buffer but has an
	// indent level that is increased by one.
	Indent() Indenter

	// Indenting returns true if this indenter has an indent string with a length > 0
	Indenting() bool

	// Len returns the current number of bytes that has been appended to the indenter
	Len() int

	// Level returns the indent level for the indenter
	Level() int

	// NewLine writes a newline followed by the current indent after trimming trailing whitespaces
	NewLine()

	// Printf formats according to a format specifier and writes to the internal buffer.
	Printf(string, ...interface{})

	// Reset resets the internal buffer. It does not reset the indent
	Reset()

	// String returns the current string that has been built using the indenter. Trailing whitespaces
	// are deleted from all lines.
	String() string

	// Write appends a slice of bytes to the internal buffer without checking for newlines
	Write(p []byte) (n int, err error)

	// WriteString appends a string to the internal buffer without checking for newlines
	WriteString(s string) (n int, err error)
}

type indenter struct {
	b *bytes.Buffer
	i int
	s string
}

type erpIndenter struct {
	indenter
	seen []Indentable
}

// An Indentable can create build a string representation of itself using an indenter
type Indentable interface {
	// AppendTo appends a string representation of the Node to the indenter
	AppendTo(w Indenter)
}

// ToString will produce an unindented string from an Indentable
func ToString(ia Indentable) string {
	i := NewIndenter(``)
	ia.AppendTo(i)
	return i.String()
}

// ToStringERP will produce an unindented string from an Indentable using an indenter returned
// by NewERPIndenter
func ToStringERP(ia Indentable) string {
	i := NewERPIndenter(``)
	ia.AppendTo(i)
	return i.String()
}

// ToIndentedString will produce a string from an Indentable using an indenter initialized
// with a two space indentation.
func ToIndentedString(ia Indentable) string {
	i := NewIndenter(` `)
	ia.AppendTo(i)
	return i.String()
}

// ToIndentedStringERP will produce a string from an Indentable using an indenter returned
// by NewERPIndenter that has been initialized with a two space indentation.
func ToIndentedStringERP(ia Indentable) string {
	i := NewERPIndenter(` `)
	ia.AppendTo(i)
	return i.String()
}

// NewIndenter creates a new indenter for indent level zero using the given string to perform
// one level of indentation. An empty string will yield unindented output
func NewIndenter(indent string) Indenter {
	return &indenter{b: &bytes.Buffer{}, i: 0, s: indent}
}

// NewERPIndenter creates an endless recursion protected indenter capable of indenting self referencing
// values. When an endless recursion is encountered, the string <recursive self reference> is emitted
// rather than the value itself.
func NewERPIndenter(indent string) Indenter {
	return &erpIndenter{indenter: indenter{b: &bytes.Buffer{}, i: 0, s: indent}}
}

func (i *indenter) Len() int {
	return i.b.Len()
}

func (i *indenter) Level() int {
	return i.i
}

func (i *indenter) Reset() {
	i.b.Reset()
}

func (i *indenter) String() string {
	n := bytes.NewBuffer(make([]byte, 0, i.b.Len()))
	wb := &bytes.Buffer{}
	for {
		r, _, err := i.b.ReadRune()
		if err == io.EOF {
			break
		}
		if r == ' ' || r == '\t' {
			// Defer whitespace output
			WriteByte(wb, byte(r))
			continue
		}
		if r == '\n' {
			// Truncate trailing space
			wb.Reset()
		} else if wb.Len() > 0 {
			_, _ = n.Write(wb.Bytes())
			wb.Reset()
		}
		WriteRune(n, r)
	}
	return n.String()
}

func (i *indenter) WriteString(s string) (n int, err error) {
	return i.b.WriteString(s)
}

func (i *indenter) Write(p []byte) (n int, err error) {
	return i.b.Write(p)
}

func (i *indenter) AppendRune(r rune) {
	WriteRune(i.b, r)
}

func (i *indenter) Append(s string) {
	WriteString(i.b, s)
}

func (i *indenter) AppendValue(v interface{}) {
	if vi, ok := v.(Indentable); ok {
		vi.AppendTo(i)
	} else if vs, ok := v.(fmt.Stringer); ok {
		WriteString(i.b, vs.String())
	} else {
		Fprintf(i.b, "%v", v)
	}
}

func (i *indenter) AppendIndented(s string) {
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

func (i *indenter) AppendBool(b bool) {
	var s string
	if b {
		s = `true`
	} else {
		s = `false`
	}
	WriteString(i.b, s)
}

func (i *indenter) AppendInt(b int) {
	WriteString(i.b, strconv.Itoa(b))
}

func (i *indenter) Indent() Indenter {
	c := *i
	c.i++
	return &c
}

func (i *indenter) Indenting() bool {
	return len(i.s) > 0
}

func (i *indenter) Printf(s string, args ...interface{}) {
	Fprintf(i.b, s, args...)
}

func (i *indenter) NewLine() {
	if len(i.s) > 0 {
		WriteByte(i.b, '\n')
		for n := 0; n < i.i; n++ {
			WriteString(i.b, i.s)
		}
	}
}

func (i *erpIndenter) Indent() Indenter {
	c := *i
	c.i++
	return &c
}

func (i *erpIndenter) AppendValue(v interface{}) {
	if vi, ok := v.(Indentable); ok {
		s := i.seen
		for n := range s {
			if s[n] == v {
				WriteString(i.b, `<recursive self reference>`)
				return
			}
		}
		i.seen = append(i.seen, vi)
		vi.AppendTo(i)
		i.seen = s
	} else if vs, ok := v.(fmt.Stringer); ok {
		WriteString(i.b, vs.String())
	} else {
		Fprintf(i.b, "%v", v)
	}
}
