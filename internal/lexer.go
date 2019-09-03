package internal

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/lyraproj/dgo/util"
)

type tokenType int

const (
	end = iota
	identifier
	integer
	float
	regexpLiteral
	stringLiteral
	dotdot
)

type token struct {
	s string
	i tokenType
}

func (t token) String() (s string) {
	switch t.i {
	case end:
		s = "end"
	case identifier, integer, float, dotdot:
		s = t.s
	case regexpLiteral:
		sb := &strings.Builder{}
		RegexpSlashQuote(sb, t.s)
		s = sb.String()
	case stringLiteral:
		s = strconv.Quote(t.s)
	default:
		s = fmt.Sprintf(`'%c'`, rune(t.i))
	}
	return
}

func badToken(r rune) error {
	return fmt.Errorf("unexpected character '%c'", r)
}

func nextToken(sr *util.StringReader) (t *token) {
	for {
		r := sr.Next()
		if r == 0 {
			return &token{``, end}
		}

		switch r {
		case ' ', '\t', '\n':
			continue
		case '`':
			t = &token{consumeRawString(sr), stringLiteral}
		case '"':
			t = &token{consumeQuotedString(sr), stringLiteral}
		case '/':
			t = &token{consumeRegexp(sr), regexpLiteral}
		case '.':
			if sr.Peek() == '.' {
				sr.Next()
				t = &token{`..`, dotdot}
			} else {
				t = &token{i: tokenType(r)}
			}
		case '-', '+':
			n := sr.Next()
			if n < '0' || n > '9' {
				panic(badToken(r))
			}
			buf := bytes.NewBufferString(string(r))
			tkn := consumeNumber(sr, n, buf, integer)
			t = &token{buf.String(), tkn}
		default:
			if r >= '0' && r <= '9' {
				buf := bytes.NewBufferString(``)
				tkn := consumeNumber(sr, r, buf, integer)
				t = &token{buf.String(), tkn}
			} else if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
				buf := bytes.NewBufferString(``)
				consumeIdentifier(sr, r, buf)
				t = &token{buf.String(), identifier}
			} else {
				t = &token{i: tokenType(r)}
			}
		}
		break
	}

	return t
}

func consumeUnsignedInteger(sr *util.StringReader, buf *bytes.Buffer) {
	for {
		r := sr.Peek()
		switch r {
		case 0:
		case '.':
			panic(badToken(r))
		default:
			if r >= '0' && r <= '9' {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			if unicode.IsLetter(r) {
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
			if r >= '0' && r <= '9' {
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
			if r >= '0' && r <= '9' || r >= 'A' && r <= 'F' || r >= 'a' && r <= 'f' {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			return
		}
	}
}

func consumeNumber(sr *util.StringReader, start rune, buf *bytes.Buffer, t tokenType) tokenType {
	buf.WriteRune(start)
	firstZero := t != float && start == '0'
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return t
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
				if r >= '0' && r <= '9' || r >= 'A' && r <= 'F' || r >= 'a' && r <= 'f' {
					buf.WriteRune(r)
					consumeHexInteger(sr, buf)
					return t
				}
			}
			panic(badToken(r))
		case '.':
			if sr.Peek2() == '.' {
				return t
			}
			if t == float {
				panic(badToken(r))
			}
			sr.Next()
			buf.WriteRune(r)
			r = sr.Next()
			if r >= '0' && r <= '9' {
				return consumeNumber(sr, r, buf, float)
			}
			panic(badToken(r))
		default:
			if r >= '0' && r <= '9' {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			return t
		}
	}
}

func consumeRegexp(sr *util.StringReader) string {
	buf := bytes.NewBufferString(``)
	for {
		r := sr.Next()
		switch r {
		case '/':
			return buf.String()
		case '\\':
			r = sr.Next()
			switch r {
			case 0:
				panic(errors.New("unterminated regexp"))
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

func consumeQuotedString(sr *util.StringReader) string {
	buf := bytes.NewBufferString(``)
	for {
		r := sr.Next()
		if r == '"' {
			return buf.String()
		}
		switch r {
		default:
			buf.WriteRune(r)
		case 0:
			panic(errors.New("unterminated string"))
		case '\n':
			panic(errors.New("unterminated string"))
		case '\\':
			r = sr.Next()
			switch r {
			default:
				panic(fmt.Errorf("illegal escape '\\%c'", r))
			case 0:
				panic(errors.New("unterminated string"))
			case 'n':
				r = '\n'
			case 'r':
				r = '\r'
			case 't':
				r = '\t'
			case '"':
			case '\\':
			}
			buf.WriteRune(r)
		}
	}
}

func consumeRawString(sr *util.StringReader) string {
	buf := bytes.NewBufferString(``)
	for {
		r := sr.Next()
		if r == '`' {
			return buf.String()
		}
		if r == 0 {
			panic(errors.New("unterminated string"))
		}
		buf.WriteRune(r)
	}
}

func consumeIdentifier(sr *util.StringReader, start rune, buf *bytes.Buffer) {
	buf.WriteRune(start)
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return
		default:
			if r == '_' || r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' {
				sr.Next()
				buf.WriteRune(r)
				continue
			}
			return
		}
	}
}
