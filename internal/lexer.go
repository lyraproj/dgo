package internal

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

const (
	digit = 1 << iota
	hexDigit
	letter
	idCharStart
	idChar
	exprEnd
)

var charTypes = [128]uint8{}

func init() {
	i := int('$')
	charTypes[i] = idCharStart | idChar
	i = int('_')
	charTypes[i] = idCharStart | idChar
	for i = '0'; i <= '9'; i++ {
		charTypes[i] = digit | hexDigit | idChar
	}
	for i = 'A'; i <= 'F'; i++ {
		charTypes[i] = hexDigit | letter | idCharStart | idChar
	}
	for i = 'G'; i <= 'Z'; i++ {
		charTypes[i] = letter | idCharStart | idChar
	}
	for i = 'a'; i <= 'f'; i++ {
		charTypes[i] = hexDigit | letter | idCharStart | idChar
	}
	for i = 'g'; i <= 'z'; i++ {
		charTypes[i] = letter | idCharStart | idChar
	}
	for _, i := range []int{end, ')', '}', ']', ',', ':', '?', '|', '&', '^', '.'} {
		charTypes[i] = exprEnd
	}
}

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
		switch r {
		case 0:
			return &token{``, end}
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
			if !isDigit(n) {
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
	case isDigit(r):
		buf := bytes.NewBufferString(``)
		tkn := consumeNumber(sr, r, buf, integer)
		return &token{buf.String(), tkn}
	case isIdentifierStart(r):
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
		switch {
		case r == '.' || isLetter(r):
			panic(badToken(r))
		case isDigit(r):
			sr.Next()
			util.WriteRune(buf, r)
		default:
			return
		}
	}
}

func isExpressionEnd(r rune) bool {
	return r < 128 && (charTypes[r]&exprEnd) != 0
}

func isDigit(r rune) bool {
	return r < 128 && (charTypes[r]&digit) != 0
}

func isLetter(r rune) bool {
	return r < 128 && (charTypes[r]&letter) != 0
}

func isHex(r rune) bool {
	return r < 128 && (charTypes[r]&hexDigit) != 0
}

func isIdentifier(r rune) bool {
	return r < 128 && (charTypes[r]&idChar) != 0
}

func isIdentifierStart(r rune) bool {
	return r < 128 && (charTypes[r]&idCharStart) != 0
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
	for isHex(sr.Peek()) {
		util.WriteRune(buf, sr.Next())
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
		switch r {
		default:
			util.WriteRune(buf, r)
		case 0, '\n':
			panic(errors.New("unterminated string"))
		case '"':
			return buf.String()
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
	for isIdentifier(sr.Peek()) {
		util.WriteRune(buf, sr.Next())
	}
}
