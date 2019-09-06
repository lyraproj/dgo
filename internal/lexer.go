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
	dotdotdot
)

type token struct {
	s string
	i tokenType
}

func (t *token) String() (s string) {
	if t == nil || t.i == end {
		return "EOT"
	}
	switch t.i {
	case identifier, integer, float, dotdot, dotdotdot:
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
	if r == 0 {
		return errors.New(`unexpected end`)
	}
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
				if sr.Peek() == '.' {
					sr.Next()
					t = &token{`...`, dotdotdot}
				} else {
					t = &token{`..`, dotdot}
				}
			} else {
				t = &token{i: tokenType(r)}
			}
		case '-', '+':
			n := sr.Next()
			if n < '0' || n > '9' {
				panic(badToken(n))
			}
			buf := bytes.NewBufferString(``)
			if r == '-' {
				util.WriteRune(buf, r)
			}
			tkn := consumeNumber(sr, n, buf, integer)
			return &token{buf.String(), tkn}
		default:
			t = buildToken(r, sr)
		}
		break
	}
	return t
}

func buildToken(r rune, sr *util.StringReader) *token {
	switch {
	case r >= '0' && r <= '9':
		buf := bytes.NewBufferString(``)
		tkn := consumeNumber(sr, r, buf, integer)
		return &token{buf.String(), tkn}
	case r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z':
		buf := bytes.NewBufferString(``)
		consumeIdentifier(sr, r, buf)
		return &token{buf.String(), identifier}
	default:
		return &token{i: tokenType(r)}
	}
}

func consumeUnsignedInteger(sr *util.StringReader, buf *bytes.Buffer) {
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return
		case '.':
			panic(badToken(r))
		default:
			if r >= '0' && r <= '9' {
				sr.Next()
				util.WriteRune(buf, r)
				continue
			}
			if unicode.IsLetter(r) {
				panic(badToken(r))
			}
			return
		}
	}
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHex(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'A' && r <= 'F' || r >= 'a' && r <= 'f'
}

func consumeExponent(sr *util.StringReader, buf *bytes.Buffer) {
	for {
		r := sr.Next()
		switch r {
		case 0:
			panic(errors.New("unexpected end"))
		case '+', '-':
			util.WriteRune(buf, r)
			r = sr.Next()
			fallthrough
		default:
			if isDigit(r) {
				util.WriteRune(buf, r)
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
			if isHex(r) {
				sr.Next()
				util.WriteRune(buf, r)
				continue
			}
			return
		}
	}
}

func consumeNumber(sr *util.StringReader, start rune, buf *bytes.Buffer, t tokenType) tokenType {
	util.WriteRune(buf, start)
	firstZero := t != float && start == '0'

	for r := sr.Peek(); r != 0; r = sr.Peek() {
		switch r {
		case '0':
			sr.Next()
			util.WriteRune(buf, r)
		case 'e', 'E':
			sr.Next()
			util.WriteRune(buf, r)
			consumeExponent(sr, buf)
			return float
		case 'x', 'X':
			if firstZero {
				sr.Next()
				util.WriteRune(buf, r)
				r = sr.Next()
				if isHex(r) {
					util.WriteRune(buf, r)
					consumeHexInteger(sr, buf)
					return t
				}
			}
			panic(badToken(r))
		case '.':
			if sr.Peek2() == '.' {
				return t
			}
			if t != float {
				sr.Next()
				util.WriteRune(buf, r)
				r = sr.Next()
				if isDigit(r) {
					return consumeNumber(sr, r, buf, float)
				}
			}
			panic(badToken(r))
		default:
			if !isDigit(r) {
				return t
			}
			sr.Next()
			util.WriteRune(buf, r)
		}
	}
	return t
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
				util.WriteRune(buf, '\\')
			}
			util.WriteRune(buf, r)
		case 0, '\n':
			panic(errors.New("unterminated regexp"))
		default:
			util.WriteRune(buf, r)
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
			util.WriteRune(buf, r)
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
			util.WriteRune(buf, r)
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
		util.WriteRune(buf, r)
	}
}

func consumeIdentifier(sr *util.StringReader, start rune, buf *bytes.Buffer) {
	util.WriteRune(buf, start)
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return
		default:
			if r == '_' || r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' {
				sr.Next()
				util.WriteRune(buf, r)
				continue
			}
			return
		}
	}
}
