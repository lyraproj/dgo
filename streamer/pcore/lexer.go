package pcore

import (
	"bytes"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/lyraproj/dgo/parser"
	"github.com/lyraproj/dgo/util"
)

const (
	end = iota
	name
	identifier
	integer
	float
	regexpLiteral
	stringLiteral
	rocket
)

func tokenTypeString(t int) (s string) {
	switch t {
	case end:
		s = "end"
	case name:
		s = "name"
	case identifier:
		s = "identifier"
	case integer:
		s = "integer"
	case float:
		s = "float"
	case regexpLiteral:
		s = "regexp"
	case stringLiteral:
		s = "string"
	case rocket:
		s = "rocket"
	default:
		s = string(rune(t))
	}
	return
}

func tokenString(t *parser.Token) string {
	return fmt.Sprintf("%s: '%s'", tokenTypeString(t.Type), t.Value)
}

func badToken(r rune) error {
	return fmt.Errorf("unexpected character '%c'", r)
}

func nextToken(sr *util.StringReader) (t *parser.Token) {
	for {
		r := sr.Next()
		if r == utf8.RuneError {
			panic(errors.New("unicode error"))
		}
		if r == 0 {
			return &parser.Token{``, end}
		}

		switch r {
		case ' ', '\t', '\n':
			continue
		case '#':
			consumeLineComment(sr)
			continue
		case '\'', '"':
			t = &parser.Token{consumeString(sr, r), stringLiteral}
		case '/':
			t = &parser.Token{consumeRegexp(sr), regexpLiteral}
		case '=':
			if sr.Peek() == '>' {
				sr.Next()
				t = &parser.Token{`=>`, rocket}
			} else {
				t = &parser.Token{Type: int(r)}
			}
		case '-', '+':
			n := sr.Next()
			if n < '0' || n > '9' {
				panic(badToken(r))
			}
			buf := bytes.NewBufferString(string(r))
			tkn := consumeNumber(sr, n, buf, integer)
			t = &parser.Token{buf.String(), tkn}
		default:
			t = buildToken(r, sr)
		}
		break
	}
	return t
}

func buildToken(r rune, sr *util.StringReader) *parser.Token {
	switch {
	case parser.IsDigit(r):
		buf := bytes.NewBufferString(``)
		tkn := consumeNumber(sr, r, buf, integer)
		return &parser.Token{buf.String(), tkn}
	case parser.IsUpperCase(r):
		buf := bytes.NewBufferString(``)
		consumeTypeName(sr, r, buf)
		return &parser.Token{buf.String(), name}
	case parser.IsLowerCase(r):
		buf := bytes.NewBufferString(``)
		consumeIdentifier(sr, r, buf)
		return &parser.Token{buf.String(), identifier}
	default:
		return &parser.Token{Type: int(r)}
	}
}

func consumeLineComment(sr *util.StringReader) {
	for {
		switch sr.Next() {
		case 0, '\n':
			return
		case utf8.RuneError:
			panic(errors.New("unicode error"))
		}
	}
}

func consumeUnsignedInteger(sr *util.StringReader, buf *bytes.Buffer) {
	for {
		r := sr.Peek()
		switch r {
		case utf8.RuneError:
			panic(errors.New("unicode error"))
		case 0:
		case '.':
			panic(badToken(r))
		default:
			if r >= '0' && r <= '9' {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			if parser.IsLetter(r) {
				sr.Next()
				panic(badToken(r))
			}
			return
		}
	}
}

func consumeExponent(sr *util.StringReader, buf *bytes.Buffer) {
	for {
		r := sr.Next()
		switch r {
		case 0:
			panic(errors.New("unexpected end"))
		case '+', '-':
			buf.WriteRune(r)
			r = sr.Next()
			fallthrough
		default:
			if parser.IsDigit(r) {
				buf.WriteRune(r)
				consumeUnsignedInteger(sr, buf)
				return
			}
			panic(badToken(r))
		}
	}
}

func consumeHexInteger(sr *util.StringReader, buf *bytes.Buffer) {
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return
		default:
			if parser.IsHex(r) {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			return
		}
	}
}

