package util

import (
	"fmt"
	"io"
	"unicode/utf8"
)

// Fprintf is like fmt.Fprintf but it panics in case of error instead of returning it.
func Fprintf(b io.Writer, format string, args ...interface{}) int {
	n, err := fmt.Fprintf(b, format, args...)
	if err != nil {
		panic(err)
	}
	return n
}

// Fprintln is like fmt.Fprintln but it panics in case of error instead of returning it.
func Fprintln(b io.Writer, args ...interface{}) int {
	n, err := fmt.Fprintln(b, args...)
	if err != nil {
		panic(err)
	}
	return n
}

// WriteByte writes the given byte on the given Writer and panics if an error occurs
func WriteByte(b io.Writer, v byte) {
	_, err := b.Write([]byte{v})
	if err != nil {
		panic(err)
	}
}

// WriteRune writes the given rune on the given Writer and panics if an error occurs. It
// returns the number of bytes written.
func WriteRune(b io.Writer, v rune) int {
	n := 1
	if v < utf8.RuneSelf {
		WriteByte(b, byte(v))
	} else {
		buf := make([]byte, utf8.UTFMax)
		var err error
		n, err = b.Write(buf[:utf8.EncodeRune(buf, v)])
		if err != nil {
			panic(err)
		}
	}
	return n
}

// WriteString writes the given string on the given Writer and panics if an error occurs. It
// returns the number of bytes written.
func WriteString(b io.Writer, s string) int {
	n, err := io.WriteString(b, s)
	if err != nil {
		panic(err)
	}
	return n
}

// WriteQuotedString writes the given string enclosed in double quotes on the given Writer and panics
// if an error occurs.
func WriteQuotedString(b io.Writer, s string) {
	WriteByte(b, '"')
	for _, c := range s {
		switch c {
		case '\t':
			WriteString(b, `\t`)
		case '\n':
			WriteString(b, `\n`)
		case '\r':
			WriteString(b, `\r`)
		case '"':
			WriteString(b, `\"`)
		case '\\':
			WriteString(b, `\\`)
		case '\a':
			WriteString(b, `\a`)
		case '\b':
			WriteString(b, `\b`)
		case '\f':
			WriteString(b, `\f`)
		case '\v':
			WriteString(b, `\v`)
		default:
			if c < 0x20 {
				Fprintf(b, `\x%.2X`, c)
			} else {
				WriteRune(b, c)
			}
		}
	}
	WriteByte(b, '"')
}
