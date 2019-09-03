package util

import (
	"fmt"
	"io"
	"unicode/utf8"
)

func Fprintf(b io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(b, format, args...)
	if err != nil {
		panic(err)
	}
}

func Fprintln(b io.Writer, args ...interface{}) {
	_, err := fmt.Fprintln(b, args...)
	if err != nil {
		panic(err)
	}
}

func WriteByte(b io.Writer, v byte) {
	_, err := b.Write([]byte{v})
	if err != nil {
		panic(err)
	}
}

func WriteRune(b io.Writer, v rune) {
	if v < utf8.RuneSelf {
		WriteByte(b, byte(v))
	} else {
		buf := make([]byte, utf8.UTFMax)
		n := utf8.EncodeRune(buf, v)
		_, err := b.Write(buf[:n])
		if err != nil {
			panic(err)
		}
	}
}

func WriteString(b io.Writer, s string) {
	_, err := io.WriteString(b, s)
	if err != nil {
		panic(err)
	}
}