func consumeNumber(sr *util.StringReader, start rune, buf *bytes.Buffer, tt int) int {
	buf.WriteRune(start)
	firstZero := tt != float && start == '0'
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return tt
		case '0':
			sr.Next()
			buf.WriteRune(r)
			continue
		case 'e', 'E':
			sr.Next()
			buf.WriteRune(r)
			consumeExponent(sr, buf)
			return float
		case 'x', 'X':
			if firstZero {
				sr.Next()
				buf.WriteRune(r)
				r = sr.Next()
				if parser.IsHex(r) {
					buf.WriteRune(r)
					consumeHexInteger(sr, buf)
					return tt
				}
			}
			panic(badToken(r))
		case '.':
			if tt == float {
				panic(badToken(r))
			}
			sr.Next()
			buf.WriteRune(r)
			r = sr.Next()
			if parser.IsDigit(r) {
				return consumeNumber(sr, r, buf, float)
			}
			panic(badToken(r))
		default:
			if parser.IsDigit(r) {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			return tt
		}
	}
}

func consumeRegexp(sr *util.StringReader) string {
	buf := bytes.NewBufferString(``)
	for {
		r := sr.Next()
		switch r {
		case utf8.RuneError:
			panic(badToken(r))
		case '/':
			return buf.String()
		case '\\':
			r = sr.Next()
			switch r {
			case 0:
				panic(errors.New("unterminated regexp"))
			case utf8.RuneError:
				panic(badToken(r))
			case '/': // Escape is removed
			default:
				buf.WriteByte('\\')
			}
			buf.WriteRune(r)
		case 0, '\n':
			panic(errors.New("unterminated regexp"))
		default:
			buf.WriteRune(r)
		}
	}
}

func consumeString(sr *util.StringReader, end rune) string {
	buf := bytes.NewBufferString(``)
	for {
		r := sr.Next()
		if r == end {
			return buf.String()
		}
		switch r {
		case 0:
			panic(errors.New("unterminated string"))
		case utf8.RuneError:
			panic(badToken(r))
		case '\\':
			r := sr.Next()
			switch r {
			case 0:
				panic(errors.New("unterminated string"))
			case utf8.RuneError:
				panic(badToken(r))
			case 'n':
				r = '\n'
			case 'r':
				r = '\r'
			case 't':
				r = '\t'
			case '\\':
			default:
				if r != end {
					panic(fmt.Errorf("illegal escape '\\%c'", r))
				}
			}
			buf.WriteRune(r)
		case '\n':
			panic(errors.New("unterminated string"))
		default:
			buf.WriteRune(r)
		}
	}
}

func consumeIdentifier(sr *util.StringReader, start rune, buf *bytes.Buffer) {
	buf.WriteRune(start)
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return
		case ':':
			sr.Next()
			buf.WriteRune(r)
			r = sr.Next()
			if r == ':' {
				buf.WriteRune(r)
				r = sr.Next()
				if r == '_' || parser.IsLowerCase(r) {
					buf.WriteRune(r)
					continue
				}
			}
			panic(badToken(r))
		default:
			if r == '_' || parser.IsLetterOrDigit(r) {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			return
		}
	}
}

func consumeTypeName(sr *util.StringReader, start rune, buf *bytes.Buffer) {
	buf.WriteRune(start)
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return
		case ':':
			sr.Next()
			buf.WriteRune(r)
			r = sr.Next()
			if r == ':' {
				buf.WriteRune(r)
				r = sr.Next()
				if parser.IsUpperCase(r) {
					buf.WriteRune(r)
					continue
				}
			}
			panic(badToken(r))
		default:
			if r == '_' || parser.IsLetterOrDigit(r) {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			return
		}
	}
}
